package resourcemanagerdatumapiscom

import (
	"context"
	"encoding/json"
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
var projectlog = logf.Log.WithName("project-resource")

// +kubebuilder:webhook:path=/webhooks/resourcemanager/validate-v1alpha1-project,mutating=false,failurePolicy=fail,sideEffects=None,groups=resourcemanager.datumapis.com,resources=projects,verbs=create;update,versions=v1alpha1,name=vproject.datum.net,admissionReviewVersions={v1,v1beta1}

// ProjectValidator validates Projects and creates associated PolicyBindings for owners.
type ProjectValidator struct {
	Client               dynamic.Interface
	decoder              admission.Decoder
	SystemNamespace      string
	ProjectOwnerRoleName string
}

// +kubebuilder:webhook:path=/webhooks/resourcemanager/mutate-v1alpha1-project,mutating=true,failurePolicy=fail,sideEffects=None,groups=resourcemanager.datumapis.com,resources=projects,verbs=create;update,versions=v1alpha1,name=mproject.datum.net,admissionReviewVersions={v1,v1beta1}

// ProjectMutator mutates Projects to add owner references.
type ProjectMutator struct {
	Client  dynamic.Interface
	decoder admission.Decoder
}

// Handle validates the incoming Project create/update request and creates PolicyBinding on create.
func (v *ProjectValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	projectlog.Info("Validating Project", "name", req.Name, "namespace", req.Namespace)

	project := &v1alpha1.Project{}
	if err := v.decoder.Decode(req, project); err != nil {
		projectlog.Error(err, "Failed to decode Project", "name", req.Name)
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Validate project structure and parent reference
	if resp := v.validateProject(project); !resp.Allowed {
		return resp
	}

	// Create PolicyBinding on Project Create operation
	if req.Operation == admissionv1.Create {
		if resp := v.createOwnerPolicyBinding(ctx, req, project); !resp.Allowed {
			return resp
		}
	}

	projectlog.Info("Project validation successful", "name", project.Name)
	return admission.Allowed("")
}

// validateProject performs basic validation of the project structure
func (v *ProjectValidator) validateProject(project *v1alpha1.Project) admission.Response {
	if project.Spec.Parent == nil {
		err := fmt.Errorf("project %s must have a parent reference", project.Name)
		projectlog.Error(err, "Validation failed for project", "name", project.Name)
		return admission.Denied(err.Error())
	}

	return v.validateParentReference(project)
}

// validateParentReference validates the parent reference
func (v *ProjectValidator) validateParentReference(project *v1alpha1.Project) admission.Response {
	if project.Spec.Parent.ResourceRef.Kind != "Organization" {
		err := fmt.Errorf("project parent reference kind must be 'Organization', got '%s'", project.Spec.Parent.ResourceRef.Kind)
		projectlog.Error(err, "Validation failed for project", "name", project.Name)
		return admission.Denied(err.Error())
	}

	if project.Spec.Parent.ResourceRef.APIGroup != v1alpha1.GroupVersion.Group {
		err := fmt.Errorf("parent reference apiGroup must be '%s', got '%s'", v1alpha1.GroupVersion.Group, project.Spec.Parent.ResourceRef.APIGroup)
		projectlog.Error(err, "Validation failed for project", "name", project.Name)
		return admission.Denied(err.Error())
	}

	return admission.Allowed("")
}

// lookupUser retrieves the User resource from the iam.datumapis.com API
func (v *ProjectValidator) lookupUser(ctx context.Context, username string) (*unstructured.Unstructured, admission.Response) {
	// TODO: Determine if we can actually use the UID from the User object in
	//       the UserInfo of the request. Likely need to configure the OIDC
	//       authorization to map the UID from the JWT claims.
	userGVR := schema.GroupVersionResource{
		Group:    iamv1alpha1.SchemeGroupVersion.Group,
		Version:  iamv1alpha1.SchemeGroupVersion.Version,
		Resource: "users",
	}

	foundUser, err := v.Client.Resource(userGVR).Get(ctx, username, metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("failed to get user '%s' from iam.datumapis.com API: %s", username, err.Error())
		projectlog.Error(err, errMsg)
		return nil, admission.Denied(errMsg)
	}

	userUID := foundUser.GetUID()
	if userUID == "" {
		errMsg := fmt.Sprintf("user '%s' found but has no UID", username)
		projectlog.Error(fmt.Errorf(errMsg), errMsg)
		return nil, admission.Denied(errMsg)
	}

	projectlog.Info("Found user in iam.datumapis.com API", "username", username, "userUID", userUID)
	return foundUser, admission.Allowed("")
}

// createOwnerPolicyBinding creates a PolicyBinding for the project owner
func (v *ProjectValidator) createOwnerPolicyBinding(ctx context.Context, req admission.Request, project *v1alpha1.Project) admission.Response {
	projectlog.Info("Attempting to create PolicyBinding for new project", "project", project.Name, "user", req.UserInfo.Username)

	// Look up the user in the iam API
	foundUser, resp := v.lookupUser(ctx, req.UserInfo.Username)
	if !resp.Allowed {
		return resp
	}

	// Build the PolicyBinding
	policyBinding := v.buildPolicyBinding(req, project, foundUser)

	// Create the PolicyBinding resource
	return v.createPolicyBindingResource(ctx, policyBinding, project.Name)
}

// buildPolicyBinding constructs a PolicyBinding object
func (v *ProjectValidator) buildPolicyBinding(req admission.Request, project *v1alpha1.Project, user *unstructured.Unstructured) *iamv1alpha1.PolicyBinding {
	policyBindingName := fmt.Sprintf("%s-owner", project.Name)
	userUID := user.GetUID()

	return &iamv1alpha1.PolicyBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: policyBindingName,
			// Create the policy binding in the organization's namespace that the
			// project belongs to.
			//
			// TODO: Will need to re-consider this when the folder type can be
			//       introduced as a parent. Maybe we have an Owner field in the spec?
			Namespace: fmt.Sprintf("organization-%s", project.Spec.Parent.ResourceRef.Name),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: v1alpha1.GroupVersion.String(),
					Kind:       "Project",
					Name:       project.Name,
					UID:        project.UID,
				},
			},
		},
		Spec: iamv1alpha1.PolicyBindingSpec{
			RoleRef: iamv1alpha1.RoleReference{
				Name:      v.ProjectOwnerRoleName,
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
				Kind:     "Project",
				Name:     project.Name,
				UID:      string(project.UID),
			},
		},
	}
}

