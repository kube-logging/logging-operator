// Copyright © 2021 Banzai Cloud
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
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/output"
	"github.com/banzaicloud/operator-tools/pkg/utils"
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

	"github.com/banzaicloud/logging-operator/e2e/common"
	"github.com/banzaicloud/logging-operator/e2e/common/cond"
	"github.com/banzaicloud/logging-operator/e2e/common/setup"
)

func TestVolumeDrain_Downscale(t *testing.T) {
	common.WithCluster(t, func(t *testing.T, c common.Cluster) {
		ns := "testing-1"

		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Config.DisableWebhook = true
			options.Config.Namespace = ns
		}))

		consumer := setup.LogConsumer(t, c.GetClient(), setup.LogConsumerOptionFunc(func(options *setup.LogConsumerOptions) {
			options.Namespace = ns
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
		require.NoError(t, c.GetClient().Create(ctx, &logging))
		tags := "time"
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-output",
				Namespace: ns,
			},
			Spec: v1beta1.OutputSpec{
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint: consumer.InputURL(),
					Buffer: &output.Buffer{
						Type:        "file",
						Tags:        &tags,
						Timekey:     "10s",
						TimekeyWait: "0s",
					},
				},
			},
		}
		require.NoError(t, c.GetClient().Create(ctx, &output))
		flow := v1beta1.Flow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-flow",
				Namespace: ns,
			},
			Spec: v1beta1.FlowSpec{
				Match: []v1beta1.Match{
					{
						Select: &v1beta1.Select{
							Labels: map[string]string{
								"my-unique-label": "log-producer",
							},
						},
					},
				},
				LocalOutputRefs: []string{output.Name},
			},
		}
		require.NoError(t, c.GetClient().Create(ctx, &flow))

		fluentdReplicaName := logging.Name + "-fluentd-1"
		require.Eventually(t, cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: fluentdReplicaName}), 2*time.Minute, 5*time.Second)

		setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = ns
			options.Labels = flow.Spec.Match[0].Select.Labels
		}))

		require.Eventually(t, func() bool {
			rawOut, err := exec.Command("kubectl", "-n", consumer.PodKey.Namespace, "logs", consumer.PodKey.Name).Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), "got request")
		}, 5*time.Minute, 2*time.Second)

		require.NoError(t, exec.Command("kubectl", "-n", consumer.PodKey.Namespace, "exec", consumer.PodKey.Name, "--", "curl", "-sS", "http://localhost:8082/off").Run())

		require.Eventually(t, func() bool {
			rawOut, err := exec.Command("kubectl", "-n", ns, "exec", fluentdReplicaName, "-c", "fluentd", "--", "ls", "-1", "/buffers").Output()
			if err != nil {
				t.Logf("failed to list buffer directory: %v", err)
				return false
			}
			return strings.Count(string(rawOut), "\n") > 2
		}, 3*time.Minute, 10*time.Second)

		patch := client.MergeFrom(logging.DeepCopy())
		logging.Spec.FluentdSpec.Scaling.Replicas = 1
		require.NoError(t, c.GetClient().Patch(ctx, &logging, patch))

		drainerJobName := fluentdReplicaName + "-drainer"
		require.Eventually(t, func() bool {
			var job batchv1.Job
			present := cond.ResourceShouldBePresent(t, c.GetClient(), common.Resource(&job, ns, drainerJobName))()
			return present && job.Status.Active > 0
		}, 2*time.Minute, 1*time.Second)

		require.Eventually(t, cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: fluentdReplicaName}), 30*time.Second, time.Second/2)

		require.NoError(t, exec.Command("kubectl", "-n", consumer.PodKey.Namespace, "exec", consumer.PodKey.Name, "--", "curl", "-sS", "http://localhost:8082/on").Run())

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(batchv1.Job), ns, drainerJobName)), 5*time.Minute, 30*time.Second)

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(corev1.Pod), ns, fluentdReplicaName)), 30*time.Second, time.Second/2)

		pvc := common.Resource(new(corev1.PersistentVolumeClaim), ns, logging.Name+"-fluentd-buffer-"+fluentdReplicaName)
		require.NoError(t, c.GetClient().Get(ctx, client.ObjectKeyFromObject(pvc), pvc))
		assert.Equal(t, "drained", pvc.GetLabels()["logging.banzaicloud.io/drain-status"])
	}, func(o *cluster.Options) {
		if o.Scheme == nil {
			o.Scheme = runtime.NewScheme()
		}
		require.NoError(t, v1beta1.AddToScheme(o.Scheme))
		require.NoError(t, apiextensionsv1.AddToScheme(o.Scheme))
		require.NoError(t, appsv1.AddToScheme(o.Scheme))
		require.NoError(t, batchv1.AddToScheme(o.Scheme))
		require.NoError(t, corev1.AddToScheme(o.Scheme))
		require.NoError(t, rbacv1.AddToScheme(o.Scheme))
	})
}

