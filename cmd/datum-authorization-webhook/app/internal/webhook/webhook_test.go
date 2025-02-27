package webhook_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.datumapis.com/datum/cmd/datum-authorization-webhook/app/internal/webhook"
	v1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

func TestWebhook(t *testing.T) {
	mux := http.NewServeMux()
	// Register the project specific webhook mux
	mux.Handle("/project/v1alpha/projects/{project}/webhook", webhook.NewAuthorizerWebhook(authorizer.AuthorizerFunc(func(ctx context.Context, a authorizer.Attributes) (authorizer.Decision, string, error) {
		return authorizer.DecisionAllow, "", nil
	})))

	srv := &http.Server{
		Addr:    "127.0.0.1:11000",
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("failed to start webhook server: %s", err)
		}
	}()

	defer srv.Shutdown(context.Background())

	client := &http.Client{}

	// Define test cases to execute against the webhook
	testCases := []struct {
		desc               string
		url                url.URL
		request            webhook.Request
		expectedResponse   webhook.Response
		expectedStatusCode int
	}{
		{
			desc: "Validate server returns correct authorization response when extra data is empty",
			url: url.URL{
				Path: "/project/v1alpha/projects/my-personal-project/webhook",
			},
			request: webhook.Request{
				SubjectAccessReview: v1.SubjectAccessReview{
					Spec: v1.SubjectAccessReviewSpec{
						User: "user:test-user@datum.net",
						UID:  "525c3cc0-6960-4950-8999-f50af1bc050d",
						ResourceAttributes: &v1.ResourceAttributes{
							Namespace: "default",
							Verb:      "use",
							Group:     "networking.datumapis.com",
							Version:   "v1alpha",
							Resource:  "networks",
							Name:      "default",
						},
					},
				},
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: webhook.Response{
				SubjectAccessReview: v1.SubjectAccessReview{
					TypeMeta: metav1.TypeMeta{
						Kind:       "SubjectAccessReview",
						APIVersion: "authorization.k8s.io/v1",
					},
					Status: v1.SubjectAccessReviewStatus{
						Allowed: true,
						Denied:  false,
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		// Configure to what the server is listening on
		testCase.url.Host = "localhost:11000"
		testCase.url.Scheme = "http"

		t.Run(testCase.desc, func(t *testing.T) {
			body, err := json.Marshal(testCase.request)
			if err != nil {
				t.Errorf("failed to marshall authorization request: %s", err)
				return
			}

			resp, err := client.Post(testCase.url.String(), "application/json", strings.NewReader(string(body)))
			if err != nil {
				t.Errorf("failed to call webhook endpoint: %s", err)
				return
			}

			if resp.StatusCode != testCase.expectedStatusCode {
				t.Errorf("got status code %d but expected %d", resp.StatusCode, &testCase.expectedStatusCode)
				return
			}

			webhookResponse := webhook.Response{}
			decoder := json.NewDecoder(resp.Body)
			if err := decoder.Decode(&webhookResponse); err != nil {
				t.Errorf("failed to decode webhook response: %s", err)
				return
			}

			if !cmp.Equal(webhookResponse, testCase.expectedResponse) {
				t.Error("webhook response did not match expected response, see logs for details")
				t.Log("response diff: ", cmp.Diff(webhookResponse, testCase.expectedResponse))
			}
		})
	}
}
