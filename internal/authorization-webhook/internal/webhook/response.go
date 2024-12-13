// This file is inspired by https://github.com/kubernetes-sigs/controller-runtime/blob/v0.16.3/pkg/webhook/authentication/response.go
package webhook

import (
	authorizationv1 "k8s.io/api/authorization/v1"
)

// Allowed constructs a response indicating that the given user is authorized
// to perform the given action. The reason parameter is optional.
func Allowed(reason string) Response {
	return AuthorizationResponse(true, false, reason, "")
}

// NoOpinion constructs a response indicating that the authorizer has no rule
// allowing the request, but also no rule specifically denying the request
func NoOpinion() Response {
	return AuthorizationResponse(false, false, "", "")
}

// Denied constructs a response indicating that the given user is denied
// to perform the given action. The reason parameter is optional.
func Denied(reason string) Response {
	return AuthorizationResponse(false, true, reason, "")
}

// Errored creates a new Response for error-handling a request.
func Errored(err error) Response {
	return AuthorizationResponse(false, false, "", err.Error())
}

// AuthorizationResponse returns a response an authorization request.
func AuthorizationResponse(allowed, denied bool, reason, evaluationError string) Response {
	return Response{
		authorizationv1.SubjectAccessReview{
			Status: authorizationv1.SubjectAccessReviewStatus{
				Allowed:         allowed,
				Denied:          denied,
				Reason:          reason,
				EvaluationError: evaluationError,
			},
		},
	}
}
