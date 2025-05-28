package resourcemanagerdatumapiscom

import (
	"context"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	iamv1alpha1 "go.datum.net/datum/pkg/apis/iam.datumapis.com/v1alpha1"
	"go.datum.net/datum/pkg/apis/resourcemanager.datumapis.com/v1alpha1"
)

// log is for logging in this package.
var organizationlog = logf.Log.WithName("organization-resource")

// +kubebuilder:webhook:path=/webhooks/resourcemanager/validate-v1alpha1-organization,mutating=false,failurePolicy=fail,sideEffects=None,groups=resourcemanager.datumapis.com,resources=organizations,verbs=create;update,versions=v1alpha1,name=vorganization.datum.net,admissionReviewVersions={v1,v1beta1}

// OrganizationValidator validates Organizations
type OrganizationValidator struct {
	Client                    dynamic.Interface
	decoder                   admission.Decoder
	SystemNamespace           string
	OrganizationOwnerRoleName string
}

// Handle validates the incoming Organization create/update request
func (v *OrganizationValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	organizationlog.Info("Validating Organization", "name", req.Name)
	org := &v1alpha1.Organization{}

	err := v.decoder.Decode(req, org)
	if err != nil {
		organizationlog.Error(err, "Failed to decode Organization", "name", req.Name)
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Create PolicyBinding on Organization Create operation
	if req.Operation == admissionv1.Create {
		if resp := v.createOwnerPolicyBinding(ctx, req, org); !resp.Allowed {
			return resp
		}
	}

	organizationlog.Info("Organization validation successful", "name", org.Name)
	return admission.Allowed("")
}

// lookupUser retrieves the User resource from the iam.datumapis.com API
func (v *OrganizationValidator) lookupUser(ctx context.Context, username string) (*unstructured.Unstructured, admission.Response) {
	// TODO: Determine if we can actually use the UID from the User object in
	//       the UserInfo of the request. Likely need to configure the OIDC
	//       authorization to map the UID from the JWT claims.
	userGVR := schema.GroupVersionResource{
		Group:    iamv1alpha1.SchemeGroupVersion.Group,
		Version:  iamv1alpha1.SchemeGroupVersion.Version,
		Resource: "users",
	}

	// Get the user directly by name since resource name matches username
	foundUser, err := v.Client.Resource(userGVR).Get(ctx, username, metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("failed to get user '%s' from iam.datumapis.com API: %s", username, err.Error())
		organizationlog.Error(err, errMsg)
		return nil, admission.Denied(errMsg)
	}

	userUID := foundUser.GetUID()
	if userUID == "" {
		errMsg := fmt.Sprintf("user '%s' found but has no UID", username)
		organizationlog.Error(fmt.Errorf(errMsg), errMsg)
		return nil, admission.Denied(errMsg)
	}

	organizationlog.Info("Found user in iam.datumapis.com API", "username", username, "userUID", userUID)
	return foundUser, admission.Allowed("")
}

// createOwnerPolicyBinding creates a PolicyBinding for the organization owner
func (v *OrganizationValidator) createOwnerPolicyBinding(ctx context.Context, req admission.Request, org *v1alpha1.Organization) admission.Response {
	organizationlog.Info("Attempting to create PolicyBinding for new organization", "organization", org.Name, "user", req.UserInfo.Username)

	// Look up the user in the iam API
	foundUser, resp := v.lookupUser(ctx, req.UserInfo.Username)
	if !resp.Allowed {
		return resp
	}

	// Build the PolicyBinding
	policyBinding := v.buildPolicyBinding(req, org, foundUser)

	// Create the PolicyBinding resource
	return v.createPolicyBindingResource(ctx, policyBinding, org.Name)
}

// buildPolicyBinding constructs a PolicyBinding object for organization owner
func (v *OrganizationValidator) buildPolicyBinding(req admission.Request, org *v1alpha1.Organization, user *unstructured.Unstructured) *iamv1alpha1.PolicyBinding {
	policyBindingName := fmt.Sprintf("%s-owner", org.Name)
	userUID := user.GetUID()

	return &iamv1alpha1.PolicyBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyBindingName,
			Namespace: fmt.Sprintf("organization-%s", org.Name),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: v1alpha1.GroupVersion.String(),
					Kind:       "Organization",
					Name:       org.Name,
					UID:        org.UID,
				},
			},
		},
		Spec: iamv1alpha1.PolicyBindingSpec{
			RoleRef: iamv1alpha1.RoleReference{
				Name:      v.OrganizationOwnerRoleName,
				Namespace: v.SystemNamespace,
			},
			Subjects: []iamv1alpha1.Subject{
				{
					Kind: "User",
					Name: req.UserInfo.Username,
					UID:  string(userUID),
				},
			},
			TargetRef: iamv1alpha1.TargetReference{
				APIGroup: v1alpha1.GroupVersion.Group,
				Kind:     "Organization",
				Name:     org.Name,
				UID:      string(org.UID),
			},
		},
	}
}

// createPolicyBindingResource creates the PolicyBinding resource in the cluster
func (v *OrganizationValidator) createPolicyBindingResource(ctx context.Context, policyBinding *iamv1alpha1.PolicyBinding, orgName string) admission.Response {
	// Ensure TypeMeta is set for conversion
	policyBinding.APIVersion = iamv1alpha1.SchemeGroupVersion.String()
	policyBinding.Kind = "PolicyBinding"

	organizationlog.Info("Constructed PolicyBinding",
		"policyBindingName", policyBinding.Name,
		"policyBindingNamespace", policyBinding.Namespace,
		"targetKind", policyBinding.Spec.TargetRef.Kind,
		"targetName", policyBinding.Spec.TargetRef.Name,
		"targetUID", policyBinding.Spec.TargetRef.UID)

	policyBindingGVR := schema.GroupVersionResource{
		Group:    iamv1alpha1.SchemeGroupVersion.Group,
		Version:  iamv1alpha1.SchemeGroupVersion.Version,
		Resource: "policybindings",
	}

	unstructuredPolicyBindingMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policyBinding)
	if err != nil {
		organizationlog.Error(err, "Failed to convert PolicyBinding to unstructured", "name", policyBinding.Name)
		return admission.Errored(http.StatusInternalServerError, fmt.Errorf("failed to convert PolicyBinding to unstructured: %w", err))
	}
	policyBindingUnstructured := &unstructured.Unstructured{Object: unstructuredPolicyBindingMap}

	_, err = v.Client.Resource(policyBindingGVR).Namespace(policyBinding.GetNamespace()).Create(ctx, policyBindingUnstructured, metav1.CreateOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("failed to create default owner policy binding '%s' for organization '%s': %s", policyBinding.Name, orgName, err.Error())
		organizationlog.Error(err, errMsg)
		return admission.Denied(errMsg)
	}

	organizationlog.Info("Successfully created PolicyBinding for organization", "organization", orgName, "policyBinding", policyBinding.Name)
	return admission.Allowed("")
}

// InjectDecoder injects the decoder.
func (v *OrganizationValidator) InjectDecoder(d admission.Decoder) error {
	v.decoder = d
	return nil
}

// InjectClient injects the client.
func (v *OrganizationValidator) InjectClient(c dynamic.Interface) error {
	v.Client = c
	return nil
}
