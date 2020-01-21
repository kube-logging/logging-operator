module github.com/banzaicloud/logging-operator

go 1.12

require (
	emperror.dev/errors v0.7.0
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0
	github.com/banzaicloud/operator-tools v0.1.0
	github.com/coreos/prometheus-operator v0.34.0
	github.com/go-logr/logr v0.1.0
	github.com/go-logr/zapr v0.1.1
	github.com/onsi/gomega v1.5.0
	github.com/pborman/uuid v1.2.0
	github.com/spf13/cast v1.3.0
	go.uber.org/zap v1.13.0
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	k8s.io/api v0.16.4
	k8s.io/apiextensions-apiserver v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk
	k8s.io/client-go => k8s.io/client-go v0.16.4
)
