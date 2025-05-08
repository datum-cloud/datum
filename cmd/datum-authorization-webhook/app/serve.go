package app

import (
	"context"
	"fmt"
	"strings"

	"buf.build/gen/go/datum-cloud/iam/grpc/go/datum/iam/v1alpha/iamv1alphagrpc"
	"go.datumapis.com/datum/cmd/datum-authorization-webhook/app/internal/auth"
	"go.datumapis.com/datum/cmd/datum-authorization-webhook/app/internal/iam"
	authwebhook "go.datumapis.com/datum/cmd/datum-authorization-webhook/app/internal/webhook"

	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"k8s.io/api/authentication/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	crlog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func serveCommand() *cobra.Command {
	var iamEndpoint, certDir, certFile, keyFile string
	var iamInsecure bool
	var iamAuthKeyFile, iamAuthTokenURL, iamAuthJwtAudience, iamAuthTokenScopes string

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the Authorization Webhook API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if iamEndpoint == "" {
				return fmt.Errorf("`--iam-endpoint` is required")
			}

			if iamAuthKeyFile != "" {
				if iamAuthTokenURL == "" {
					return fmt.Errorf("`--iam-auth-token-url` is required when `--iam-auth-key-file` is provided")
				}
				if iamAuthJwtAudience == "" {
					return fmt.Errorf("`--iam-auth-jwt-audience` is required when `--iam-auth-key-file` is provided")
				}
			}

			exporter, err := otlptrace.New(cmd.Context(), otlptracegrpc.NewClient())
			if err != nil {
				return err
			}

			otel.SetTracerProvider(trace.NewTracerProvider(
				trace.WithSampler(trace.AlwaysSample()),
				trace.WithResource(resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceName("authorization-webhook.datumapis.com"),
				)),
				trace.WithBatcher(exporter),
			))
			otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			))

			appLogger := crlog.Log.WithName("datum-auth-webhook")

			dialOptions := []grpc.DialOption{
				grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
			}

			if iamAuthKeyFile != "" {
				appLogger.Info("Configuring IAM authentication using private key JWT from key file", "keyFile", iamAuthKeyFile)

				scopes := []string{}
				if iamAuthTokenScopes != "" {
					scopes = strings.Split(iamAuthTokenScopes, ",")
				}

				jwtTokenSource, err := auth.NewPrivateKeyJwtTokenSource(
					appLogger.WithName("private-key-jwt-token-source"),
					iamAuthKeyFile,
					iamAuthTokenURL,
					iamAuthJwtAudience,
					scopes,
					nil,
				)
				if err != nil {
					return fmt.Errorf("failed to create private key JWT token source: %w", err)
				}

				cachedTokenSource := oauth2.ReuseTokenSource(nil, jwtTokenSource)

				dialOptions = append(dialOptions, grpc.WithPerRPCCredentials(oauth.TokenSource{TokenSource: cachedTokenSource}))
				appLogger.Info("Successfully configured private key JWT token source for IAM authentication")

				// Transport security is handled after this block
			} else {
				appLogger.Info("Using IAM endpoint without service account JWT authentication")
				// No specific PerRPCCredentials in this case. Transport security handled below.
			}

			// Configure transport credentials (TLS or insecure) - applied regardless of Zitadel auth
			if iamInsecure {
				appLogger.Info("Using insecure transport to IAM endpoint.")
				dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
			} else {
				appLogger.Info("Using secure transport (TLS) to IAM endpoint.")
				// For production, ensure proper TLS config (e.g., system cert pool or specific CA)
				dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
			}

			iamConnection, err := grpc.NewClient(iamEndpoint, dialOptions...)
			if err != nil {
				return fmt.Errorf("failed to create new IAM client: %w", err)
			}
			defer iamConnection.Close()

			crlog.SetLogger(zap.New(zap.JSONEncoder()))

			iamClient := iamv1alphagrpc.NewAccessCheckClient(iamConnection)

			entryLog := appLogger.WithName("entrypoint")

			restConfig, err := config.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to get rest config: %s", err)
			}

			runtimeScheme := runtime.NewScheme()
			v1beta1.AddToScheme(runtimeScheme)

			entryLog.Info("setting up manager")

			// Prepare webhook options
			// If a certDir is specified but certFile or keyFile are not, default them
			// to tls.crt and tls.key, as controller-runtime often expects these for auto-generation.
			actualCertFile := certFile
			actualKeyFile := keyFile

			if certDir != "" {
				if actualCertFile == "" {
					actualCertFile = "tls.crt"
					entryLog.Info("Defaulting cert-file to tls.crt in specified cert-dir", "certDir", certDir)
				}
				if actualKeyFile == "" {
					actualKeyFile = "tls.key"
					entryLog.Info("Defaulting key-file to tls.key in specified cert-dir", "certDir", certDir)
				}
			}

			mgr, err := manager.New(restConfig, manager.Options{
				Scheme: runtimeScheme,
				Metrics: server.Options{
					BindAddress: ":8999",
				},
				WebhookServer: webhook.NewServer(webhook.Options{
					CertDir:  certDir,        // User-specified directory
					CertName: actualCertFile, // Defaulted if not provided by user
					KeyName:  actualKeyFile,  // Defaulted if not provided by user
				}),
			})
			if err != nil {
				return fmt.Errorf("failed to setup manager: %s", err)
			}

			entryLog.Info("setting up webhook server")
			hookServer := mgr.GetWebhookServer()

			entryLog.Info("registering webhooks to the webhook server")

			hookServer.Register("/project/v1alpha/projects/{project}/webhook", authwebhook.NewAuthorizerWebhook(&iam.ProjectControlPlaneAuthorizer{
				IAMClient: iamClient,
			}))
			hookServer.Register("/core/v1alpha/webhook", authwebhook.NewAuthorizerWebhook(&iam.CoreControlPlaneAuthorizer{
				IAMClient: iamClient,
			}))

			return mgr.Start(context.Background())
		},
	}

	serveCmd.Flags().StringVar(&iamEndpoint, "iam-endpoint", "", "Endpoint to use for connecting to the datum gRPC API endpoint")
	serveCmd.Flags().BoolVar(&iamInsecure, "iam-endpoint-insecure", false, "Whether to use an insecure connection to the IAM endpoint (transport layer only)")

	serveCmd.Flags().StringVar(&iamAuthKeyFile, "iam-auth-key-file", "", "Path to the service account key JSON file for IAM authentication (using private key JWT).")
	serveCmd.Flags().StringVar(&iamAuthTokenURL, "iam-auth-token-url", "", "URL of the OAuth2 token endpoint. Required if --iam-auth-key-file is set.")
	serveCmd.Flags().StringVar(&iamAuthJwtAudience, "iam-auth-jwt-audience", "", "Audience for the JWT assertion sent to the token endpoint. Required if --iam-auth-key-file is set.")
	serveCmd.Flags().StringVar(&iamAuthTokenScopes, "iam-auth-token-scopes", "openid", "Comma-separated list of OAuth scopes to request for the access token.")

	serveCmd.Flags().StringVar(&certDir, "cert-dir", "", "Directory that contains the TLS certs to use for serving the webhook")
	serveCmd.Flags().StringVar(&certFile, "cert-file", "", "Filename in the directory that contains the TLS cert")
	serveCmd.Flags().StringVar(&keyFile, "key-file", "", "Filename in the directory that contains the TLS private key")

	return serveCmd
}
