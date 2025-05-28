package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleReference contains information that points to the Role being used
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
type RoleReference struct {
	// Name is the name of resource being referenced
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Namespace of the referenced Role. If empty, it is assumed to be in the PolicyBinding's namespace.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// Subject contains a reference to the object or user identities a role binding applies to.
// This can be a User or Group.
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
type Subject struct {
	// Kind of object being referenced. Values defined in Kind constants.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=User
	Kind string `json:"kind"`
	// Name of the object being referenced.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// Namespace of the referenced object. If DNE, then for an SA it refers to the PolicyBinding resource's namespace.
	// For a User or Group, it is ignored.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
	// UID of the referenced object.
	// +kubebuilder:validation:Required
	UID string `json:"uid"`
}

// TargetReference contains enough information to let you identify an API resource.
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen=true
type TargetReference struct {
	// APIGroup is the group for the resource being referenced.
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	// +kubebuilder:validation:Optional
	APIGroup string `json:"apiGroup,omitempty"`
	// Kind is the type of resource being referenced.
	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
	// Name is the name of resource being referenced.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// UID is the unique identifier of the resource being referenced.
	// +kubebuilder:validation:Required
	UID string `json:"uid"`
	// Namespace is the namespace of resource being referenced.
	// Required for namespace-scoped resources. Omitted for cluster-scoped resources.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PolicyBinding is the Schema for the policybindings API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Role",type="string",JSONPath=".spec.roleRef.name"
// +kubebuilder:printcolumn:name="Target API Group",type="string",JSONPath=".spec.targetRef.apiGroup"
// +kubebuilder:printcolumn:name="Target Kind",type="string",JSONPath=".spec.targetRef.kind"
// +kubebuilder:printcolumn:name="Target Name",type="string",JSONPath=".spec.targetRef.name"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:path=policybindings,scope=Namespaced
type PolicyBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   PolicyBindingSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status PolicyBindingStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PolicyBindingSpec defines the desired state of PolicyBinding
// +k8s:deepcopy-gen=true
// +k8s:openapi-gen=true
type PolicyBindingSpec struct {
	// RoleRef is a reference to the Role that is being bound.
	// This can be a reference to a Role custom resource.
	// +kubebuilder:validation:Required
	RoleRef RoleReference `json:"roleRef"`

	// Subjects holds references to the objects the role applies to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Subjects []Subject `json:"subjects"`

	// TargetRef is a reference to the resource to which this policy binding applies.
	// This allows the binding to be about a resource in any namespace or a cluster-scoped resource.
	// +kubebuilder:validation:Required
	TargetRef TargetReference `json:"targetRef"`
}

// PolicyBindingStatus defines the observed state of PolicyBinding
// +k8s:deepcopy-gen=true
// +k8s:openapi-gen=true
type PolicyBindingStatus struct {
	// ObservedGeneration is the most recent generation observed for this PolicyBinding by the controller.
	// +kubebuilder:validation:Optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions provide conditions that represent the current status of the PolicyBinding.
	// +kubebuilder:default={{type: "Ready", status: "Unknown", reason: "Unknown", message: "Waiting for control plane to reconcile", lastTransitionTime: "1970-01-01T00:00:00Z"}}
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// PolicyBindingList contains a list of PolicyBinding
type PolicyBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PolicyBinding `json:"items"`
}
