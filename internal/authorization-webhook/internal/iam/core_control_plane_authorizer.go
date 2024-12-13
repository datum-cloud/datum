package iam

import (
	"context"
	"fmt"
	"log/slog"

	"buf.build/gen/go/datum-cloud/iam/grpc/go/datum/iam/v1alpha/iamv1alphagrpc"
	iampb "buf.build/gen/go/datum-cloud/iam/protocolbuffers/go/datum/iam/v1alpha"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

var _ authorizer.Authorizer = &CoreControlPlaneAuthorizer{}

type CoreControlPlaneAuthorizer struct {
	IAMClient iamv1alphagrpc.AccessCheckClient
}

// Authorize implements authorizer.Authorizer.
func (o *CoreControlPlaneAuthorizer) Authorize(ctx context.Context, attributes authorizer.Attributes) (authorizer.Decision, string, error) {
	req := &iampb.CheckAccessRequest{
		Resource:   fmt.Sprintf("resourcemanager.datumapis.com/%s/%s", attributes.GetResource(), attributes.GetName()),
		Subject:    "user:" + attributes.GetUser().GetName(),
		Permission: attributes.GetVerb(),
	}
	ctx, span := otel.Tracer("go.datum.net/k8s-authz-webhook").Start(ctx, "datum.k8s-authz-webhook.global.Authorize", trace.WithAttributes(
		attribute.String("subject", req.Subject),
		attribute.String("resource", req.Resource),
		attribute.String("permission", req.Permission),
	))
	defer span.End()

	if attributes.GetAPIGroup() != "resourcemanager.datumapis.com" {
		slog.DebugContext(ctx, "No opinion on auth webhook request since API Group is not managed by webhook", slog.String("api_group", attributes.GetAPIGroup()))
		return authorizer.DecisionNoOpinion, "", nil
	}

	resp, err := o.IAMClient.CheckAccess(ctx, req)
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
