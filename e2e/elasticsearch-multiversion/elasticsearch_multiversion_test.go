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
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/filter"
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

func TestElasticsearch_MultiVersion(t *testing.T) {
	common.Initialize(t)
	ns := "logging"
	releaseNameOverride := "e2e"
	common.WithCluster("elasticsearch-multiversion", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		es7Service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "elasticsearch7",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "elasticsearch7",
					"version": "7.17.16",
				},
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
				Selector: map[string]string{
					"app": "elasticsearch7",
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, es7Service))

		es7Deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "elasticsearch7",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "elasticsearch7",
					"version": "7.17.16",
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "elasticsearch7",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":     "elasticsearch7",
							"version": "7.17.16",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "elasticsearch",
								Image: "docker.elastic.co/elasticsearch/elasticsearch:7.17.16",
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
										corev1.ResourceMemory: resource.MustParse("1Gi"),
										corev1.ResourceCPU:    resource.MustParse("1000m"),
									},
								},
								LivenessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/_cluster/health",
											Port: intstr.FromInt(9200),
										},
									},
									InitialDelaySeconds: 60,
									PeriodSeconds:       10,
								},
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
		common.RequireNoError(t, c.GetClient().Create(ctx, es7Deployment))

		es8Service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "elasticsearch8",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "elasticsearch8",
					"version": "8.12.0",
				},
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
				Selector: map[string]string{
					"app": "elasticsearch8",
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, es8Service))

		es8Deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "elasticsearch8",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "elasticsearch8",
					"version": "8.12.0",
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "elasticsearch8",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":     "elasticsearch8",
							"version": "8.12.0",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "elasticsearch",
								Image: "docker.elastic.co/elasticsearch/elasticsearch:8.12.0",
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
										corev1.ResourceMemory: resource.MustParse("1Gi"),
										corev1.ResourceCPU:    resource.MustParse("1000m"),
									},
								},
								LivenessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/_cluster/health",
											Port: intstr.FromInt(9200),
										},
									},
									InitialDelaySeconds: 60,
									PeriodSeconds:       10,
								},
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
		common.RequireNoError(t, c.GetClient().Create(ctx, es8Deployment))

		es9Service := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "elasticsearch9",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "elasticsearch9",
					"version": "9.1.5",
				},
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
				Selector: map[string]string{
					"app": "elasticsearch9",
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, es9Service))

		es9Deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "elasticsearch9",
				Namespace: ns,
				Labels: map[string]string{
					"app":     "elasticsearch9",
					"version": "9.1.5",
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: utils.IntPointer(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "elasticsearch9",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app":     "elasticsearch9",
							"version": "9.1.5",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "elasticsearch",
								Image: "docker.elastic.co/elasticsearch/elasticsearch:9.1.5",
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
										corev1.ResourceMemory: resource.MustParse("1Gi"),
										corev1.ResourceCPU:    resource.MustParse("1000m"),
									},
								},
								LivenessProbe: &corev1.Probe{
									ProbeHandler: corev1.ProbeHandler{
										HTTPGet: &corev1.HTTPGetAction{
											Path: "/_cluster/health",
											Port: intstr.FromInt(9200),
										},
									},
									InitialDelaySeconds: 60,
									PeriodSeconds:       10,
								},
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
		common.RequireNoError(t, c.GetClient().Create(ctx, es9Deployment))

		t.Log("Waiting for Elasticsearch 7 deployment to be ready...")
		require.Eventually(t, func() bool {
			return cond.DeploymentAvailable(t, c.GetClient(), &ctx, ns, "elasticsearch7")()
		}, 30*time.Minute, 10*time.Second)

		t.Log("Waiting for Elasticsearch 8 deployment to be ready...")
		require.Eventually(t, func() bool {
			return cond.DeploymentAvailable(t, c.GetClient(), &ctx, ns, "elasticsearch8")()
		}, 30*time.Minute, 10*time.Second)

		t.Log("Waiting for Elasticsearch 9 deployment to be ready...")
		require.Eventually(t, func() bool {
			return cond.DeploymentAvailable(t, c.GetClient(), &ctx, ns, "elasticsearch9")()
		}, 30*time.Minute, 10*time.Second)

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name: "all-to-es",
			},
			Spec: v1beta1.LoggingSpec{
				ControlNamespace: ns,
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
					LogLevel: "debug",
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))

		agent := v1beta1.FluentbitAgent{
			ObjectMeta: metav1.ObjectMeta{
				Name: "all-to-es",
			},
			Spec: v1beta1.FluentbitSpec{},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &agent))

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
					SuppressTypeName:            utils.BoolPointer(false),
					TypeName:                    "_doc",
					IndexName:                   "test-logs-es7",
					LogstashFormat:              true,
					LogstashPrefix:              "fluentd-es7",
					LogstashDateformat:          "%Y.%m.%d",
					IncludeTimestamp:            true,
					ReconnectOnError:            true,
					ReloadConnections:           utils.BoolPointer(false),
					ReloadOnFailure:             true,
					VerifyEsVersionAtStartup:    utils.BoolPointer(false),
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
		common.RequireNoError(t, c.GetClient().Create(ctx, &es7Output))

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
					SuppressTypeName:            utils.BoolPointer(true),
					DataStreamEnable:            utils.BoolPointer(true),
					DataStreamName:              "logs-fluentd-es8",
					DataStreamTemplateName:      "logs-fluentd-template",
					IncludeTimestamp:            true,
					ReconnectOnError:            true,
					ReloadConnections:           utils.BoolPointer(false),
					ReloadOnFailure:             true,
					VerifyEsVersionAtStartup:    utils.BoolPointer(false),
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
		common.RequireNoError(t, c.GetClient().Create(ctx, &es8Output))

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
					SuppressTypeName:            utils.BoolPointer(true),
					DataStreamEnable:            utils.BoolPointer(true),
					DataStreamName:              "logs-fluentd-es9",
					DataStreamTemplateName:      "logs-fluentd-template",
					IncludeTimestamp:            true,
					ReconnectOnError:            true,
					ReloadConnections:           utils.BoolPointer(false),
					ReloadOnFailure:             true,
					VerifyEsVersionAtStartup:    utils.BoolPointer(false),
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
		common.RequireNoError(t, c.GetClient().Create(ctx, &es9Output))

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

		t.Log("Waiting for components to be ready...")
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
			return true
		}, 5*time.Minute, 3*time.Second)

		const (
			pollInterval = 5 * time.Second
			pollTimeout  = 2 * time.Minute
		)
		curlPod, err := common.SetupCurlPod(ctx, c.GetClient(), ns, "es-tester", pollInterval, pollTimeout)
		common.RequireNoError(t, err)

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl", "exec", curlPod.Name, "-n", ns, "--",
				"curl", "-s", "http://elasticsearch7.logging.svc:9200/_cat/count/fluentd-es7-*?h=count"), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("Error checking ES7: %v", err)
				return false
			}
			count := strings.TrimSpace(string(rawOut))
			t.Logf("ES7 document count: %s", count)
			return count != "" && count != "0"
		}, 3*time.Minute, 10*time.Second)

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl", "exec", curlPod.Name, "-n", ns, "--",
				"curl", "-s", "http://elasticsearch8.logging.svc:9200/_cat/count/.ds-logs-fluentd-es8-*?h=count"), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("Error checking ES8: %v", err)
				return false
			}
			count := strings.TrimSpace(string(rawOut))
			t.Logf("ES8 document count: %s", count)
			return count != "" && count != "0"
		}, 3*time.Minute, 10*time.Second)

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl", "exec", curlPod.Name, "-n", ns, "--",
				"curl", "-s", "http://elasticsearch9.logging.svc:9200/_cat/count/.ds-logs-fluentd-es9-*?h=count"), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("Error checking ES9: %v", err)
				return false
			}
			count := strings.TrimSpace(string(rawOut))
			t.Logf("ES9 document count: %s", count)
			return count != "" && count != "0"
		}, 3*time.Minute, 10*time.Second)

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
		return err

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
