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

package elasticsearch_multiversion

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/common/setup"
	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

const (
	// esReadyFallback applies when the binary has no deadline to derive from.
	esReadyFallback = 15 * time.Minute

	// esReadyMargin is the room left before the deadline. 30s covered the wait
	// reporting its own failure but not teardown, which is a stern dump, a
	// kubectl exec plus tar for coverage, and a kind delete.
	esReadyMargin = 90 * time.Second
)

// logContainerRestarts names a restart loop while a readiness wait is still
// running, so it is a reported number rather than an unexplained stall.
func logContainerRestarts(t *testing.T, env *harness.Env, ns, app string) {
	var pods corev1.PodList
	if err := env.Client.List(env.Ctx, &pods, client.InNamespace(ns), client.MatchingLabels{"app": app}); err != nil {
		t.Logf("listing %s pods failed: %v", app, err)
		return
	}

	for _, pod := range pods.Items {
		for _, status := range pod.Status.ContainerStatuses {
			if status.RestartCount == 0 {
				continue
			}
			reason := "unknown"
			if terminated := status.LastTerminationState.Terminated; terminated != nil {
				reason = fmt.Sprintf("%s, exit %d", terminated.Reason, terminated.ExitCode)
			}
			t.Logf("%s/%s restarted %d times, last termination: %s", pod.Name, status.Name, status.RestartCount, reason)
		}
	}
}

// esReadyBudget shares what is left of the -timeout between the deployments
// still to check. Handing it all to the first let a slow one spend the
// package's budget, leaving the rest to fail instantly under the wrong name.
// The fallback is already per-wait, so it is not divided.
func esReadyBudget(t *testing.T, waitsLeft int) time.Duration {
	deadline, ok := t.Deadline()
	if !ok {
		return esReadyFallback
	}
	return budgetWithin(time.Until(deadline), waitsLeft)
}

// budgetWithin floors after dividing, not before. A third of a second is
// positive but expires on the first tick, so flooring on the undivided
// remainder would put the misattribution this split removes back at the
// boundary: elasticsearch7 blamed for a package that had already run out.
func budgetWithin(remaining time.Duration, waitsLeft int) time.Duration {
	if budget := (remaining - esReadyMargin) / time.Duration(waitsLeft); budget > time.Second {
		return budget
	}
	return time.Second
}

func TestBudgetWithin(t *testing.T) {
	for _, c := range []struct {
		name      string
		remaining time.Duration
		waitsLeft int
		want      time.Duration
	}{
		{"shared between the waits still to come", 20 * time.Minute, 3, (20*time.Minute - esReadyMargin) / 3},
		{"the last wait gets what is left", 20 * time.Minute, 1, 20*time.Minute - esReadyMargin},
		// The floor is a second per wait, not a second shared between them.
		{"just above the margin", esReadyMargin + time.Second, 3, time.Second},
		{"a slice under a second", esReadyMargin + 2*time.Second, 3, time.Second},
		{"exactly the margin", esReadyMargin, 3, time.Second},
		{"already past the deadline", -time.Minute, 3, time.Second},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, budgetWithin(c.remaining, c.waitsLeft))
			require.Positive(t, budgetWithin(c.remaining, c.waitsLeft))
		})
	}
}

// The whole point of the split: no single wait may take the budget the others
// still need.
func TestBudgetLeavesRoomForTheRemainingWaits(t *testing.T) {
	const remaining = 20 * time.Minute

	first := budgetWithin(remaining, 3)

	require.Less(t, first, remaining-esReadyMargin)
	require.LessOrEqual(t, 3*first, remaining-esReadyMargin)
}

// esClient reads Elasticsearch through the curl pod, which is the only way in
// from the test binary. It stays suite-local: no other suite has a curl pod, so
// on the harness it would be a helper with one caller.
type esClient struct {
	pod string
	ns  string
	env *harness.Env
}

// hasDocuments asks _cat/count for one index pattern. An empty body is what a
// pattern that matches nothing returns, so it counts as not yet rather than as
// an error.
func (c esClient) hasDocuments(t *testing.T, host, index string) bool {
	t.Helper()

	url := fmt.Sprintf("http://%s.%s.svc:9200/_cat/count/%s?h=count", host, c.ns, index)
	rawOut, err := common.CmdEnv(exec.Command("kubectl", "exec", c.pod, "-n", c.ns, "--",
		"curl", "-s", url), c.env.Cluster).Output()
	if err != nil {
		t.Logf("Error checking %s: %v", host, err)
		return false
	}

	count := strings.TrimSpace(string(rawOut))
	t.Logf("%s document count: %s", host, count)
	return count != "" && count != "0"
}

// esVersion is the whole difference between the three deployments. Container
// spec, probes, resources, env and ports are identical across them, so the
// versions are a table and the objects are built from it.
type esVersion struct {
	name    string
	version string
}

