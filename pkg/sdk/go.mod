module github.com/banzaicloud/logging-operator/pkg/sdk

go 1.13

require (
	emperror.dev/errors v0.7.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/operator-tools v0.10.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/zapr v0.1.1
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/spf13/cast v1.3.0
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9
	k8s.io/api v0.17.4
	k8s.io/apiextensions-apiserver v0.17.4
	k8s.io/apimachinery v0.17.4
	k8s.io/client-go v0.17.4
	sigs.k8s.io/controller-runtime v0.5.0
)

replace github.com/shurcooL/vfsgen => github.com/banzaicloud/vfsgen v0.0.0-20200203103248-c48ce8603af1