func TestVolumeDrain_Downscale_DeleteVolume(t *testing.T) {
	common.WithCluster(t, func(t *testing.T, c common.Cluster) {
		ns := "testing-2"

		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Config.DisableWebhook = true
			options.Config.Namespace = ns
		}))

		consumer := setup.LogConsumer(t, c.GetClient(), setup.LogConsumerOptionFunc(func(options *setup.LogConsumerOptions) {
			options.Namespace = ns
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
							Enabled:      true,
							DeleteVolume: true,
						},
					},
				},
			},
		}
		require.NoError(t, c.GetClient().Create(ctx, &logging))
		tags := "time"
		output := v1beta1.Output{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-output",
				Namespace: ns,
			},
			Spec: v1beta1.OutputSpec{
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint: consumer.InputURL(),
					Buffer: &output.Buffer{
						Type:        "file",
						Tags:        &tags,
						Timekey:     "10s",
						TimekeyWait: "0s",
					},
				},
			},
		}
		require.NoError(t, c.GetClient().Create(ctx, &output))
		flow := v1beta1.Flow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-flow",
				Namespace: ns,
			},
			Spec: v1beta1.FlowSpec{
				Match: []v1beta1.Match{
					{
						Select: &v1beta1.Select{
							Labels: map[string]string{
								"my-unique-label": "log-producer",
							},
						},
					},
				},
				LocalOutputRefs: []string{output.Name},
			},
		}
		require.NoError(t, c.GetClient().Create(ctx, &flow))

		fluentdReplicaName := logging.Name + "-fluentd-1"
		require.Eventually(t, cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: fluentdReplicaName}), 2*time.Minute, 5*time.Second)

		setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = ns
			options.Labels = flow.Spec.Match[0].Select.Labels
		}))

		require.Eventually(t, func() bool {
			rawOut, err := exec.Command("kubectl", "-n", consumer.PodKey.Namespace, "logs", consumer.PodKey.Name).Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), "got request")
		}, 5*time.Minute, 2*time.Second)

		require.NoError(t, exec.Command("kubectl", "-n", consumer.PodKey.Namespace, "exec", consumer.PodKey.Name, "--", "curl", "-sS", "http://localhost:8082/off").Run())

		require.Eventually(t, func() bool {
			rawOut, err := exec.Command("kubectl", "-n", ns, "exec", fluentdReplicaName, "-c", "fluentd", "--", "ls", "-1", "/buffers").Output()
			if err != nil {
				t.Logf("failed to list buffer directory: %v", err)
				return false
			}
			return strings.Count(string(rawOut), "\n") > 2
		}, 3*time.Minute, 10*time.Second)

		patch := client.MergeFrom(logging.DeepCopy())
		logging.Spec.FluentdSpec.Scaling.Replicas = 1
		require.NoError(t, c.GetClient().Patch(ctx, &logging, patch))

		drainerJobName := fluentdReplicaName + "-drainer"
		require.Eventually(t, func() bool {
			var job batchv1.Job
			present := cond.ResourceShouldBePresent(t, c.GetClient(), common.Resource(&job, ns, drainerJobName))()
			return present && job.Status.Active > 0
		}, 2*time.Minute, 1*time.Second)

		require.Eventually(t, cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: fluentdReplicaName}), 30*time.Second, time.Second/2)

		require.NoError(t, exec.Command("kubectl", "-n", consumer.PodKey.Namespace, "exec", consumer.PodKey.Name, "--", "curl", "-sS", "http://localhost:8082/on").Run())

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(batchv1.Job), ns, drainerJobName)), 5*time.Minute, 30*time.Second)

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(corev1.Pod), ns, fluentdReplicaName)), 30*time.Second, time.Second/2)

		require.Eventually(t, cond.ResourceShouldBeAbsent(t, c.GetClient(), common.Resource(new(corev1.PersistentVolumeClaim), ns, logging.Name+"-fluentd-buffer-"+fluentdReplicaName)), 30*time.Second, time.Second/2)
	}, func(o *cluster.Options) {
		if o.Scheme == nil {
			o.Scheme = runtime.NewScheme()
		}
		require.NoError(t, v1beta1.AddToScheme(o.Scheme))
		require.NoError(t, apiextensionsv1.AddToScheme(o.Scheme))
		require.NoError(t, appsv1.AddToScheme(o.Scheme))
		require.NoError(t, batchv1.AddToScheme(o.Scheme))
		require.NoError(t, corev1.AddToScheme(o.Scheme))
		require.NoError(t, rbacv1.AddToScheme(o.Scheme))
	})
}
