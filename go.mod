module github.com/banzaicloud/logging-operator

go 1.14

require (
	emperror.dev/errors v0.7.0
	github.com/Azure/go-autorest v11.1.2+incompatible // indirect
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/ant31/crd-validation v0.0.0-20180702145049-30f8a35d0ac2 // indirect
	github.com/banzaicloud/logging-operator/pkg/sdk v0.0.0
	github.com/banzaicloud/operator-tools v0.12.0
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/brancz/gojsontoyaml v0.0.0-20190425155809-e8bd32d46b3d // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/coreos/prometheus-operator v0.29.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/go-logr/logr v0.1.0
	github.com/hashicorp/go-version v1.1.0 // indirect
	github.com/improbable-eng/thanos v0.3.2 // indirect
	github.com/jsonnet-bundler/jsonnet-bundler v0.1.0 // indirect
	github.com/kylelemons/godebug v0.0.0-20170820004349-d65d576e9348 // indirect
	github.com/mitchellh/hashstructure v0.0.0-20170609045927-2bca23e0e452 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/onsi/gomega v1.10.1
	github.com/openshift/prom-label-proxy v0.1.1-0.20191016113035-b8153a7f39f1 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/prometheus/client_golang v1.1.0 // indirect
	github.com/prometheus/tsdb v0.8.0 // indirect
	github.com/spf13/cast v1.3.1
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0 // indirect
	k8s.io/api v0.18.6
	k8s.io/apiextensions-apiserver v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
)

replace github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk

//replace github.com/banzaicloud/operator-tools => ../operator-tools
