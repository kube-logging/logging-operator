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

package fluentd_aggregator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/cisco-open/operator-tools/pkg/utils"
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

func TestFluentdAggregator_MultiWorker(t *testing.T) {
	common.Initialize(t)
	ns := "testing-1"
	releaseNameOverride := "e2e"
	testTag := "test.fluentd_aggregator_multiworker"
	outputName := "test-output"
	flowName := "test-flow"
	common.WithCluster("fluentd-multiworker", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fluentd-aggregator-multiworker-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				EnableRecreateWorkloadOnImmutableFieldChange: true,
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Network: &v1beta1.FluentbitNetwork{
						Keepalive: utils.BoolPointer(false),
					},
					ConfigHotReload: &v1beta1.HotReload{
						Image: v1beta1.ImageSpec{
							Repository: common.ConfigReloaderRepo,
							Tag:        common.ConfigReloaderTag,
						},
					},
					BufferVolumeImage: v1beta1.ImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
				},
				FluentdSpec: &v1beta1.FluentdSpec{
					Image: v1beta1.ImageSpec{
						Repository: common.FluentdImageRepo,
						Tag:        common.FluentdImageTag,
					},
					ConfigReloaderImage: v1beta1.ImageSpec{
						Repository: common.ConfigReloaderRepo,
						Tag:        common.ConfigReloaderTag,
					},
					BufferVolumeImage: v1beta1.ImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
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
						Replicas: 1,
						Drain: v1beta1.FluentdDrainConfig{
							Enabled: true,
							Image: v1beta1.ImageSpec{
								Repository: common.FluentdDrainWatchRepo,
								Tag:        common.FluentdDrainWatchTag,
							},
						},
					},
					Workers: 2,
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		tags := "time"
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      outputName,
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
				Name:      flowName,
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

		aggregatorLabels := map[string]string{
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
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggregatorLabels)); !aggregatorRunning() {
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

	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		err := c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{ns, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
		if err != nil {
			return err
		}

		loggingOperatorName := "logging-operator-" + releaseNameOverride
		t.Logf("Collecting coverage files from logging-operator: %s/%s", ns, loggingOperatorName)
		err = c.CollectTestCoverageFiles(ns, loggingOperatorName)
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
func TestFluentdAggregator_ConfigChecks(t *testing.T) {
	common.Initialize(t)
	ns := "testing-2"
	releaseNameOverride := "e2e"
	outputName := "test-output"
	flowName := "test-flow"
	common.WithCluster("fluentd-configcheck", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fluentd-aggregator-configchecks-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				EnableRecreateWorkloadOnImmutableFieldChange: true,
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Network: &v1beta1.FluentbitNetwork{
						Keepalive: utils.BoolPointer(false),
					},
					ConfigHotReload: &v1beta1.HotReload{
						Image: v1beta1.ImageSpec{
							Repository: common.ConfigReloaderRepo,
							Tag:        common.ConfigReloaderTag,
						},
					},
					BufferVolumeImage: v1beta1.ImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
				},
				FluentdSpec: &v1beta1.FluentdSpec{
					Image: v1beta1.ImageSpec{
						Repository: common.FluentdImageRepo,
						Tag:        common.FluentdImageTag,
					},
					ConfigReloaderImage: v1beta1.ImageSpec{
						Repository: common.ConfigReloaderRepo,
						Tag:        common.ConfigReloaderTag,
					},
					BufferVolumeImage: v1beta1.ImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
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
						Replicas: 1,
						Drain: v1beta1.FluentdDrainConfig{
							Enabled: true,
							Image: v1beta1.ImageSpec{
								Repository: common.FluentdDrainWatchRepo,
								Tag:        common.FluentdDrainWatchTag,
							},
						},
					},
					Workers: 1,
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      outputName,
				Namespace: ns,
			},
			Spec: v1beta1.OutputSpec{
				FileOutput: &output.FileOutputConfig{
					Path:   "/tmp/logs/${tag}/%Y/%m/%d.%H.%M",
					Append: true,
					Buffer: &output.Buffer{
						Type:        "file",
						Timekey:     "1m",
						TimekeyWait: "10s",
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
				Name:      flowName,
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

		aggregatorLabels := map[string]string{
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
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggregatorLabels)); !aggregatorRunning() {
				t.Log("waiting for the aggregator")
				return false
			}

			return logging.Status.ProblemsCount == 0
		}, 5*time.Minute, 3*time.Second)

		t.Logf("Breaking File Output with an invalid config")
		patch := client.MergeFrom(output.DeepCopy())
		output.Spec.FileOutput.Path = "/tmp/zzz"
		common.RequireNoError(t, c.GetClient().Patch(ctx, &output, patch))
		require.Eventually(t, func() bool {
			common.RequireNoError(t, c.GetClient().Get(ctx, utils.ObjectKeyFromObjectMeta(&logging), &logging))
			if logging.Status.ProblemsCount > 0 {
				for _, problem := range logging.Status.Problems {
					match, err := regexp.MatchString(`^Configuration with checksum (.+) has failed. .*`, problem)
					common.RequireNoError(t, err)
					if match {
						t.Logf("Found the problem in Logging status: %v", logging.Status)
						return true
					}
				}
			}
			t.Logf("Waiting for the problem to appear in Logging status: %v", logging.Status.Problems)
			return false
		}, 5*time.Minute, 3*time.Second)

		t.Logf("Fixing Output")
		patch = client.MergeFrom(output.DeepCopy())
		output.Spec.FileOutput.Path = "/tmp/logs/${tag}/%Y/%m/%d.%H.%M"
		common.RequireNoError(t, c.GetClient().Patch(ctx, &output, patch))
		require.Eventually(t, func() bool {
			common.RequireNoError(t, c.GetClient().Get(ctx, utils.ObjectKeyFromObjectMeta(&logging), &logging))
			if logging.Status.ProblemsCount > 0 {
				for _, problem := range logging.Status.Problems {
					match, err := regexp.MatchString(`^Configuration with checksum (.+) has failed. .*`, problem)
					common.RequireNoError(t, err)
					if match {
						t.Logf("Waiting for the problem to be cleared in Logging status: %v", logging.Status.Problems)
						return false
					}
				}
			}
			t.Logf("Problem cleared in Logging status: %v", logging.Status)
			return true
		}, 5*time.Minute, 3*time.Second)
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
func TestFluentdAggregator_ConfigChecks_WhenReadOnlyRootFilesystemIsConfigured(t *testing.T) {
	common.Initialize(t)
	ns := "testing-3"
	releaseNameOverride := "e2e"
	outputName := "test-output"
	flowName := "test-flow"
	common.WithCluster("fluentd-configcheck-readonly-root-filesystem", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "fluentd-aggregator-configchecks-ro-rootfs-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				EnableRecreateWorkloadOnImmutableFieldChange: true,
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Network: &v1beta1.FluentbitNetwork{
						Keepalive: utils.BoolPointer(false),
					},
					ConfigHotReload: &v1beta1.HotReload{
						Image: v1beta1.ImageSpec{
							Repository: common.ConfigReloaderRepo,
							Tag:        common.ConfigReloaderTag,
						},
					},
					BufferVolumeImage: v1beta1.ImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
				},
				FluentdSpec: &v1beta1.FluentdSpec{
					Image: v1beta1.ImageSpec{
						Repository: common.FluentdImageRepo,
						Tag:        common.FluentdImageTag,
					},
					ConfigReloaderImage: v1beta1.ImageSpec{
						Repository: common.ConfigReloaderRepo,
						Tag:        common.ConfigReloaderTag,
					},
					BufferVolumeImage: v1beta1.ImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
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
						Replicas: 1,
						Drain: v1beta1.FluentdDrainConfig{
							Enabled: true,
							Image: v1beta1.ImageSpec{
								Repository: common.FluentdDrainWatchRepo,
								Tag:        common.FluentdDrainWatchTag,
							},
						},
					},
					Workers: 1,
					Security: &v1beta1.Security{
						SecurityContext: &corev1.SecurityContext{
							ReadOnlyRootFilesystem: utils.BoolPointer(true),
						},
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      outputName,
				Namespace: ns,
			},
			Spec: v1beta1.OutputSpec{
				FileOutput: &output.FileOutputConfig{
					Path:   "/tmp/logs/${tag}/%Y/%m/%d.%H.%M",
					Append: true,
					Buffer: &output.Buffer{
						Type:        "file",
						Timekey:     "1m",
						TimekeyWait: "10s",
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
				Name:      flowName,
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

		aggregatorLabels := map[string]string{
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
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggregatorLabels)); !aggregatorRunning() {
				t.Log("waiting for the aggregator")
				return false
			}

			return logging.Status.ProblemsCount == 0
		}, 5*time.Minute, 3*time.Second)
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
