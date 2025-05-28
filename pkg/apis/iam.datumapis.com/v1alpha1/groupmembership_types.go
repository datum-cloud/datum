package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GroupMembershipSpec defines the desired state of GroupMembership
type GroupMembershipSpec struct {
	// UserRef is a reference to the User that is a member of the Group.
	// User is a cluster-scoped resource.
	// +kubebuilder:validation:Required
	UserRef UserReference `json:"userRef"`

	// GroupRef is a reference to the Group.
	// Group is a namespaced resource.
	// +kubebuilder:validation:Required
	GroupRef GroupReference `json:"groupRef"`
}

// UserReference contains information that points to the User being referenced.
// User is a cluster-scoped resource, so Namespace is not needed.
type UserReference struct {
	// Name is the name of the User being referenced.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// GroupReference contains information that points to the Group being referenced.
// Group is a namespaced resource.
type GroupReference struct {
	// Name is the name of the Group being referenced.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Namespace of the referenced Group.
	// +kubebuilder:validation:Required
	Namespace string `json:"namespace"`
}

// GroupMembershipStatus defines the observed state of GroupMembership
type GroupMembershipStatus struct {
	// Conditions represent the latest available observations of an object's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced
// +kubebuilder:printcolumn:name="User",type="string",JSONPath=".spec.userRef.name"
// +kubebuilder:printcolumn:name="Group",type="string",JSONPath=".spec.groupRef.name"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// GroupMembership is the Schema for the groupmemberships API
type GroupMembership struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GroupMembershipSpec   `json:"spec,omitempty"`
	Status GroupMembershipStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GroupMembershipList contains a list of GroupMembership
type GroupMembershipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GroupMembership `json:"items"`
}
