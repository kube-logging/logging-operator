module github.com/banzaicloud/logging-operator

go 1.14

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0
	github.com/banzaicloud/operator-tools v0.16.2
	github.com/go-logr/logr v0.2.1
	github.com/imdario/mergo v0.3.9
	github.com/onsi/gomega v1.10.1
	github.com/pborman/uuid v1.2.0
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/spf13/cast v1.3.1
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.6.2
)

replace k8s.io/client-go => k8s.io/client-go v0.19.2

replace github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk

replace github.com/banzaicloud/operator-tools => /Users/ahma/Projects/bc/operator-tools