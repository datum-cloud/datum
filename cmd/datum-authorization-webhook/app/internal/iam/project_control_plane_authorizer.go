package iam

import (
	"context"
	"fmt"
	"log/slog"

	"buf.build/gen/go/datum-cloud/iam/grpc/go/datum/iam/v1alpha/iamv1alphagrpc"
	iampb "buf.build/gen/go/datum-cloud/iam/protocolbuffers/go/datum/iam/v1alpha"
	"go.datumapis.com/datum/cmd/datum-authorization-webhook/app/internal/webhook"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

var _ authorizer.Authorizer = &ProjectControlPlaneAuthorizer{}

type ProjectControlPlaneAuthorizer struct {
	IAMClient iamv1alphagrpc.AccessCheckClient
}

// Contains a mapping of Kubernetes APIGroups to the service name that should be
// used by the webhook to perform authorization checks.
var serviceNameMapping = map[string]string{
	// An empty APIGroup is used for the core/v1 Kubernetes API Group.
	"": "core.datumapis.com",
}

// Authorize implements authorizer.Authorizer.
func (o *ProjectControlPlaneAuthorizer) Authorize(
	ctx context.Context, attributes authorizer.Attributes,
) (authorizer.Decision, string, error) {

	ctx, span := otel.Tracer("go.datum.net/datum/cmd/datum-authorization-webhook").Start(ctx, "datum.authz-webhook.Authorize", trace.WithAttributes(
		attribute.String("subject", attributes.GetUser().GetName()),
	))
	defer span.End()

	var projectName string
	if projectNames, set := attributes.GetUser().GetExtra()[webhook.ProjectExtraKey]; !set {
		span.SetStatus(codes.Error, "no project ID present in webhook request")
		return authorizer.DecisionDeny, "", fmt.Errorf("extra '%s' is required by core control plane authorizer", webhook.ProjectExtraKey)
	} else if len(projectNames) > 1 {
		span.SetStatus(codes.Error, "multiple project IDs present in webhook request")
		return authorizer.DecisionDeny, "", fmt.Errorf("extra '%s' only supports one value, but multiple were provided: %v", webhook.ProjectExtraKey, projectNames)
	} else {
		projectName = projectNames[0]
	}

	serviceName := attributes.GetAPIGroup()
	if override, set := serviceNameMapping[serviceName]; set {
		serviceName = override
	}

	resourceURL := "resourcemanager.datumapis.com/" + projectName
	permissionName := fmt.Sprintf("%s/%s.%s", serviceName, attributes.GetResource(), attributes.GetVerb())

	span.SetAttributes(
		attribute.String("resource", resourceURL),
		attribute.String("permission", permissionName),
	)

	resp, err := o.IAMClient.CheckAccess(ctx, &iampb.CheckAccessRequest{
		Resource:   resourceURL,
		Subject:    "user:" + attributes.GetUser().GetName(),
		Permission: permissionName,
	})
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(ctx, "failed to check subject access in IAM system", slog.String("error", err.Error()))
		return authorizer.DecisionNoOpinion, "", err
	}
	span.SetAttributes(attribute.Bool("allowed", resp.GetAllowed()))

	if resp.GetAllowed() {
		slog.DebugContext(ctx, "subject was granted access through IAM service")
		return authorizer.DecisionAllow, "", nil
	}

	return authorizer.DecisionDeny, "", nil
}
