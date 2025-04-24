package app

import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission/plugin/namespace/lifecycle"
	mutatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/mutating"
	validatingwebhook "k8s.io/apiserver/pkg/admission/plugin/webhook/validating"
	"k8s.io/kubernetes/pkg/kubeapiserver/options"
)

// DefaultOffAdmissionPlugins get admission plugins off by default for kube-apiserver.
func DefaultOffAdmissionPlugins() sets.Set[string] {
	defaultOnPlugins := sets.New[string](
		lifecycle.PluginName, // NamespaceLifecycle
		// defaulttolerationseconds.PluginName, // DefaultTolerationSeconds
		mutatingwebhook.PluginName,   // MutatingAdmissionWebhook
		validatingwebhook.PluginName, // ValidatingAdmissionWebhook
	// resourcequota.PluginName,            // ResourceQuota
	// certapproval.PluginName,              // CertificateApproval
	// certsigning.PluginName,               // CertificateSigning
	// ctbattest.PluginName,                 // ClusterTrustBundleAttest
	// certsubjectrestriction.PluginName,    // CertificateSubjectRestriction
	// validatingadmissionpolicy.PluginName, // ValidatingAdmissionPolicy, only active when feature gate ValidatingAdmissionPolicy is enabled
	)

	return sets.New(options.AllOrderedPlugins...).Difference(defaultOnPlugins)
}
