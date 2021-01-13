module github.com/banzaicloud/logging-operator/pkg/sdk

go 1.13

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/operator-tools v0.16.2
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/zapr v0.2.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/spf13/cast v1.3.1
	go.uber.org/zap v1.13.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	k8s.io/api v0.19.2
	k8s.io/apiextensions-apiserver v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.6.2
)

replace github.com/shurcooL/vfsgen => github.com/banzaicloud/vfsgen v0.0.0-20200203103248-c48ce8603af1

replace github.com/banzaicloud/operator-tools => /Users/ahma/Projects/bc/operator-tools