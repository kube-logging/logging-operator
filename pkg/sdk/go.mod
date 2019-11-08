module github.com/banzaicloud/logging-operator/pkg/sdk

go 1.13

require (
	emperror.dev/errors v0.4.2
	github.com/MakeNowJust/heredoc v0.0.0-20171113091838-e9091a26100e
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/ghodss/yaml v1.0.0
	github.com/go-logr/zapr v0.1.1
	github.com/goph/emperror v0.17.2
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/sergi/go-diff v1.0.0 // indirect
	go.uber.org/zap v1.12.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b
	k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.2.0
)
