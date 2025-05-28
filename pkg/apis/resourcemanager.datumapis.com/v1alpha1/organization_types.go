// +kubebuilder:object:generate=true
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OrganizationSpec defines the desired state of Organization
// +k8s:protobuf=true
type OrganizationSpec struct {
	// The type of organization.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Personal;Business;Government;Research;Education;Nonprofit;Other
	Type string `json:"type"`
}

// OrganizationStatus defines the observed state of Organization
// +k8s:protobuf=true
type OrganizationStatus struct {
	// ObservedGeneration is the most recent generation observed for this Organization by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represents the observations of an organization's current state.
	// Known condition types are: "Ready"
	// +kubebuilder:default={{type: "Ready", status: "Unknown", reason: "Unknown", message: "Waiting for control plane to reconcile", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:protobuf=true

// +kubebuilder:subresource:status
// Use lowercase for path, which influences plural name. Ensure kind is Organization.
// +kubebuilder:resource:path=organizations,scope=Cluster,categories=datum,singular=organization
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
// Organization is the Schema for the Organizations API
// +kubebuilder:object:root=true
type Organization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrganizationSpec   `json:"spec,omitempty"`
	Status OrganizationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:protobuf=true

// +kubebuilder:object:root=true
// OrganizationList contains a list of Organization
type OrganizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Organization `json:"items"`
}
