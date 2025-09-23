// Copyright © 2025 Kube logging authors
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

package fluentbit_nodeagents_namespace

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	v1beta1 "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/common/cond"
	"github.com/kube-logging/logging-operator/e2e/common/setup"
)

var TestTempDir string

func init() {
	var ok bool
	TestTempDir, ok = os.LookupEnv("PROJECT_DIR")
	if !ok {
		TestTempDir = "../.."
	}
	TestTempDir = filepath.Join(TestTempDir, "build/_test")
	err := os.MkdirAll(TestTempDir, os.FileMode(0755))
	if err != nil {
		panic(err)
	}
}

var tags = "time"
var realTimeBuffer = &output.Buffer{
	Tags:        &tags,
	Timekey:     "1s",
	TimekeyWait: "0s",
}
var producerLabels = map[string]string{
	"my-unique-label": "log-producer",
}

// Verifies that Fluent Bit node agents are deployed to a dedicated namespace when
// logging.spec.nodeAgentNamespace is set, while the aggregator stays in the control namespace.
func TestFluentbitNodeAgentsDedicatedNamespace(t *testing.T) {
	common.Initialize(t)

	nsControl := "logging"
	nsAgents := "logging-node-agents"
	release := "fluentbit-nodeagents-namespace"
	tag := "tag_nodeagents"

	common.WithCluster(release, t, func(t *testing.T, c common.Cluster) {
		// Install logging-operator Helm chart with test-receiver in control namespace
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = nsControl
			options.NameOverride = release
		}))

		ctx := context.Background()

		// Ensure node agents namespace exists
		common.RequireNoError(t, c.GetClient().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: nsAgents,
			},
		}))

		// Create ClusterOutput pointing to test-receiver in control namespace
		httpOut := v1beta1.ClusterOutput{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "http",
				Namespace: nsControl,
			},
			Spec: v1beta1.ClusterOutputSpec{
				OutputSpec: v1beta1.OutputSpec{
					LoggingRef: "infra",
					HTTPOutput: &output.HTTPOutputConfig{
						Endpoint:    fmt.Sprintf("http://%s-test-receiver:8080/%s", release, tag),
						ContentType: "application/json",
						Buffer:      realTimeBuffer,
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &httpOut))

		// Match producer pods cluster-wide
		clusterFlow := v1beta1.ClusterFlow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "flow",
				Namespace: nsControl,
			},
			Spec: v1beta1.ClusterFlowSpec{
				LoggingRef: "infra",
				Match: []v1beta1.ClusterMatch{
					{
						ClusterSelect: &v1beta1.ClusterSelect{
							Labels: producerLabels,
						},
					},
				},
				GlobalOutputRefs: []string{httpOut.Name},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &clusterFlow))

		// Create Fluent Bit agent (standalone) bound to the same loggingRef
		agent := v1beta1.FluentbitAgent{
			ObjectMeta: metav1.ObjectMeta{
				Name: "infra",
			},
			Spec: v1beta1.FluentbitSpec{
				LoggingRef: "infra",
				ConfigHotReload: &v1beta1.HotReload{
					Image: v1beta1.ImageSpec{Repository: common.ConfigReloaderRepo, Tag: common.ConfigReloaderTag},
				},
				BufferVolumeImage: v1beta1.ImageSpec{Repository: common.NodeExporterRepo, Tag: common.NodeExporterTag},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &agent))

		// Create Logging with control namespace and dedicated nodeAgentNamespace
		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{Name: "infra"},
			Spec: v1beta1.LoggingSpec{
				LoggingRef:         "infra",
				ControlNamespace:   nsControl,
				NodeAgentNamespace: nsAgents,
				FluentdSpec: &v1beta1.FluentdSpec{
					Image:               v1beta1.ImageSpec{Repository: common.FluentdImageRepo, Tag: common.FluentdImageTag},
					ConfigReloaderImage: v1beta1.ImageSpec{Repository: common.ConfigReloaderRepo, Tag: common.ConfigReloaderTag},
					BufferVolumeImage:   v1beta1.ImageSpec{Repository: common.NodeExporterRepo, Tag: common.NodeExporterTag},
					DisablePvc:          true,
					Resources: corev1.ResourceRequirements{Requests: corev1.ResourceList{
						corev1.ResourceCPU:    apiresource.MustParse("50m"),
						corev1.ResourceMemory: apiresource.MustParse("50M"),
					}},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))

		// Start a log producer in the default namespace
		go setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = "default"
			options.Labels = producerLabels
		}))

		fluentBitLabels := map[string]string{
			"app.kubernetes.io/name": "fluentbit",
		}
		aggregatorLabels := map[string]string{
			"app.kubernetes.io/name":      "fluentd",
			"app.kubernetes.io/component": "fluentd",
		}
		operatorLabels := map[string]string{
			"app.kubernetes.io/name": release,
		}

		require.Eventually(t, func() bool {
			if operatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(operatorLabels))(); !operatorRunning {
				t.Log("waiting for the operator")
				return false
			}
			if producerRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(producerLabels))(); !producerRunning {
				t.Log("waiting for the producer")
				return false
			}
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggregatorLabels), client.InNamespace(nsControl)); !aggregatorRunning() {
				t.Log("waiting for the aggregator in control namespace")
				return false
			}
			if fluentbitRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(fluentBitLabels), client.InNamespace(nsAgents)); !fluentbitRunning() {
				t.Log("waiting for the fluentbit daemonset in node agents namespace")
				return false
			}

			cmd := common.CmdEnv(exec.Command("kubectl",
				"logs",
				"-n", nsControl,
				"--tail", "30",
				"-l", fmt.Sprintf("app.kubernetes.io/name=%s-test-receiver", release)), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), tag)
		}, 5*time.Minute, 3*time.Second)

	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		err := c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{nsControl, nsAgents, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
		if err != nil {
			return err
		}

		loggingOperatorName := "logging-operator-" + release
		t.Logf("Collecting coverage files from logging-operator: %s/%s", nsControl, loggingOperatorName)
		err = c.CollectTestCoverageFiles(nsControl, loggingOperatorName)
		if err != nil {
			t.Logf("Failed collecting coverage files: %s", err)
		}
		return nil
	}, func(o *cluster.Options) {
		if o.Scheme == nil {
			o.Scheme = runtime.NewScheme()
		}
		common.RequireNoError(t, v1beta1.AddToScheme(o.Scheme))
		common.RequireNoError(t, apiextensionsv1.AddToScheme(o.Scheme))
		common.RequireNoError(t, appsv1.AddToScheme(o.Scheme))
		common.RequireNoError(t, batchv1.AddToScheme(o.Scheme))
		common.RequireNoError(t, corev1.AddToScheme(o.Scheme))
		common.RequireNoError(t, rbacv1.AddToScheme(o.Scheme))
	})
}
