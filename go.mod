module github.com/banzaicloud/logging-operator

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0
	github.com/banzaicloud/operator-tools v0.28.2
	github.com/go-logr/logr v1.2.2
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pborman/uuid v1.2.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/spf13/cast v1.3.1
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11
	k8s.io/api v0.23.4
	k8s.io/apiextensions-apiserver v0.23.1
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	k8s.io/klog/v2 v2.40.1
	sigs.k8s.io/controller-runtime v0.11.1
)

replace github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk
