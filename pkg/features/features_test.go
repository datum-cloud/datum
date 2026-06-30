package features_test

import (
	"testing"

	utilfeature "k8s.io/apiserver/pkg/util/feature"

	"go.datum.net/datum/pkg/features"

	_ "go.datum.net/datum/pkg/features"
)

func TestUnifiedOrganizationsDefaultDisabled(t *testing.T) {
	if utilfeature.DefaultFeatureGate.Enabled(features.UnifiedOrganizations) {
		t.Fatal("UnifiedOrganizations should default to false")
	}
}
