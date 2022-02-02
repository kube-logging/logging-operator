module github.com/banzaicloud/logging-operator/e2e

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/banzaicloud/logging-operator/pkg/sdk v0.7.4
	github.com/banzaicloud/operator-tools v0.26.0
	github.com/go-logr/logr v0.4.0
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.3
	sigs.k8s.io/controller-runtime v0.9.5
)

replace github.com/banzaicloud/logging-operator/pkg/sdk => ../pkg/sdk
