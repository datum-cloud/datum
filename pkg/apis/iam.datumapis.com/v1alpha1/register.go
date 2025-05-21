package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "iam.datumapis.com", Version: "v1alpha1"}
	// SchemeBuilder initializes a scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return GroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return GroupVersion.WithResource(resource).GroupResource()
}

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&User{},
		&UserList{},
		&Group{},
		&GroupList{},
		&Service{},
		&ServiceList{},
		&Role{},
		&RoleList{},
		&PolicyBinding{},
		&PolicyBindingList{},
		&UserInvitation{},
		&UserInvitationList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}