// createPolicyBindingResource creates the PolicyBinding resource in the cluster
func (v *ProjectValidator) createPolicyBindingResource(ctx context.Context, policyBinding *iamv1alpha1.PolicyBinding, projectName string) admission.Response {
	// Ensure TypeMeta is set for conversion
	policyBinding.APIVersion = iamv1alpha1.SchemeGroupVersion.String()
	policyBinding.Kind = "PolicyBinding"

	projectlog.Info("Constructed PolicyBinding",
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
		projectlog.Error(err, "Failed to convert PolicyBinding to unstructured", "name", policyBinding.Name)
		return admission.Errored(http.StatusInternalServerError, fmt.Errorf("failed to convert PolicyBinding to unstructured: %w", err))
	}
	policyBindingUnstructured := &unstructured.Unstructured{Object: unstructuredPolicyBindingMap}

	_, err = v.Client.Resource(policyBindingGVR).Namespace(policyBinding.GetNamespace()).Create(ctx, policyBindingUnstructured, metav1.CreateOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("failed to create default owner policy binding '%s' for project '%s': %s", policyBinding.Name, projectName, err.Error())
		projectlog.Error(err, errMsg)
		return admission.Denied(errMsg)
	}

	projectlog.Info("Successfully created PolicyBinding for project", "project", projectName, "policyBinding", policyBinding.Name)
	return admission.Allowed("")
}

// InjectDecoder injects the decoder.
func (v *ProjectValidator) InjectDecoder(d admission.Decoder) error {
	v.decoder = d
	return nil
}

// InjectClient injects the client.
func (v *ProjectValidator) InjectClient(c dynamic.Interface) error {
	v.Client = c
	return nil
}

// Handle mutates the incoming Project to add owner references based on parent
func (m *ProjectMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	projectlog.Info("Mutating Project", "name", req.Name, "namespace", req.Namespace)

	project := &v1alpha1.Project{}
	if err := m.decoder.Decode(req, project); err != nil {
		projectlog.Error(err, "Failed to decode Project", "name", req.Name)
		return admission.Errored(http.StatusBadRequest, err)
	}

	// Only add owner reference if parent is specified and it's an Organization
	if project.Spec.Parent != nil &&
		project.Spec.Parent.ResourceRef.Kind == "Organization" &&
		project.Spec.Parent.ResourceRef.APIGroup == v1alpha1.GroupVersion.Group {

		resp := m.addOrganizationOwnerReference(ctx, project)
		if !resp.Allowed {
			return resp
		}
	}

	// Convert the modified project back to raw bytes
	projectBytes, err := json.Marshal(project)
	if err != nil {
		projectlog.Error(err, "Failed to marshal modified project", "name", project.Name)
		return admission.Errored(http.StatusInternalServerError, err)
	}

	projectlog.Info("Project mutation successful", "name", project.Name)
	return admission.PatchResponseFromRaw(req.Object.Raw, projectBytes)
}

// addOrganizationOwnerReference adds an owner reference to the parent organization
func (m *ProjectMutator) addOrganizationOwnerReference(ctx context.Context, project *v1alpha1.Project) admission.Response {
	orgName := project.Spec.Parent.ResourceRef.Name

	// Look up the organization to get its UID
	orgGVR := schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: "organizations",
	}

	org, err := m.Client.Resource(orgGVR).Get(ctx, orgName, metav1.GetOptions{})
	if err != nil {
		errMsg := fmt.Sprintf("failed to get parent organization '%s': %s", orgName, err.Error())
		projectlog.Error(err, errMsg)
		return admission.Denied(errMsg)
	}

	orgUID := org.GetUID()
	if orgUID == "" {
		errMsg := fmt.Sprintf("parent organization '%s' found but has no UID", orgName)
		projectlog.Error(fmt.Errorf(errMsg), errMsg)
		return admission.Denied(errMsg)
	}

	// Create the owner reference
	ownerRef := metav1.OwnerReference{
		APIVersion: v1alpha1.GroupVersion.String(),
		Kind:       "Organization",
		Name:       orgName,
		UID:        orgUID,
	}

	// Check if owner reference already exists
	ownerRefExists := false
	for _, existingRef := range project.OwnerReferences {
		if existingRef.UID == orgUID {
			ownerRefExists = true
			break
		}
	}

	// Add owner reference if it doesn't exist
	if !ownerRefExists {
		project.OwnerReferences = append(project.OwnerReferences, ownerRef)
		projectlog.Info("Added organization owner reference",
			"project", project.Name,
			"organization", orgName,
			"organizationUID", orgUID)
	}

	return admission.Allowed("")
}

// InjectDecoder injects the decoder into ProjectMutator
func (m *ProjectMutator) InjectDecoder(d admission.Decoder) error {
	m.decoder = d
	return nil
}

// InjectClient injects the client into ProjectMutator
func (m *ProjectMutator) InjectClient(c dynamic.Interface) error {
	m.Client = c
	return nil
}
