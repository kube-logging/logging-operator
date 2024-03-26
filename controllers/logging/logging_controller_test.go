// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/MakeNowJust/heredoc"
	"github.com/andreyvit/diff"
	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/onsi/gomega"
	"github.com/pborman/uuid"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	controllers "github.com/kube-logging/logging-operator/controllers/logging"
	"github.com/kube-logging/logging-operator/pkg/resources/fluentd"
	"github.com/kube-logging/logging-operator/pkg/resources/model"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

var (
	err error
	mgr ctrl.Manager
)

const (
	timeout = 5 * time.Second
)

func TestFluentdResourcesCreatedAndRemoved(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	defer ensureCreated(t, logging)()

	cm := &corev1.Secret{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.SecretConfigName), cm)()

	g.Expect(cm.Data["fluent.conf"]).Should(gomega.And(
		gomega.ContainSubstring("@include /fluentd/etc/input.conf"),
		gomega.ContainSubstring("@include /fluentd/app-config/*"),
		gomega.ContainSubstring("@include /fluentd/etc/devnull.conf"),
	))

	deployment := &appsv1.StatefulSet{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.StatefulSetName), deployment)()
}

func TestSingleFlowWithoutOutputRefs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestSingleFlowWithoutExistingLoggingRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			LoggingRef: "nonexistent",
			Selectors: map[string]string{
				"a": "b",
			},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).ShouldNot(gomega.ContainSubstring("namespace " + testNamespace))
}

func TestSingleFlowWithOutputRefDefaultLoggingRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			LocalOutputRefs: []string{"test-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestSingleFlowWithClusterOutput(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-cluster-output",
			Namespace: controlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: output.NewNullOutputConfig(),
			},
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			GlobalOutputRefs: []string{"test-cluster-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestLogginResourcesWithNonUniqueLoggingRefs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging1 := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-1",
		},
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: controlNamespace,
			WatchNamespaces:  []string{"a", "b"},
		},
	}
	logging2 := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-2",
		},
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: controlNamespace,
			WatchNamespaces:  []string{"b", "c"},
		},
	}
	logging3 := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-3",
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       "test",
			ControlNamespace: controlNamespace,
		},
	}

	defer ensureCreated(t, logging1)()
	defer ensureCreated(t, logging2)()
	defer ensureCreated(t, logging3)()

	// The first logging resource won't be populated with a warning initially since at the time of creation it is unique
	g.Eventually(getLoggingProblems(logging1)).WithPolling(time.Second).WithTimeout(5 * time.Second).Should(gomega.BeEmpty())
	g.Eventually(getLoggingProblems(logging2)).WithPolling(time.Second).WithTimeout(5 * time.Second).Should(gomega.ContainElement(gomega.ContainSubstring(model.LoggingRefConflict)))
	g.Eventually(getLoggingProblems(logging3)).WithPolling(time.Second).WithTimeout(5 * time.Second).Should(gomega.BeEmpty())

	g.Eventually(func() error {
		l := &v1beta1.Logging{}
		if err := mgr.GetClient().Get(context.TODO(), client.ObjectKeyFromObject(logging1), l); err != nil {
			return err
		}
		l.Spec.ErrorOutputRef = "trigger reconcile"
		return mgr.GetClient().Update(context.TODO(), l)
	}).WithPolling(time.Second).WithTimeout(5 * time.Second).Should(gomega.Succeed())
	g.Eventually(getLoggingProblems(logging1)).WithPolling(time.Second).WithTimeout(5 * time.Second).Should(gomega.ContainElement(gomega.ContainSubstring(model.LoggingRefConflict)))
}

func getLoggingProblems(logging *v1beta1.Logging) func() ([]string, error) {
	return func() ([]string, error) {
		l := &v1beta1.Logging{}
		err := mgr.GetClient().Get(context.TODO(), client.ObjectKeyFromObject(logging), l)
		return l.Status.Problems, err
	}
}

func TestSingleClusterFlowWithClusterOutputFromExternalNamespace(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:                        []string{testNamespace},
			FluentdSpec:                            &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled:                true,
			ControlNamespace:                       controlNamespace,
			AllowClusterResourcesFromAllNamespaces: true,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-cluster-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: output.NewNullOutputConfig(),
			},
		},
	}

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			GlobalOutputRefs: []string{"test-cluster-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestClusterFlowWithNamespacedOutput(t *testing.T) {
	errors := make(chan error)
	defer beforeEachWithError(t, errors)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: controlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			GlobalOutputRefs: []string{"test-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	expectError(t, "referenced clusteroutput not found: test-output", errors)
}

func TestSingleFlowWithOutputRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:              "someloggingref",
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			LoggingRef:       "someloggingref",
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			LoggingRef: "someloggingref",
			Selectors: map[string]string{
				"a": "b",
			},
			LocalOutputRefs: []string{"test-output"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()

	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("a:b"))
}

func TestSingleFlowDefaultLoggingRefInvalidOutputRef(t *testing.T) {
	errors := make(chan error)
	defer beforeEachWithError(t, errors)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			LocalOutputRefs: []string{"test-output-nonexistent"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, flow)()

	expected := fmt.Sprintf("referenced output test-output-nonexistent not found for flow %s/test-flow", testNamespace)
	expectError(t, expected, errors)
}

func TestSingleFlowWithSecretInOutput(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			S3OutputConfig: &output.S3OutputConfig{
				AwsAccessKey: &secret.Secret{
					ValueFrom: &secret.ValueFrom{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "topsecret",
							},
							Key: "key",
						},
					},
				},
				AwsSecretKey: &secret.Secret{
					MountFrom: &secret.ValueFrom{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "topsecret",
							},
							Key: "key",
						},
					},
				},
				SharedCredentials: &output.S3SharedCredentials{},
			},
		},
	}
	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			LocalOutputRefs: []string{
				"test-output",
			},
		},
	}
	topsecret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      "topsecret",
			Namespace: testNamespace,
		},
		StringData: map[string]string{
			"key": "topsecretdata",
		},
	}
	defer ensureCreated(t, logging)()
	defer ensureCreated(t, topsecret)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	secret := &corev1.Secret{}
	secretKey := fmt.Sprintf("%s-topsecret-key", testNamespace)

	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), secret)()
	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring("aws_key_id topsecretdata"))
	g.Expect(string(secret.Data[fluentd.AppConfigKey])).Should(gomega.ContainSubstring(
		fmt.Sprintf("aws_sec_key /fluentd/secret/%s", secretKey)))

	outputSecret := &corev1.Secret{}
	defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.OutputSecretName), outputSecret)()

	g.Expect(outputSecret.Data).Should(gomega.HaveKeyWithValue(secretKey, []byte("topsecretdata")))
}

func TestMultiProcessWorker(t *testing.T) {
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec: &v1beta1.FluentdSpec{
				Workers: 2,
				RootDir: "/var/log/testing-testing",
			},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}
	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			FileOutput: &output.FileOutputConfig{
				Path:   "/tmp/logs/${tag}/%Y/%m/%d.%H.%M",
				Append: true,
				Buffer: &output.Buffer{
					Timekey:       "1m",
					TimekeyWait:   "30s",
					TimekeyUseUtc: true,
					Path:          "asd",
				},
			},
		},
	}
	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Selectors: map[string]string{
				"a": "b",
			},
			LocalOutputRefs: []string{
				"test-output",
			},
		},
	}
	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	{
		var secret corev1.Secret
		defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.SecretConfigName), &secret)()

		expected := heredoc.Docf(`
		# Enable RPC endpoint (this allows to trigger config reload without restart)
		<system>
		  rpc_endpoint 127.0.0.1:24444
		  log_level info
		  workers 2
		  root_dir %s
		</system>

		# Prometheus monitoring
		`, logging.Spec.FluentdSpec.RootDir)

		if e, a := diff.TrimLinesInString(expected), diff.TrimLinesInString(string(secret.Data["input.conf"])); e != a {
			t.Errorf("input.conf does not match (-actual vs +expected):\n%s\nActual:\n%s", diff.LineDiff(a, e), a)
		}
	}

	{
		var secret corev1.Secret
		defer ensureCreatedEventually(t, controlNamespace, logging.QualifiedName(fluentd.AppSecretConfigName), &secret)()

		config := string(secret.Data[fluentd.AppConfigKey])
		t.Logf("%s:\n%s", fluentd.AppConfigKey, config)
		matches := regexp.MustCompile(`[ \t]*<buffer[^>]*>((?s:.*))</buffer>`).FindAllString(config, -1)
		pathRegexp := regexp.MustCompile(`[ \t]*path[ \t]+.*`)
		for _, match := range matches {
			if pathRegexp.MatchString(match) {
				t.Errorf("Config shouldn't contain buffer directive with path parameter. Buffer:\n%s", match)
			}
		}
	}
}

