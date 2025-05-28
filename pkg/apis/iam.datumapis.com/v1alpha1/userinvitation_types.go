package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UserInvitation is the Schema for the userinvitations API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Email",type=string,JSONPath=".spec.email"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=userinvitations,scope=Namespaced
type UserInvitation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserInvitationSpec   `json:"spec,omitempty"`
	Status UserInvitationStatus `json:"status,omitempty"`
}

// UserInvitationSpec defines the desired state of UserInvitation
type UserInvitationSpec struct {
	// The email of the user being invited.
	// +kubebuilder:validation:Required
	Email string `json:"email"`
	// The first name of the user being invited.
	// +kubebuilder:validation:Optional
	GivenName string `json:"givenName,omitempty"`
	// The last name of the user being invited.
	// +kubebuilder:validation:Optional
	FamilyName string `json:"familyName,omitempty"`

	// The roles that will be assigned to the user when they accept the invitation.
	// +kubebuilder:validation:Optional
	Roles []RoleReference `json:"roles,omitempty"`
}

// UserInvitationStatus defines the observed state of UserInvitation
type UserInvitationStatus struct {
	// Conditions provide conditions that represent the current status of the UserInvitation.
	// +kubebuilder:default={{type: "Ready", status: "Unknown", reason: "Unknown", message: "Waiting for control plane to reconcile", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// UserInvitationList contains a list of UserInvitation
type UserInvitationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserInvitation `json:"items"`
}
