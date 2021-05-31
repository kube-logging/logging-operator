module github.com/banzaicloud/logging-operator/pkg/sdk

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/operator-tools v0.23.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/zapr v0.2.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/spf13/cast v1.3.1
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	k8s.io/api v0.20.7
	k8s.io/apiextensions-apiserver v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	sigs.k8s.io/controller-runtime v0.8.3
)

replace github.com/shurcooL/vfsgen => github.com/banzaicloud/vfsgen v0.0.0-20200203103248-c48ce8603af1
