package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceResource is an entity offered by services to provide functionality to service
// consumers. Resources can have actions registered that result in permissions
// being created.
// +k8s:openapi-gen=true
type ServiceResource struct {
	// The fully qualified name of the resource.
	// This will be in the format `compute.datumapis.com/Workload`.
	// +kubebuilder:validation:Required
	Type string `json:"type"`
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
	// they can be inherited down the resource hierarchy. The resource must use
	// the fully qualified resource name (e.g. compute.datumapis.com/Workload).
	// +kubebuilder:validation:Required
	ParentResources []string `json:"parentResources"`
	// A list of resource name patterns that may be present for the resource.
	// +kubebuilder:validation:Required
	ResourceNamePatterns []string `json:"resourceNamePatterns"`
	// A list of permissions that are associated with the resource.
	// +kubebuilder:validation:Required
	Permissions []string `json:"permissions"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Service is the Schema for the services API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=services,scope=Namespaced
type Service struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServiceSpec   `json:"spec,omitempty"`
	Status ServiceStatus `json:"status,omitempty"`
}

// ServiceSpec defines the desired state of Service
type ServiceSpec struct {
	// List of resources offered by a service.
	// +kubebuilder:validation:Required
	Resources []ServiceResource `json:"resources"`
}

// ServiceStatus defines the observed state of Service
type ServiceStatus struct {
	// Conditions provide conditions that represent the current status of the Service.
	// +kubebuilder:default=`[{"type": "Ready", "status": "Unknown", "reason": "Unknown", "message": "Waiting for control plane to reconcile"}]`
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceList contains a list of Service
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Service `json:"items"`
}
