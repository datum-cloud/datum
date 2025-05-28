package v1alpha1

// ScopedRoleReference defines a reference to another Role, scoped by namespace.
// This is used for purposes like role inheritance where a simple name and namespace
// is sufficient to identify the target role.
// +k8s:deepcopy-gen=true
// +kubebuilder:object:generate=true
type ScopedRoleReference struct {
	// Name of the referenced Role.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Namespace of the referenced Role.
	// If not specified, it defaults to the namespace of the resource containing this reference.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty"`
}

// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster

// IAMSetting is the Schema for the iamsettings API