func TestClusterOutputWithoutPlugin(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("no output target configured"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.ProblemsCount == len(output.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestOutputWithoutPlugin(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("no output target configured"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.ProblemsCount == len(output.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestClusterOutputWithMultiplePlugins(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				FileOutput: &output.FileOutputConfig{
					Path: "/dev/null",
				},
				NullOutputConfig: &output.NullOutputConfig{},
			},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("multiple output targets configured: [file nullout]"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.ProblemsCount == len(output.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestOutputWithMultiplePlugins(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			FileOutput: &output.FileOutputConfig{
				Path: "/dev/null",
			},
			NullOutputConfig: &output.NullOutputConfig{},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("multiple output targets configured: [file nullout]"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.ProblemsCount == len(output.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestClusterOutputWithMissingSecret(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				SyslogOutputConfig: &output.SyslogOutputConfig{
					Host: "localhost",
					TrustedCaPath: &secret.Secret{
						ValueFrom: &secret.ValueFrom{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "no-such-secret",
								},
								Key: "the-value",
							},
						},
					},
				},
			},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf(gomega.ContainSubstring("Secret \"no-such-secret\" not found")))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.ProblemsCount == len(output.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestOutputWithMissingSecret(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			SyslogOutputConfig: &output.SyslogOutputConfig{
				Host: "localhost",
				TrustedCaPath: &secret.Secret{
					ValueFrom: &secret.ValueFrom{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "no-such-secret",
							},
							Key: "the-value",
						},
					},
				},
			},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf(gomega.ContainSubstring("Secret \"no-such-secret\" not found")))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(output), output)
		return output.Status.ProblemsCount == len(output.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestClusterFlowWithLegacyOutputRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: &output.NullOutputConfig{},
			},
		},
	}

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{},
				},
			},
			OutputRefs: []string{output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("\"outputRefs\" field is deprecated, use \"globalOutputRefs\" instead"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.ProblemsCount == len(flow.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestFlowWithLegacyOutputRef(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	clusterOutput := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-clusteroutput",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: &output.NullOutputConfig{},
			},
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: &output.NullOutputConfig{},
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{},
				},
			},
			OutputRefs: []string{clusterOutput.Name, output.Name},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, clusterOutput)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("\"outputRefs\" field is deprecated, use \"globalOutputRefs\" and \"localOutputRefs\" instead"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.ProblemsCount == len(flow.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestClusterFlowWithDanglingGlobalOutputRefs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	output := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: &output.NullOutputConfig{},
			},
		},
	}

	flow := &v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterFlowSpec{
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{},
				},
			},
			GlobalOutputRefs: []string{"no-such-output-1", output.Name, "no-such-output-2"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf("dangling global output reference: no-such-output-1", "dangling global output reference: no-such-output-2"))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.ProblemsCount == len(flow.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestFlowWithDanglingLocalAndGlobalOutputRefs(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	logging := &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			WatchNamespaces:         []string{testNamespace},
			ControlNamespace:        controlNamespace,
		},
	}

	clusterOutput := &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-clusteroutput",
			Namespace: logging.Spec.ControlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: &output.NullOutputConfig{},
			},
		},
	}

	output := &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: &output.NullOutputConfig{},
		},
	}

	flow := &v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-flow",
			Namespace: testNamespace,
		},
		Spec: v1beta1.FlowSpec{
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{},
				},
			},
			GlobalOutputRefs: []string{"no-such-output-1", clusterOutput.Name, "no-such-output-2"},
			LocalOutputRefs:  []string{"no-such-output-1", output.Name, "no-such-output-2"},
		},
	}

	defer ensureCreated(t, logging)()
	defer ensureCreated(t, clusterOutput)()
	defer ensureCreated(t, output)()
	defer ensureCreated(t, flow)()

	g.Eventually(func() ([]string, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.Problems, err
	}, timeout).Should(gomega.ConsistOf(
		"dangling global output reference: no-such-output-1",
		"dangling global output reference: no-such-output-2",
		"dangling local output reference: no-such-output-1",
		"dangling local output reference: no-such-output-2",
	))

	g.Eventually(func() (bool, error) {
		err := mgr.GetClient().Get(context.TODO(), utils.ObjectKeyFromObjectMeta(flow), flow)
		return flow.Status.ProblemsCount == len(flow.Status.Problems), err
	}, timeout).Should(gomega.BeTrue())
}

func TestWatchNamespaces(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	defer beforeEach(t)()

	defer ensureCreated(t, &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-bylabel-1",
			Labels: map[string]string{
				"bylabel": "test1",
			},
		},
	})()
	defer ensureCreated(t, &corev1.Namespace{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-bylabel-2",
			Labels: map[string]string{
				"bylabel": "test2",
			},
		},
	})()

	type ReturnVal struct {
		namespaces []string
		err        error
	}

	cases := []struct {
		name           string
		logging        *v1beta1.Logging
		expectedResult func() ReturnVal
		expectError    bool
	}{
		{
			name: "full list",
			logging: &v1beta1.Logging{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-" + uuid.New()[:8],
				},
				Spec: v1beta1.LoggingSpec{
					WatchNamespaces:        []string{},
					WatchNamespaceSelector: nil,
				},
			},
			expectedResult: func() ReturnVal {
				allNamespaces := &corev1.NamespaceList{}
				err := mgr.GetClient().List(context.TODO(), allNamespaces)
				if err != nil {
					t.Fatalf("unexpected error when getting namespaces %s", err)
				}
				items := []string{}
				for _, i := range allNamespaces.Items {
					items = append(items, i.Name)
				}
				slices.Sort(items)
				return ReturnVal{
					namespaces: items,
				}
			},
		},
		{
			name: "explicit list",
			logging: &v1beta1.Logging{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-" + uuid.New()[:8],
				},
				Spec: v1beta1.LoggingSpec{
					WatchNamespaces:        []string{"test-explicit-1", "test-explicit-2"},
					WatchNamespaceSelector: nil,
				},
			},
			expectedResult: func() ReturnVal {
				return ReturnVal{
					namespaces: []string{"test-explicit-1", "test-explicit-2"},
				}
			},
		},
		{
			name: "bylabel list",
			logging: &v1beta1.Logging{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-" + uuid.New()[:8],
				},
				Spec: v1beta1.LoggingSpec{
					WatchNamespaces: []string{},
					WatchNamespaceSelector: &v1.LabelSelector{
						MatchLabels: map[string]string{
							"bylabel": "test1",
						},
					},
				},
			},
			expectedResult: func() ReturnVal {
				return ReturnVal{
					namespaces: []string{"test-bylabel-1"},
				}
			},
		},
		{
			name: "bylabel negative list (label exists but value should be different)",
			logging: &v1beta1.Logging{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-" + uuid.New()[:8],
				},
				Spec: v1beta1.LoggingSpec{
					WatchNamespaces: []string{},
					WatchNamespaceSelector: &v1.LabelSelector{
						MatchExpressions: []v1.LabelSelectorRequirement{
							{
								Key:      "bylabel",
								Operator: v1.LabelSelectorOpExists,
							},
							{
								Key:      "bylabel",
								Operator: v1.LabelSelectorOpNotIn,
								Values:   []string{"test1"},
							},
						},
					},
				},
			},
			expectedResult: func() ReturnVal {
				return ReturnVal{
					namespaces: []string{"test-bylabel-2"},
				}
			},
		},
		{
			name: "merge two sets uniquely",
			logging: &v1beta1.Logging{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-" + uuid.New()[:8],
				},
				Spec: v1beta1.LoggingSpec{
					WatchNamespaces: []string{"a", "b", "c", "test-bylabel-1"},
					WatchNamespaceSelector: &v1.LabelSelector{
						MatchExpressions: []v1.LabelSelectorRequirement{
							{
								Key:      "bylabel",
								Operator: v1.LabelSelectorOpExists,
							},
						},
					},
				},
			},
			expectedResult: func() ReturnVal {
				return ReturnVal{
					namespaces: []string{"a", "b", "c", "test-bylabel-1", "test-bylabel-2"},
				}
			},
		},
	}

	for _, c := range cases {
		if c.expectError {
			_, err := model.UniqueWatchNamespaces(context.TODO(), mgr.GetClient(), c.logging)
			if c.expectError && err == nil {
				t.Fatalf("expected error for test case %s", c.name)
			}
			continue
		}

		g.Eventually(func() ReturnVal {
			n, e := model.UniqueWatchNamespaces(context.TODO(), mgr.GetClient(), c.logging)
			return ReturnVal{
				namespaces: n,
				err:        e,
			}
		}, timeout).Should(gomega.Equal(
			c.expectedResult(),
		))
	}
}

