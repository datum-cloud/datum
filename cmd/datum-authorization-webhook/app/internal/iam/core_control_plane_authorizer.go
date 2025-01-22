package iam

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"buf.build/gen/go/datum-cloud/iam/grpc/go/datum/iam/v1alpha/iamv1alphagrpc"
	iampb "buf.build/gen/go/datum-cloud/iam/protocolbuffers/go/datum/iam/v1alpha"
	"go.datumapis.com/datum/cmd/datum-authorization-webhook/app/internal/webhook"
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
	ctx, span := otel.Tracer("go.datum.net/k8s-authz-webhook").Start(ctx, "datum.k8s-authz-webhook.global.Authorize", trace.WithAttributes(
		attribute.String("api_group", attributes.GetAPIGroup()),
		attribute.String("resource_kind", attributes.GetResource()),
	))
	defer span.End()

	if attributes.GetAPIGroup() != "resourcemanager.datumapis.com" {
		slog.DebugContext(ctx, "No opinion on auth webhook request since API Group is not managed by webhook", slog.String("api_group", attributes.GetAPIGroup()))
		return authorizer.DecisionNoOpinion, "", nil
	}

	var organizationID string
	if orgIDs, set := attributes.GetUser().GetExtra()[webhook.OrganizationIDExtraKey]; !set {
		return authorizer.DecisionDeny, "", fmt.Errorf("extra '%s' is required by core control plane authorizer", webhook.OrganizationIDExtraKey)
	} else if len(orgIDs) > 1 {
		return authorizer.DecisionDeny, "", fmt.Errorf("extra '%s' only supports one value, but multiple were provided: %v", webhook.OrganizationIDExtraKey, orgIDs)
	} else {
		organizationID = orgIDs[0]
	}

	req := getCheckAccessRequest(attributes, organizationID)

	span.SetAttributes(
		attribute.String("subject", req.Subject),
		attribute.String("resource", req.Resource),
		attribute.String("permission", req.Permission),
	)

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

func getCheckAccessRequest(attributes authorizer.Attributes, organizationID string) *iampb.CheckAccessRequest {
	req := &iampb.CheckAccessRequest{
		Subject:    "user:" + attributes.GetUser().GetName(),
		Permission: fmt.Sprintf("%s/%s.%s", attributes.GetAPIGroup(), attributes.GetResource(), attributes.GetVerb()),
	}

	// Use the organization resource URL when acting on resource collections.
	if slices.Contains([]string{"list", "create", "watch"}, attributes.GetVerb()) {
		req.Resource = "resourcemanager.datumapis.com/organizations/" + organizationID
	} else {
		req.Resource = fmt.Sprintf("resourcemanager.datumapis.com/%s/%s", attributes.GetResource(), attributes.GetName())
		req.Context = []*iampb.CheckContext{{
			ContextType: &iampb.CheckContext_ParentRelationship{
				ParentRelationship: &iampb.ParentRelationship{
					ParentResource: "resourcemanager.datumapis.com/organizations/" + organizationID,
					ChildResource:  req.Resource,
				},
			},
		}}
	}

	return req
}
