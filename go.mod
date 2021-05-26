module github.com/banzaicloud/logging-operator

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0
	github.com/banzaicloud/operator-tools v0.23.0
	github.com/go-logr/logr v0.4.0
	github.com/onsi/gomega v1.10.2
	github.com/pborman/uuid v1.2.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/spf13/cast v1.3.1
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
	k8s.io/api v0.20.7
	k8s.io/apiextensions-apiserver v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk
