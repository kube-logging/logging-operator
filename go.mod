module github.com/banzaicloud/logging-operator

go 1.17

require (
	emperror.dev/errors v0.8.0
	github.com/MakeNowJust/heredoc v1.0.0
	github.com/andreyvit/diff v0.0.0-20170406064948-c7f18ee00883
	github.com/banzaicloud/logging-operator/pkg/sdk v0.7.26
	github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config v0.0.0
	github.com/banzaicloud/operator-tools v0.28.2
	github.com/go-logr/logr v1.2.2
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pborman/uuid v1.2.1
	github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring v0.43.0
	github.com/prometheus/client_golang v1.11.0
	github.com/spf13/cast v1.3.1
	golang.org/x/time v0.0.0-20211116232009-f0f3c7e86c11
	k8s.io/api v0.23.4
	k8s.io/apiextensions-apiserver v0.23.1
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	k8s.io/klog/v2 v2.40.1
	sigs.k8s.io/controller-runtime v0.11.1
)

require (
	cloud.google.com/go v0.81.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.18 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.13 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/banzaicloud/k8s-objectmatcher v1.8.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/briandowns/spinner v1.12.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cppforlife/go-patch v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.3+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/zapr v1.2.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/siliconbrain/go-seqs v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/wayneashleyberry/terminal-dimensions v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/exp v0.0.0-20220706164943-b4a6d9510983 // indirect
	golang.org/x/net v0.0.0-20220114011407-0dd24b26b47d // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.23.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220114203427-a0453230fd26 // indirect
	k8s.io/utils v0.0.0-20211208161948-7d6a63dca704 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/banzaicloud/logging-operator/pkg/sdk => ./pkg/sdk

replace github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/syslogng/config => ./pkg/sdk/logging/model/syslogng/config
