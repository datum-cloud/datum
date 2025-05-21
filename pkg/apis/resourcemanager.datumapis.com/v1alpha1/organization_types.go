// +kubebuilder:object:generate=true
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OrganizationSpec defines the desired state of Organization
// +k8s:protobuf=true
type OrganizationSpec struct {
	// Add custom validation here if needed for fields that replace DisplayName and Description
}

// OrganizationStatus defines the observed state of Organization
// +k8s:protobuf=true
type OrganizationStatus struct {
	// ObservedGeneration is the most recent generation observed for this Organization by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"bytes,1,opt,name=observedGeneration"`

	// Conditions represents the observations of an organization's current state.
	// Known condition types are: "Ready"
	// +kubebuilder:default=`[{"type": "Ready", "status": "Unknown", "reason": "Unknown", "message": "Waiting for control plane to reconcile"}]`
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:protobuf=true

// +kubebuilder:subresource:status
// Use lowercase for path, which influences plural name. Ensure kind is Organization.
// +kubebuilder:resource:path=organizations,scope=Cluster,shortName=org;orgs,categories=datum,singular=organization
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// Organization is the Schema for the Organizations API
// +kubebuilder:object:root=true
type Organization struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   OrganizationSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status OrganizationStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:protobuf=true

// +kubebuilder:object:root=true
// OrganizationList contains a list of Organization
type OrganizationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Organization `json:"items" protobuf:"bytes,2,rep,name=items"`
}
