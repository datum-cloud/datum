package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ProtectedResourceSpec defines the desired state of ProtectedResource
// +k8s:openapi-gen=true
type ProtectedResourceSpec struct {
	// ServiceRef references the service definition this protected resource belongs to.
	// +kubebuilder:validation:Required
	ServiceRef ServiceReference `json:"serviceRef"`

	// The kind of the resource.
	// This will be in the format `Workload`.
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`

	// The singular form for the resource type, e.g. 'workload'. Must follow
	// camelCase format.
	// +kubebuilder:validation:Required
	Singular string `json:"singular"`

	// The plural form for the resource type, e.g. 'workloads'. Must follow
	// camelCase format.
	// +kubebuilder:validation:Required
	Plural string `json:"plural"`

	// A list of resources that are registered with the platform that may be a
	// parent to the resource. Permissions may be bound to a parent resource so
	// they can be inherited down the resource hierarchy.
	// +kubebuilder:validation:Optional
	ParentResources []ParentResourceRef `json:"parentResources,omitempty"`

	// A list of permissions that are associated with the resource.
	// +kubebuilder:validation:Required
	Permissions []string `json:"permissions"`
}

// ProtectedResourceStatus defines the observed state of ProtectedResource
type ProtectedResourceStatus struct {
	// Conditions provide conditions that represent the current status of the ProtectedResource.
	// +kubebuilder:default={{type: "Ready", status: "Unknown", reason: "Unknown", message: "Waiting for control plane to reconcile", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ObservedGeneration is the most recent generation observed for this ProtectedResource. It corresponds to the
	// ProtectedResource's generation, which is updated on mutation by the API Server.
	// +kubebuilder:validation:Optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProtectedResource is the Schema for the protectedresources API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Kind",type="string",JSONPath=".spec.kind"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=protectedresources,scope=Cluster,singular=protectedresource
type ProtectedResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProtectedResourceSpec   `json:"spec,omitempty"`
	Status ProtectedResourceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProtectedResourceList contains a list of ProtectedResource
type ProtectedResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProtectedResource `json:"items"`
}

// ParentResourceRef defines the reference to a parent resource
// +k8s:openapi-gen=true
type ParentResourceRef struct {
	// APIGroup is the group for the resource being referenced.
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	// +kubebuilder:validation:Optional
	APIGroup string `json:"apiGroup,omitempty"`
	// Kind is the type of resource being referenced.
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
}

// ServiceReference holds a reference to a service definition.
// +k8s:openapi-gen=true
type ServiceReference struct {
	// Name is the resource name of the service definition.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}
