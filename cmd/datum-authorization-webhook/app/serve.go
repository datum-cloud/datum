package app

import (
	"context"
	"fmt"

	"buf.build/gen/go/datum-cloud/iam/grpc/go/datum/iam/v1alpha/iamv1alphagrpc"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"k8s.io/api/authentication/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

func serveCommand() *cobra.Command {
	var iamEndpoint, certDir, certFile, keyFile string
	var iamInsecure bool

	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the Authorization Webhook API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if iamEndpoint == "" {
				return fmt.Errorf("`--iam-endpoint` is required")
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

			dialOptions := []grpc.DialOption{
				grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
				grpc.WithChainUnaryInterceptor(
					func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
						logger := log.Log.WithName("grpc_client")
						logger.Info(method, ".request: ", protojson.Format(req.(proto.Message)))
						err := invoker(ctx, method, req, reply, cc, opts...)
						logger.Info(method, ".response: ", protojson.Format(reply.(proto.Message)))
						return err
					},
				),
			}

			if iamInsecure {
				dialOptions = append(dialOptions, grpc.WithTransportCredentials(insecure.NewCredentials()))
			} else {
				dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
			}

			iamConnection, err := grpc.NewClient(iamEndpoint, dialOptions...)
			if err != nil {
				return fmt.Errorf("failed to create new IAM client: %w", err)
			}
			defer iamConnection.Close()

			log.SetLogger(zap.New(zap.JSONEncoder()))

			iamClient := iamv1alphagrpc.NewAccessCheckClient(iamConnection)

			entryLog := log.Log.WithName("entrypoint")

			restConfig, err := config.GetConfig()
			if err != nil {
				return fmt.Errorf("failed to get rest config: %s", err)
			}

			runtimeScheme := runtime.NewScheme()
			v1beta1.AddToScheme(runtimeScheme)

			// Setup a Manager
			entryLog.Info("setting up manager")
			mgr, err := manager.New(restConfig, manager.Options{
				Scheme: runtimeScheme,
				Metrics: server.Options{
					BindAddress: ":8999",
				},
				WebhookServer: webhook.NewServer(webhook.Options{
					CertDir:  certDir,
					CertName: certFile,
					KeyName:  keyFile,
				}),
			})
			if err != nil {
				return fmt.Errorf("failed to setup manager: %s", err)
			}

			// Setup webhooks
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
	serveCmd.Flags().BoolVar(&iamInsecure, "iam-endpoint-insecure", false, "Whether the use an insecure connection when export")

	serveCmd.Flags().StringVar(&certDir, "cert-dir", "", "Directory that contains the TLS certs to use for serving the webhook")
	serveCmd.Flags().StringVar(&certFile, "cert-file", "", "Filename in the directory that contains the TLS cert")
	serveCmd.Flags().StringVar(&keyFile, "key-file", "", "Filename in the directory that contains the TLS private key")

	return serveCmd
}
