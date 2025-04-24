package app

import (
	"context"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/informers"
	v1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/kubernetes/pkg/serviceaccount"
)

// clientGetter implements ServiceAccountTokenGetter using a factory function
type clientGetter struct {
	secretLister         v1listers.SecretLister
	serviceAccountLister v1listers.ServiceAccountLister
}

// genericTokenGetter returns a ServiceAccountTokenGetter that does not depend
// on pods and nodes.
func genericTokenGetter(factory informers.SharedInformerFactory) serviceaccount.ServiceAccountTokenGetter {
	return clientGetter{}
	// return clientGetter{secretLister: factory.Core().V1().Secrets().Lister(), serviceAccountLister: factory.Core().V1().ServiceAccounts().Lister()}
}

func (c clientGetter) GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	return c.serviceAccountLister.ServiceAccounts(namespace).Get(name)
}

func (c clientGetter) GetPod(namespace, name string) (*v1.Pod, error) {
	return nil, apierrors.NewNotFound(v1.Resource("pods"), name)
}

func (c clientGetter) GetSecret(namespace, name string) (*v1.Secret, error) {
	return c.secretLister.Secrets(namespace).Get(name)
}

func (c clientGetter) GetNode(name string) (*v1.Node, error) {
	return nil, apierrors.NewNotFound(v1.Resource("nodes"), name)
}

type jwtTokenGenerator struct {
	iss    string
	signer jose.Signer
}

func (j *jwtTokenGenerator) GenerateToken(ctx context.Context, claims *jwt.Claims, privateClaims interface{}) (string, error) {
	return "", nil
}
