// Package features defines feature gates for the Datum controller manager.
package features

import (
	"k8s.io/apimachinery/pkg/util/runtime"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/component-base/featuregate"
)

const (
	// UnifiedOrganizations enables unified organization behavior in Datum.
	// When enabled, the PersonalOrganizationController is disabled and infra
	// should apply unified quota grant policies instead of legacy personal/standard
	// policies.
	//
	// owner: @datum-cloud/platform
	// alpha: v0.1.0
	UnifiedOrganizations featuregate.Feature = "UnifiedOrganizations"
)

func init() {
	runtime.Must(utilfeature.DefaultMutableFeatureGate.Add(defaultFeatureGates))
}

var defaultFeatureGates = map[featuregate.Feature]featuregate.FeatureSpec{
	UnifiedOrganizations: {
		Default:    false,
		PreRelease: featuregate.Alpha,
	},
}
