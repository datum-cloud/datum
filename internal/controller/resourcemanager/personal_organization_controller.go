// SPDX-License-Identifier: AGPL-3.0-only

package resourcemanager

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
		if apierrors.IsNotFound(err) {
			logger.Info("User not found, skipping reconciliation", "user", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("failed to get user: %w", err)
	}

	if !user.DeletionTimestamp.IsZero() {
		logger.Info("User is being deleted, skipping reconciliation", "user", user.Name)
		return ctrl.Result{}, nil
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

	// If the user is not active, we should not create a personal project,
	// as the impersonated client will not have the correct permissions.
	if user.Status.RegistrationApproval != iamv1alpha1.RegistrationApprovalStateApproved {
		logger.Info("User is not active, skipping personal project creation", "user", user.Name, "state", user.Status.State)
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	// Create a default personal project in the personal organization.
	personalProjectID := hashPersonalOrgName(string(user.UID))
	personalProject := &resourcemanagerv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("personal-project-%s", personalProjectID),
		},
	}

	// Use the controller's own client (cluster-scope RBAC) to check whether
	// the project already exists. The impersonated user only has org-scoped
	// permissions and cannot GET projects at the cluster scope.
	existingProject := &resourcemanagerv1alpha1.Project{}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(personalProject), existingProject)
	if err != nil {
		if !apierrors.IsNotFound(err) {
			return ctrl.Result{}, fmt.Errorf("failed to check for existing personal project: %w", err)
		}

		// The project webhook requires parent context in UserInfo.Extra fields,
		// and also looks up the requesting user by UID to create a PolicyBinding
		// granting them ownership. We impersonate the actual user so the webhook
		// sees the correct identity and creates the right PolicyBinding.
		impersonatedConfig := rest.CopyConfig(r.RestConfig)
		impersonatedConfig.Impersonate = rest.ImpersonationConfig{
			UserName: user.Spec.Email,
			UID:      user.Name,
			Groups:   []string{"system:authenticated"},
			Extra: map[string][]string{
				"iam.miloapis.com/parent-name":      {personalOrg.Name},
				"iam.miloapis.com/parent-type":      {"Organization"},
				"iam.miloapis.com/parent-api-group": {"resourcemanager.miloapis.com"},
			},
		}

		impersonatedClient, err := client.New(impersonatedConfig, client.Options{Scheme: r.Scheme})
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create impersonated client: %w", err)
		}

		// Project does not exist — create it via the impersonated client so
		// the webhook sees the actual user's identity.
		logger.Info("Creating personal project", "organization", personalOrg.Name, "project", personalProject.Name)
		metav1.SetMetaDataAnnotation(&personalProject.ObjectMeta, "kubernetes.io/display-name", "Personal Project")
		metav1.SetMetaDataAnnotation(&personalProject.ObjectMeta, "kubernetes.io/description", fmt.Sprintf("%s %s's Personal Project", user.Spec.GivenName, user.Spec.FamilyName))

		if err := impersonatedClient.Create(ctx, personalProject); err != nil {
			if apierrors.IsAlreadyExists(err) {
				logger.Info("Personal project already exists (race)", "project", personalProject.Name)
			} else {
				logger.Error(err, "Failed to create personal project")
				return ctrl.Result{}, fmt.Errorf("failed to create personal project: %w", err)
			}
		}
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
