// SPDX-License-Identifier: AGPL-3.0-only

package resourcemanager

import (
	"context"
	"encoding/hex"
	"fmt"
	"hash/fnv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	iamv1alpha1 "go.miloapis.com/milo/pkg/apis/iam/v1alpha1"
	resourcemanagerv1alpha1 "go.miloapis.com/milo/pkg/apis/resourcemanager/v1alpha1"
)

type PersonalOrganizationControllerConfig struct {
	RoleName      string `json:"roleName"`
	RoleNamespace string `json:"roleNamespace"`
}

// PersonalOrganizationController reconciles a User object
type PersonalOrganizationController struct {
	Client client.Client

	Config PersonalOrganizationControllerConfig
}

// +kubebuilder:rbac:groups=iam.datumapis.com,resources=users,verbs=get;list;watch
// +kubebuilder:rbac:groups=resourcemanager.datumapis.com,resources=organizations,verbs=create

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
			Annotations: map[string]string{
				"kubernetes.io/display-name": fmt.Sprintf("Personal Organization - %s %s", user.Spec.GivenName, user.Spec.FamilyName),
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, personalOrg, func() error {
		logger.Info("Creating or updating personal organization", "organization", personalOrg.Name)
		personalOrg.Annotations["kubernetes.io/display-name"] = fmt.Sprintf("%s %s's Personal Org", user.Spec.GivenName, user.Spec.FamilyName)
		personalOrg.ObjectMeta.OwnerReferences = []metav1.OwnerReference{
			// The owner reference is used to ensure that the personal organization
			// is deleted when the user is deleted.
			{
				APIVersion: iamv1alpha1.SchemeGroupVersion.String(),
				Kind:       "User",
				Name:       user.Name,
				UID:        user.UID,
			},
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
		}
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create or update organization membership: %w", err)
	}

	// Assign the default role to the user.
	policyBinding := &iamv1alpha1.PolicyBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("policy-binding-%s", user.Name),
			Namespace: fmt.Sprintf("organization-%s", personalOrg.Name),
		},
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, policyBinding, func() error {
		logger.Info("Creating or updating personal organization policy binding", "organization", personalOrg.Name)
		policyBinding.Spec = iamv1alpha1.PolicyBindingSpec{
			RoleRef: iamv1alpha1.RoleReference{
				Name:      r.Config.RoleName,
				Namespace: r.Config.RoleNamespace,
			},
			ResourceSelector: iamv1alpha1.ResourceSelector{
				ResourceRef: &iamv1alpha1.ResourceReference{
					APIGroup: resourcemanagerv1alpha1.GroupVersion.Group,
					Kind:     "Organization",
					Name:     personalOrg.Name,
					UID:      string(personalOrg.UID),
				},
			},
			Subjects: []iamv1alpha1.Subject{
				{
					Kind: "User",
					Name: user.Name,
					UID:  string(user.UID),
				},
			},
		}

		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to create or update policy binding: %w", err)
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
