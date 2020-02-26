module github.com/banzaicloud/logging-operator

go 1.12

require (
	emperror.dev/errors v0.7.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0
	github.com/banzaicloud/operator-tools v0.9.2
	github.com/coreos/prometheus-operator v0.34.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/gomega v1.8.1
	github.com/pborman/uuid v1.2.0
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/spf13/cast v1.3.0
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	k8s.io/api v0.17.2
	k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v12.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.5.0
)

replace (
	github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk
	k8s.io/client-go => k8s.io/client-go v0.16.4
)

//replace github.com/banzaicloud/operator-tools => ../operator-tools
