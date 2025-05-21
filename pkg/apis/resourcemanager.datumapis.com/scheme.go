package resourcemanager

import (
	"k8s.io/apimachinery/pkg/runtime"

	"go.datum.net/datum/pkg/apis/resourcemanager.datumapis.com/v1alpha1"
)

func Install(scheme *runtime.Scheme) {
	v1alpha1.AddToScheme(scheme)
}
