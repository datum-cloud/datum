package webhook

import (
	"context"
	"log/slog"
	"net/http"

	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

const OrganizationIDExtraKey = "resourcemanager.datumapis.com/organization-id"
const ProjectExtraKey = "resourcemanager.datumapis.com/project-name"

func NewAuthorizerWebhook(authzer authorizer.Authorizer) *Webhook {
	return &Webhook{
		Handler: HandlerFunc(func(ctx context.Context, r Request) Response {
			if r.Spec.ResourceAttributes != nil && r.Spec.NonResourceAttributes != nil {
				return Denied("must specify oneof resource or non-resource attributes, not both")
			}

			extra := map[string][]string{}
			for key, val := range r.Spec.Extra {
				extra[key] = []string(val)
			}

			attrs := authorizer.AttributesRecord{
				User: &user.DefaultInfo{
					Name:   r.Spec.User,
					UID:    r.Spec.UID,
					Groups: r.Spec.Groups,
					Extra:  extra,
				},
			}

			if resourceAttrs := r.Spec.ResourceAttributes; resourceAttrs != nil {
				attrs.Verb = resourceAttrs.Verb
				attrs.Namespace = resourceAttrs.Namespace
				attrs.APIGroup = resourceAttrs.Group
				attrs.APIVersion = resourceAttrs.Version
				attrs.Resource = resourceAttrs.Resource
				attrs.Subresource = resourceAttrs.Subresource
				attrs.Name = resourceAttrs.Name
				attrs.ResourceRequest = true
			}

			if nonResourceAttrs := r.Spec.NonResourceAttributes; nonResourceAttrs != nil {
				attrs.Verb = nonResourceAttrs.Verb
				attrs.Path = nonResourceAttrs.Path
				attrs.ResourceRequest = false
			}

			decision, reason, err := authzer.Authorize(
				ctx,
				attrs,
			)
			if err != nil {
				slog.ErrorContext(ctx, "error authorizing request", slog.String("error", err.Error()))
				return Errored(err)
			}

			status := authorizationv1.SubjectAccessReviewStatus{
				Reason: reason,
			}

			switch decision {
			case authorizer.DecisionAllow:
				status.Allowed = true
			case authorizer.DecisionDeny:
				status.Denied = true
			}

			return Response{
				SubjectAccessReview: authorizationv1.SubjectAccessReview{
					Status: status,
				},
			}
		}),
	}
}

// Webhook represents each individual webhook.
type Webhook struct {
	// Handler actually processes an authorization request returning whether it was authorized
	Handler Handler

	// WithContextFunc will allow you to take the http.Request.Context() and
	// add any additional information such as passing the request path or
	// headers thus allowing you to read them from within the handler
	WithContextFunc func(context.Context, *http.Request) context.Context
}

// Handle processes SubjectAccessReview.
func (wh *Webhook) Handle(ctx context.Context, req Request) Response {
	return wh.Handler.Handle(ctx, req)
}
