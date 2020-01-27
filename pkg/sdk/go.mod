module github.com/banzaicloud/logging-operator/pkg/sdk

go 1.13

require (
	emperror.dev/errors v0.7.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/operator-tools v0.4.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/zapr v0.1.1
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/spf13/cast v1.3.0
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc
	k8s.io/api v0.16.4
	k8s.io/apimachinery v0.16.4
	k8s.io/client-go v0.16.4
	sigs.k8s.io/controller-runtime v0.4.0
)
