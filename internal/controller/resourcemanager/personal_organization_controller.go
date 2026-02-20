// SPDX-License-Identifier: AGPL-3.0-only

package resourcemanager

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	iamv1alpha1 "go.miloapis.com/milo/pkg/apis/iam/v1alpha1"
	resourcemanagerv1alpha1 "go.miloapis.com/milo/pkg/apis/resourcemanager/v1alpha1"
)

type PersonalOrganizationControllerConfig struct {
	// The name of the role to use when assigning owner permissions to the user
	// this organization is being created for. This role should be used to grant
	// the default set of permissions that should be granted to the user the
	// personal organization is being created.
	RoleName string `json:"roleName"`

	// The namespace the owner role exists in that will be assigned to the user
	// the organization is being created for.
	RoleNamespace string `json:"roleNamespace"`
}

// PersonalOrganizationController reconciles a User object
type PersonalOrganizationController struct {
	Client client.Client

	Config PersonalOrganizationControllerConfig

	// The scheme is used to set the controller reference on the personal
	// organization.
	Scheme *runtime.Scheme

	// RestConfig is used to create an impersonated client for project creation.
	RestConfig *rest.Config
}

// +kubebuilder:rbac:groups=iam.datumapis.com,resources=users,verbs=get;list;watch
// +kubebuilder:rbac:groups=resourcemanager.datumapis.com,resources=organizations,verbs=create
// +kubebuilder:rbac:groups=resourcemanager.datumapis.com,resources=projects,verbs=create;get;list;watch;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.4/pkg/reconcile
func (r *PersonalOrganizationController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Get the user.
	user := &iamv1alpha1.User{}
	if err := r.Client.Get(ctx, req.NamespacedName, user); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to get user: %w", err)
	}

	// Automatically create a personal organization for the user. They should not
	// be able to modify or delete the organization.
	personalOrg := &resourcemanagerv1alpha1.Organization{
		ObjectMeta: metav1.ObjectMeta{
			// Create a unique name for the personal organization.
			Name: fmt.Sprintf("personal-org-%s", hashPersonalOrgName(string(user.UID))),
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, personalOrg, func() error {
		logger.Info("Creating or updating personal organization", "organization", personalOrg.Name)
		// TODO: Remove once portal uses the description annotation
		metav1.SetMetaDataAnnotation(&personalOrg.ObjectMeta, "kubernetes.io/display-name", fmt.Sprintf("%s %s's Personal Org", user.Spec.GivenName, user.Spec.FamilyName))
		metav1.SetMetaDataAnnotation(&personalOrg.ObjectMeta, "kubernetes.io/description", fmt.Sprintf("%s %s's Personal Org", user.Spec.GivenName, user.Spec.FamilyName))
		if err := controllerutil.SetControllerReference(user, personalOrg, r.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
		personalOrg.Spec.Type = "Personal"
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create or update personal organization: %w", err)
	}

	// Now we need to create the OrganizationMembership for the user to grant them
	// access to the personal organization.
	membership := &resourcemanagerv1alpha1.OrganizationMembership{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("membership-%s", user.Name),
			Namespace: fmt.Sprintf("organization-%s", personalOrg.Name),
		},
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, membership, func() error {
		logger.Info("Creating or updating personal organization membership", "organization", personalOrg.Name)
		membership.Spec = resourcemanagerv1alpha1.OrganizationMembershipSpec{
			OrganizationRef: resourcemanagerv1alpha1.OrganizationReference{
				Name: personalOrg.Name,
			},
			UserRef: resourcemanagerv1alpha1.MemberReference{
				Name: user.Name,
			},
			Roles: []resourcemanagerv1alpha1.RoleReference{
				{
					Name:      r.Config.RoleName,
					Namespace: r.Config.RoleNamespace,
				},
			},
		}
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create or update organization membership: %w", err)
	}

	// Create a default personal project in the personal organization.
	// An impersonated client is required so that the API server's project filter
	// and webhook receive the correct parent context (Organization) in the
	// user's extra fields for authorization and defaulting.
	personalProjectID := hashPersonalOrgName(string(user.UID))
	personalProject := &resourcemanagerv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("personal-project-%s", personalProjectID),
		},
	}

	impersonatedConfig := rest.CopyConfig(r.RestConfig)
	impersonatedConfig.Impersonate = rest.ImpersonationConfig{
		UserName: user.Name,
		Extra: map[string][]string{
			"iam.miloapis.com/parent-api-group": {resourcemanagerv1alpha1.GroupVersion.Group},
			"iam.miloapis.com/parent-type":      {"Organization"},
			"iam.miloapis.com/parent-name":      {personalOrg.Name},
		},
	}

	impersonatedClient, err := client.New(impersonatedConfig, client.Options{Scheme: r.Scheme})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create impersonated client: %w", err)
	}

	_, err = controllerutil.CreateOrUpdate(ctx, impersonatedClient, personalProject, func() error {
		logger.Info("Creating or updating personal project", "organization", personalOrg.Name, "project", personalProject.Name)
		metav1.SetMetaDataAnnotation(&personalProject.ObjectMeta, "kubernetes.io/display-name", "Personal Project")
		metav1.SetMetaDataAnnotation(&personalProject.ObjectMeta, "kubernetes.io/description", fmt.Sprintf("%s %s's Personal Project", user.Spec.GivenName, user.Spec.FamilyName))
		if err := controllerutil.SetControllerReference(user, personalProject, r.Scheme); err != nil {
			return fmt.Errorf("failed to set controller reference: %w", err)
		}
		personalProject.Spec = resourcemanagerv1alpha1.ProjectSpec{
			OwnerRef: resourcemanagerv1alpha1.OwnerReference{
				Kind: "Organization",
				Name: personalOrg.Name,
			},
		}
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create personal project: %w", err)
	}

	logger.Info("Successfully created or updated personal organization resources", "organization", personalOrg.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PersonalOrganizationController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iamv1alpha1.User{}).
		Named("personal-organization").
		Complete(r)
}

func hashPersonalOrgName(name string) string {
	hasher := fnv.New32a()
	//revive:disable-next-line:unhandled-error a
	hasher.Write([]byte(name))

	return hex.EncodeToString(hasher.Sum(nil))
}