var esVersions = []esVersion{
	{"elasticsearch7", "7.17.16"},
	{"elasticsearch8", "8.12.0"},
	{"elasticsearch9", "9.1.5"},
}

func (v esVersion) labels() map[string]string {
	return map[string]string{"app": v.name, "version": v.version}
}

func esService(ns string, v esVersion) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.name,
			Namespace: ns,
			Labels:    v.labels(),
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     9200,
					Protocol: corev1.ProtocolTCP,
				},
				{
					Name:     "transport",
					Port:     9300,
					Protocol: corev1.ProtocolTCP,
				},
			},
			Selector: map[string]string{"app": v.name},
		},
	}
}

func esDeployment(ns string, v esVersion) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      v.name,
			Namespace: ns,
			Labels:    v.labels(),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: new(int32(1)),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": v.name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: v.labels()},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "elasticsearch",
							Image: "docker.elastic.co/elasticsearch/elasticsearch:" + v.version,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 9200,
									Protocol:      corev1.ProtocolTCP,
								},
								{
									Name:          "transport",
									ContainerPort: 9300,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "discovery.type",
									Value: "single-node",
								},
								{
									Name:  "ES_JAVA_OPTS",
									Value: "-Xms512m -Xmx512m",
								},
								{
									Name:  "xpack.security.enabled",
									Value: "false",
								},
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("1Gi"),
									corev1.ResourceCPU:    resource.MustParse("500m"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("1536Mi"),
									corev1.ResourceCPU:    resource.MustParse("1000m"),
								},
							},
							// No liveness probe: readiness already gates the wait, and a
							// 60s + 3x10s deadline killed the JVM mid-boot on a loaded runner.
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/_cluster/health",
										Port: intstr.FromInt(9200),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
							},
						},
					},
				},
			},
		},
	}
}