func beforeEach(t *testing.T) func() {
	return beforeEachWithError(t, nil)
}

func beforeEachWithError(t *testing.T, errors chan<- error) func() {
	g := gomega.NewWithT(t)

	timeout := 1 * time.Second

	mgr, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme:                  scheme.Scheme,
		Metrics:                 metricsserver.Options{BindAddress: "0"},
		GracefulShutdownTimeout: &timeout,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	flowReconciler := controllers.NewLoggingReconciler(mgr.GetClient(), mgr.GetEventRecorderFor("logging-operator"), ctrl.Log.WithName("controllers").WithName("Flow"))

	var stopped bool
	wrappedReconciler := duplicateRequest(t, flowReconciler, &stopped, errors)

	err := controllers.SetupLoggingWithManager(mgr, ctrl.Log.WithName("manager").WithName("Setup")).
		Named(uuid.New()[:8]).Complete(wrappedReconciler)

	g.Expect(err).NotTo(gomega.HaveOccurred())

	stopMgr, mgrStopped := startTestManager(t, mgr)

	return func() {
		stopMgr()
		stopped = true
		mgrStopped.Wait()
	}
}

func ensureCreated(t *testing.T, obj runtime.Object) func() {
	object, ok := obj.(client.Object)
	if !ok {
		t.Fatalf("unable to cast runtime.Object to client.Object")
	}

	if err := mgr.GetClient().Create(context.TODO(), object); err != nil {
		t.Fatalf("%+v", err)
	}
	return func() {
		err := mgr.GetClient().Delete(context.TODO(), object)
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}
	}
}

func ensureCreatedAll[O client.Object](t *testing.T, objs []O) func() {
	for _, object := range objs {
		if err := mgr.GetClient().Create(context.TODO(), object); err != nil {
			t.Fatalf("%+v", err)
		}
	}
	return func() {
		for _, object := range objs {
			err := mgr.GetClient().Delete(context.TODO(), object)
			if err != nil {
				t.Fatalf("%+v", errors.WithStack(err))
			}
		}
	}
}

func ensureCreatedEventually(t *testing.T, ns, name string, obj runtime.Object) func() {
	object, ok := obj.(client.Object)
	if !ok {
		t.Fatalf("unable to cast runtime.Object to client.Object")
	}

	err := wait.PollUntilContextTimeout(context.TODO(), time.Second, time.Second*3, false, func(ctx context.Context) (bool, error) {
		err := mgr.GetClient().Get(ctx, types.NamespacedName{
			Name: name, Namespace: ns,
		}, object)
		if apierrors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	})
	if err != nil {
		t.Fatalf("%+v", errors.WithStack(err))
	}
	return func() {
		err := mgr.GetClient().Delete(context.TODO(), object)
		if err != nil {
			t.Fatalf("%+v", errors.WithStack(err))
		}
	}
}

func expectError(t *testing.T, expected string, reconcilerErrors <-chan error) {
	err := wait.PollUntilContextTimeout(context.TODO(), time.Second, time.Second*3, false, func(ctx context.Context) (bool, error) {
		select {
		case err := <-reconcilerErrors:
			if !strings.Contains(err.Error(), expected) {
				return false, errors.Errorf("expected `%s` but received `%s`", expected, err.Error())
			} else {
				return true, nil
			}
		case <-time.After(100 * time.Millisecond):
			return false, nil
		}
	})
	if err != nil {
		t.Fatalf("%+v", err)
	}
}

func testOutput() *v1beta1.Output {
	return &v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: testNamespace,
		},
		Spec: v1beta1.OutputSpec{
			NullOutputConfig: output.NewNullOutputConfig(),
		},
	}
}

func testClusterOutput() *v1beta1.ClusterOutput {
	return &v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-output",
			Namespace: controlNamespace,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				NullOutputConfig: output.NewNullOutputConfig(),
			},
		},
	}
}

func testLogging() *v1beta1.Logging {
	return &v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-" + uuid.New()[:8],
		},
		Spec: v1beta1.LoggingSpec{
			WatchNamespaces:         []string{testNamespace},
			FluentdSpec:             &v1beta1.FluentdSpec{},
			FlowConfigCheckDisabled: true,
			ControlNamespace:        controlNamespace,
		},
	}
}
