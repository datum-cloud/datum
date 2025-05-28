package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProjectSpec defines the desired state of Project.
type ProjectSpec struct {
	// A reference to the project's parent in the resource hierarchy.
	//
	// +kubebuilder:validation:Required
	Parent *ProjectParentReference `json:"parent,omitempty"`
}

// ProjectStatus defines the observed state of Project.
type ProjectStatus struct {
	// Represents the observations of a project's current state.
	// Known condition types are: "Ready"
	// +kubebuilder:default={{type: "Ready", status: "Unknown", reason: "Unknown", message: "Waiting for control plane to reconcile", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

const (
	// ProjectReady indicates that the project has been provisioned and is ready
	// for use.
	ProjectReady = "Ready"
)

const (
	// ProjectReadyReason indicates that the project is ready for use.
	ProjectReadyReason = "Ready"

	// ProjectProvisioningReason indicates that the project is provisioning.
	ProjectProvisioningReason = "Provisioning"

	// ProjectNameConflict indicates that the project name already exists
	ProjectNameConflictReason = "ProjectNameConflict"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Project is the Schema for the projects API.
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster
type Project struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Required
	Spec   ProjectSpec   `json:"spec,omitempty"`
	Status ProjectStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProjectList contains a list of Project.
type ProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Project `json:"items"`
}

type ProjectParentReference struct {
	// External is a reference to the parent of the project. Must be a valid
	// resource name.
	//
	// +kubebuilder:validation:Optional
	External string `json:"external,omitempty"`

	// Resource is a reference to the parent of the project. Must be a valid
	// resource.
	//
	// +kubebuilder:validation:Required
	ResourceRef ResourceReference `json:"resourceRef"`
}

type ResourceReference struct {
	// Group is the group of the resource.
	//
	// +kubebuilder:validation:Required
	APIGroup string `json:"apiGroup,omitempty"`

	// Kind is the kind of the resource.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=Organization
	Kind string `json:"kind,omitempty"`

	// Name is the name of the resource.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
}
