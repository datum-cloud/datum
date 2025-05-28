package resourcemanager

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	ctrladmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	iamv1alpha1 "go.datum.net/datum/pkg/apis/iam.datumapis.com/v1alpha1"
)

var setuplog = logf.Log.WithName("resourcemanager-webhook-setup")

// SetupWebhooksWithManager sets up all resourcemanager.datumapis.com webhooks
func SetupWebhooksWithManager(kubeConfig *rest.Config, mux *mux.PathRecorderMux, scheme *runtime.Scheme, systemNamespace string, organizationOwnerRoleName string, projectOwnerRoleName string) error {
	setuplog.Info("Setting up resourcemanager.datumapis.com webhooks")

	// Ensure IAM types are known to the scheme used by webhooks
	utilruntime.Must(iamv1alpha1.AddToScheme(scheme))

	decoder := ctrladmission.NewDecoder(scheme)

	// Create a new dynamic client using client-go
	dynamicClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client for resourcemanager webhooks: %w", err)
	}

	// Setup Project Mutator (runs first to add owner references)
	if err := setupProjectMutator(dynamicClient, decoder, mux); err != nil {
		return fmt.Errorf("failed to setup project mutator: %w", err)
	}

	// Setup Project Validator
	if err := setupProjectValidator(dynamicClient, decoder, mux, systemNamespace, projectOwnerRoleName); err != nil {
		return fmt.Errorf("failed to setup project validator: %w", err)
	}

	// Setup Organization Validator
	if err := setupOrganizationValidator(dynamicClient, decoder, mux, systemNamespace, organizationOwnerRoleName); err != nil {
		return fmt.Errorf("failed to setup organization validator: %w", err)
	}

	setuplog.Info("All resourcemanager.datumapis.com webhooks setup complete")
	return nil
}

// setupProjectMutator configures and registers the Project mutation webhook
func setupProjectMutator(dynamicClient dynamic.Interface, decoder admission.Decoder, mux *mux.PathRecorderMux) error {
	projectMutator := &ProjectMutator{
		Client: dynamicClient,
	}

	if err := projectMutator.InjectDecoder(decoder); err != nil {
		return fmt.Errorf("failed to inject decoder into ProjectMutator: %w", err)
	}

	mux.Handle("/webhooks/resourcemanager/mutate-v1alpha1-project", &ctrlwebhook.Admission{Handler: projectMutator})
	setuplog.Info("Registered project mutation webhook", "path", "/webhooks/resourcemanager/mutate-v1alpha1-project")

	return nil
}

// setupProjectValidator configures and registers the Project validation webhook
func setupProjectValidator(dynamicClient dynamic.Interface, decoder admission.Decoder, mux *mux.PathRecorderMux, systemNamespace string, projectOwnerRoleName string) error {
	projectValidator := &ProjectValidator{
		Client:               dynamicClient,
		SystemNamespace:      systemNamespace,
		ProjectOwnerRoleName: projectOwnerRoleName,
	}

	if err := projectValidator.InjectDecoder(decoder); err != nil {
		return fmt.Errorf("failed to inject decoder into ProjectValidator: %w", err)
	}

	mux.Handle("/webhooks/resourcemanager/validate-v1alpha1-project", &ctrlwebhook.Admission{Handler: projectValidator})
	setuplog.Info("Registered project validation webhook", "path", "/webhooks/resourcemanager/validate-v1alpha1-project")

	return nil
}

// setupOrganizationValidator configures and registers the Organization validation webhook
func setupOrganizationValidator(dynamicClient dynamic.Interface, decoder admission.Decoder, mux *mux.PathRecorderMux, systemNamespace string, organizationOwnerRoleName string) error {
	organizationValidator := &OrganizationValidator{
		Client:                    dynamicClient,
		SystemNamespace:           systemNamespace,
		OrganizationOwnerRoleName: organizationOwnerRoleName,
	}

	if err := organizationValidator.InjectDecoder(decoder); err != nil {
		return fmt.Errorf("failed to inject decoder into OrganizationValidator: %w", err)
	}

	mux.Handle("/webhooks/resourcemanager/validate-v1alpha1-organization", &ctrlwebhook.Admission{Handler: organizationValidator})
	setuplog.Info("Registered organization validation webhook", "path", "/webhooks/resourcemanager/validate-v1alpha1-organization")

	return nil
}
