// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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

package volumedrain

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
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

func TestVolumeDrain_Downscale(t *testing.T) {
	t.Parallel()
	ns := "testing-1"
	releaseNameOverride := "volumedrain"
	testTag := "test.volumedrain"
	common.WithCluster("drain", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "drainer-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				EnableRecreateWorkloadOnImmutableFieldChange: true,
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Network: &v1beta1.FluentbitNetwork{
						Keepalive: utils.BoolPointer(false),
					},
				},
				FluentdSpec: &v1beta1.FluentdSpec{
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("200M"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("50M"),
						},
					},
					BufferVolumeMetrics: &v1beta1.Metrics{},
					Scaling: &v1beta1.FluentdScaling{
						Replicas: 2,
						Drain: v1beta1.FluentdDrainConfig{
							Enabled: true,
						},
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		tags := "time"
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-output",
				Namespace: ns,
			},
			Spec: v1beta1.OutputSpec{
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    fmt.Sprintf("http://%s-test-receiver:8080/%s", releaseNameOverride, testTag),
					ContentType: "application/json",
					Buffer: &output.Buffer{
						Type:        "file",
						Tags:        &tags,
						Timekey:     "1s",
						TimekeyWait: "0s",
					},
				},
			},
		}

		producerLabels := map[string]string{
			"my-unique-label": "log-producer",
		}

		common.RequireNoError(t, c.GetClient().Create(ctx, &output))
		flow := v1beta1.Flow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-flow",
				Namespace: ns,
			},
			Spec: v1beta1.FlowSpec{
				Match: []v1beta1.Match{
					{
						Select: &v1beta1.Select{
							Labels: producerLabels,
						},
					},
				},
				LocalOutputRefs: []string{output.Name},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &flow))

		aggergatorLabels := map[string]string{
			"app.kubernetes.io/name":      "fluentd",
			"app.kubernetes.io/component": "fluentd",
		}
		operatorLabels := map[string]string{
			"app.kubernetes.io/name": releaseNameOverride,
		}

		go setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = ns
			options.Labels = producerLabels
		}))

		require.Eventually(t, func() bool {
			if operatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(operatorLabels))(); !operatorRunning {
				t.Log("waiting for the operator")
				return false
			}
			if producerRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(producerLabels))(); !producerRunning {
				t.Log("waiting for the producer")
				return false
			}
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggergatorLabels)); !aggregatorRunning() {
				t.Log("waiting for the aggregator")
				return false
			}

			cmd := common.CmdEnv(exec.Command("kubectl",
				"logs",
				"-n", ns,
				"-l", fmt.Sprintf("app.kubernetes.io/name=%s-test-receiver", releaseNameOverride)), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), testTag)
		}, 5*time.Minute, 3*time.Second)

		cmd := common.CmdEnv(exec.Command("kubectl", "scale",
			fmt.Sprintf("deployment/%s-test-receiver", releaseNameOverride),
			"-n", ns,
			"--replicas", "0"), c)
		common.RequireNoError(t, cmd.Run())

		fluentdReplicaName := logging.Name + "-fluentd-1"

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl",
				"exec",
				"-n", ns, fluentdReplicaName,
				"-c", "fluentd",
				"--", "ls", "-1", "/buffers"), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to list buffer directory: %v", err)
				return false
			}
			return strings.Count(string(rawOut), "\n") > 2
		}, 1*time.Minute, 3*time.Second)

		patch := client.MergeFrom(logging.DeepCopy())
		logging.Spec.FluentdSpec.Scaling.Replicas = 1
		common.RequireNoError(t, c.GetClient().Patch(ctx, &logging, patch))

		drainerJobName := fluentdReplicaName + "-drainer"
		require.Eventually(t, func() bool {
			var job batchv1.Job
			present := cond.ResourceShouldBePresent(t, c.GetClient(), common.Resource(&job, ns, drainerJobName))()
			return present && job.Status.Active > 0
		}, 1*time.Minute, 3*time.Second)

		require.Eventually(t, cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: fluentdReplicaName}), 30*time.Second, time.Second/2)

		cmd = common.CmdEnv(exec.Command("kubectl", "scale",
			fmt.Sprintf("deployment/%s-test-receiver", releaseNameOverride),
			"-n", ns,
			"--replicas", "1"), c)
		common.RequireNoError(t, cmd.Run())

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(batchv1.Job), ns, drainerJobName)), 3*time.Minute, 3*time.Second)

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(corev1.Pod), ns, fluentdReplicaName)), 30*time.Second, time.Second)

		pvc := common.Resource(new(corev1.PersistentVolumeClaim), ns, logging.Name+"-fluentd-buffer-"+fluentdReplicaName)
		common.RequireNoError(t, c.GetClient().Get(ctx, client.ObjectKeyFromObject(pvc), pvc))
		assert.Equal(t, "drained", pvc.GetLabels()["logging.banzaicloud.io/drain-status"])
	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		return c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{ns, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
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

func TestVolumeDrain_Downscale_DeleteVolume(t *testing.T) {
	t.Parallel()
	ns := "testing-2"
	releaseNameOverride := "volumedrain"
	testTag := "test.volumedrain"
	common.WithCluster("drain-2", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()
		aggergatorLabels := map[string]string{
			"app.kubernetes.io/name":      "fluentd",
			"app.kubernetes.io/component": "fluentd",
		}
		operatorLabels := map[string]string{
			"app.kubernetes.io/name": releaseNameOverride,
		}

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "drainer-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				EnableRecreateWorkloadOnImmutableFieldChange: true,
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Network: &v1beta1.FluentbitNetwork{
						Keepalive: utils.BoolPointer(false),
					},
				},
				FluentdSpec: &v1beta1.FluentdSpec{
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("500m"),
							corev1.ResourceMemory: resource.MustParse("200M"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceCPU:    resource.MustParse("250m"),
							corev1.ResourceMemory: resource.MustParse("50M"),
						},
					},
					BufferVolumeMetrics: &v1beta1.Metrics{},
					Scaling: &v1beta1.FluentdScaling{
						Replicas: 2,
						Drain: v1beta1.FluentdDrainConfig{
							Enabled:      true,
							DeleteVolume: true,
						},
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		tags := "time"
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-output",
				Namespace: ns,
			},
			Spec: v1beta1.OutputSpec{
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    fmt.Sprintf("http://%s-test-receiver:8080/%s", releaseNameOverride, testTag),
					ContentType: "application/json",
					Buffer: &output.Buffer{
						Type:        "file",
						Tags:        &tags,
						Timekey:     "10s",
						TimekeyWait: "0s",
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &output))

		producerLabels := map[string]string{
			"my-unique-label": "log-producer",
		}
		flow := v1beta1.Flow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-flow",
				Namespace: ns,
			},
			Spec: v1beta1.FlowSpec{
				Match: []v1beta1.Match{
					{
						Select: &v1beta1.Select{
							Labels: producerLabels,
						},
					},
				},
				LocalOutputRefs: []string{output.Name},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &flow))

		fluentdReplicaName := logging.Name + "-fluentd-1"

		go setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = ns
			options.Labels = producerLabels
		}))

		require.Eventually(t, func() bool {
			if operatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(operatorLabels))(); !operatorRunning {
				t.Log("waiting for the operator")
				return false
			}
			if producerRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(producerLabels))(); !producerRunning {
				t.Log("waiting for the producer")
				return false
			}
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggergatorLabels)); !aggregatorRunning() {
				t.Log("waiting for the aggregator")
				return false
			}

			cmd := common.CmdEnv(exec.Command("kubectl",
				"logs",
				"-n", ns,
				"-l", fmt.Sprintf("app.kubernetes.io/name=%s-test-receiver", releaseNameOverride)), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), testTag)
		}, 5*time.Minute, 2*time.Second)

		cmd := common.CmdEnv(exec.Command("kubectl", "scale",
			fmt.Sprintf("deployment/%s-test-receiver", releaseNameOverride),
			"-n", ns,
			"--replicas", "0"), c)
		common.RequireNoError(t, cmd.Run())

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl", "-n", ns, "exec", fluentdReplicaName, "-c", "fluentd", "--", "ls", "-1", "/buffers"), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to list buffer directory: %v", err)
				return false
			}
			return strings.Count(string(rawOut), "\n") > 2
		}, 3*time.Minute, 10*time.Second)

		patch := client.MergeFrom(logging.DeepCopy())
		logging.Spec.FluentdSpec.Scaling.Replicas = 1
		common.RequireNoError(t, c.GetClient().Patch(ctx, &logging, patch))

		drainerJobName := fluentdReplicaName + "-drainer"
		require.Eventually(t, func() bool {
			var job batchv1.Job
			present := cond.ResourceShouldBePresent(t, c.GetClient(), common.Resource(&job, ns, drainerJobName))()
			return present && job.Status.Active > 0
		}, 2*time.Minute, 1*time.Second)

		require.Eventually(t, cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: fluentdReplicaName}), 30*time.Second, time.Second/2)

		cmd = common.CmdEnv(exec.Command("kubectl", "scale",
			fmt.Sprintf("deployment/%s-test-receiver", releaseNameOverride),
			"-n", ns,
			"--replicas", "1"), c)
		common.RequireNoError(t, cmd.Run())

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(batchv1.Job), ns, drainerJobName)), 3*time.Minute, 3*time.Second)

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(corev1.Pod), ns, fluentdReplicaName)), 30*time.Second, time.Second/2)

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(corev1.PersistentVolumeClaim), ns, logging.Name+"-fluentd-buffer-"+fluentdReplicaName)), 30*time.Second, time.Second/2)
	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		return c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{ns, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
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