func TestElasticsearch_MultiVersion(t *testing.T) {
	const (
		ns      = "logging"
		release = "e2e"
	)

	env := harness.New(t).
		WithCluster("elasticsearch-multiversion").
		WithRelease(release).
		WithControlNamespace(ns).
		Start()

	for _, v := range esVersions {
		env.Create(esService(ns, v), esDeployment(ns, v))
	}
	for i, v := range esVersions {
		t.Logf("Waiting for %s deployment to be ready...", v.name)
		// Driving the Condition from the suite's own poll rather than
		// env.WaitWithin: the budget WaitWithin would take, but the restart
		// count logged on each failed poll a read-only Condition cannot
		// produce, and it is what turns a stall into a reported number.
		require.Eventuallyf(t, func() bool {
			if met, _ := wait.Deployment(ns, v.name).Met(env.Ctx, env.Client); met {
				return true
			}
			logContainerRestarts(t, env, ns, v.name)
			return false
		}, esReadyBudget(t, len(esVersions)-i), 10*time.Second, "%s never became ready", v.name)
	}

	logging := v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name: "all-to-es",
		},
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: ns,
			FluentdSpec: &v1beta1.FluentdSpec{
				Image:               image.Fluentd().Spec(),
				ConfigReloaderImage: image.ConfigReloader().Spec(),
				BufferVolumeImage:   image.NodeExporter().Spec(),
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
				LogLevel: "debug",
			},
		},
	}
	env.Create(&logging)

	agent := v1beta1.FluentbitAgent{
		ObjectMeta: metav1.ObjectMeta{
			Name: "all-to-es",
		},
		Spec: v1beta1.FluentbitSpec{},
	}
	env.Create(&agent)

	es7Output := v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "es7-output",
			Namespace: ns,
		},
		Spec: v1beta1.OutputSpec{
			ElasticsearchOutput: &output.ElasticsearchOutput{
				Host:                        "elasticsearch7.logging.svc.cluster.local",
				Port:                        9200,
				Scheme:                      "http",
				DefaultElasticsearchVersion: "7",
				SuppressTypeName:            new(false),
				TypeName:                    "_doc",
				IndexName:                   "test-logs-es7",
				LogstashFormat:              true,
				LogstashPrefix:              "fluentd-es7",
				LogstashDateformat:          "%Y.%m.%d",
				IncludeTimestamp:            true,
				ReconnectOnError:            true,
				ReloadConnections:           new(false),
				ReloadOnFailure:             true,
				VerifyEsVersionAtStartup:    new(false),
				Buffer: &output.Buffer{
					Type:             "file",
					Path:             "/buffers/es7",
					ChunkLimitSize:   "4MB",
					FlushAtShutdown:  true,
					FlushInterval:    "15s",
					FlushMode:        "interval",
					FlushThreadCount: 2,
					OverflowAction:   "block",
					RetryMaxInterval: "30s",
					RetryTimeout:     "72h",
				},
			},
		},
	}
	env.Create(&es7Output)

	es8Output := v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "es8-output",
			Namespace: ns,
		},
		Spec: v1beta1.OutputSpec{
			ElasticsearchOutput: &output.ElasticsearchOutput{
				Host:                        "elasticsearch8.logging.svc.cluster.local",
				Port:                        9200,
				Scheme:                      "http",
				DefaultElasticsearchVersion: "8",
				SuppressTypeName:            new(true),
				DataStreamEnable:            new(true),
				DataStreamName:              "logs-fluentd-es8",
				DataStreamTemplateName:      "logs-fluentd-template",
				IncludeTimestamp:            true,
				ReconnectOnError:            true,
				ReloadConnections:           new(false),
				ReloadOnFailure:             true,
				VerifyEsVersionAtStartup:    new(false),
				Buffer: &output.Buffer{
					Type:             "file",
					Path:             "/buffers/es8",
					ChunkLimitSize:   "4MB",
					FlushAtShutdown:  true,
					FlushInterval:    "15s",
					FlushMode:        "interval",
					FlushThreadCount: 2,
					OverflowAction:   "block",
					RetryMaxInterval: "30s",
					RetryTimeout:     "72h",
				},
			},
		},
	}
	env.Create(&es8Output)

	es9Output := v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "es9-output",
			Namespace: ns,
		},
		Spec: v1beta1.OutputSpec{
			ElasticsearchOutput: &output.ElasticsearchOutput{
				Host:                        "elasticsearch9.logging.svc.cluster.local",
				Port:                        9200,
				Scheme:                      "http",
				DefaultElasticsearchVersion: "9",
				SuppressTypeName:            new(true),
				DataStreamEnable:            new(true),
				DataStreamName:              "logs-fluentd-es9",
				DataStreamTemplateName:      "logs-fluentd-template",
				IncludeTimestamp:            true,
				ReconnectOnError:            true,
				ReloadConnections:           new(false),
				ReloadOnFailure:             true,
				VerifyEsVersionAtStartup:    new(false),
				Buffer: &output.Buffer{
					Type:             "file",
					Path:             "/buffers/es9",
					ChunkLimitSize:   "4MB",
					FlushAtShutdown:  true,
					FlushInterval:    "15s",
					FlushMode:        "interval",
					FlushThreadCount: 2,
					OverflowAction:   "block",
					RetryMaxInterval: "30s",
					RetryTimeout:     "72h",
				},
			},
		},
	}
	env.Create(&es9Output)

	producerLabels := map[string]string{
		"my-unique-label": "log-producer",
	}

	flow := v1beta1.Flow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "all-logs-to-elasticsearch",
			Namespace: ns,
		},
		Spec: v1beta1.FlowSpec{
			Filters: []v1beta1.Filter{
				{
					TagNormaliser: &filter.TagNormaliser{},
				},
				{
					RecordModifier: &filter.RecordModifier{
						Records: []filter.Record{
							{"cluster": "test-cluster"},
							{"environment": "development"},
						},
					},
				},
				{
					RecordTransformer: &filter.RecordTransformer{
						EnableRuby: true,
						Records: []filter.Record{
							{"kubernetes_labels_flattened": `${record.dig("kubernetes", "labels").to_json rescue "{}"}`},
						},
						RemoveKeys: "kubernetes.labels",
					},
				},
			},
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{},
				},
			},
			LocalOutputRefs: []string{"es7-output", "es8-output", "es9-output"},
		},
	}
	env.Create(&flow)

	setup.LogProducer(t, env.Client, setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
		options.Namespace = ns
		options.Labels = producerLabels
	}))

	t.Log("Waiting for components to be ready...")
	env.WaitFor(
		wait.Operator(release),
		wait.Producer(producerLabels),
		wait.FluentdAggregator(ns),
	)

	const (
		pollInterval = 5 * time.Second
		pollTimeout  = 2 * time.Minute
	)
	curlPod, err := common.SetupCurlPod(env.Ctx, env.Client, ns, "es-tester", pollInterval, pollTimeout)
	common.RequireNoError(t, err)
	es := esClient{pod: curlPod.Name, ns: ns, env: env}

	require.Eventuallyf(t, func() bool {
		return es.hasDocuments(t, "elasticsearch7", "fluentd-es7-*")
	}, 3*time.Minute, 10*time.Second, "elasticsearch7 never indexed anything")

	require.Eventuallyf(t, func() bool {
		return es.hasDocuments(t, "elasticsearch8", ".ds-logs-fluentd-es8-*")
	}, 3*time.Minute, 10*time.Second, "elasticsearch8 never indexed anything")

	require.Eventuallyf(t, func() bool {
		return es.hasDocuments(t, "elasticsearch9", ".ds-logs-fluentd-es9-*")
	}, 3*time.Minute, 10*time.Second, "elasticsearch9 never indexed anything")
}
