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

package logging_metrics_monitoring_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cisco-open/operator-tools/pkg/utils"
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/common/setup"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"sigs.k8s.io/e2e-framework/third_party/helm"
)

type metricsTester struct {
	testPod *corev1.Pod
}

type metricsEndpoint struct {
	serviceName string
	port        int
	path        string
}

type loggingResourceName string

const (
	pollInterval = 5 * time.Second
	pollTimeout  = 5 * time.Minute

	fluentbit loggingResourceName = "fluentbit"
	syslogNG  loggingResourceName = "syslog-ng"
	fluentd   loggingResourceName = "fluentd"

	fluentbitServiceName              = "metrics-monitoring-test-" + string(fluentbit) + "-metrics"
	fluentbitBufferMetricsServiceName = "metrics-monitoring-test-" + string(fluentbit) + "-buffer-metrics"
	syslogNGServiceName               = "metrics-monitoring-test-" + string(syslogNG) + "-metrics"
	syslogNGBufferMetricsServiceName  = "metrics-monitoring-test-" + string(syslogNG) + "-buffer-metrics"
	fluentdServiceName                = "metrics-monitoring-test-" + string(fluentd) + "-metrics"
	fluentdBufferMetricsServiceName   = "metrics-monitoring-test-" + string(fluentd) + "-buffer-metrics"
)

var metricServices = map[loggingResourceName]metricsEndpoint{
	fluentbit: {
		serviceName: fluentbitServiceName,
		port:        2020,
		path:        "/api/v1/metrics/prometheus",
	},
	syslogNG: {
		serviceName: syslogNGServiceName,
		port:        9577,
		path:        "/metrics",
	},
	fluentd: {
		serviceName: fluentdServiceName,
		port:        24231,
		path:        "/metrics",
	},
}

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

func TestLoggingMetrics_Monitoring(t *testing.T) {
	common.Initialize(t)
	ns := "test"
	releaseNameOverride := "e2e"
	common.WithCluster("logging-metrics-monitoring", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		common.RequireNoError(t, installPrometheusOperator(c))

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "metrics-monitoring-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Metrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
					},
					BufferVolumeMetrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
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
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					ConfigReloadImage: &v1beta1.BasicImageSpec{
						Repository: common.SyslogNGReloaderRepo,
						Tag:        common.SyslogNGReloaderTag,
					},
					BufferVolumeMetricsImage: &v1beta1.BasicImageSpec{
						Repository: common.NodeExporterRepo,
						Tag:        common.NodeExporterTag,
					},
					Metrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
					},
					BufferVolumeMetrics: &v1beta1.BufferMetrics{
						Metrics: v1beta1.Metrics{
							Enabled:        utils.BoolPointer(true),
							ServiceMonitor: true,
						},
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		common.RequireNoError(t, common.WaitForPodReady(ctx, c.GetClient(), &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      logging.Name + "-" + string(syslogNG) + "-0",
				Namespace: ns,
			},
		}, pollInterval, pollTimeout))
		serviceMonitorsSyslogNG := &v1.ServiceMonitorList{}
		common.RequireNoError(t, c.GetClient().List(ctx, serviceMonitorsSyslogNG))

		mt, err := setupMetricsTester(ctx, c, ns)
		common.RequireNoError(t, err)

		require.Eventually(t, func() bool {
			rawOut, err := mt.getMetrics(metricServices[fluentbit], c, ns)
			if err != nil {
				t.Log(err)
				return false
			}
			if err := mt.validateMetrics(rawOut, fluentbit); err != nil {
				t.Log(err)
				return false
			}
			return true
		}, pollTimeout, pollInterval)

		require.Eventually(t, func() bool {
			rawOut, err := mt.getMetrics(metricServices[syslogNG], c, ns)
			if err != nil {
				t.Log(err)
				return false
			}
			if err := mt.validateMetrics(rawOut, syslogNG); err != nil {
				t.Log(err)
				return false
			}
			return true
		}, pollTimeout, pollInterval)

		common.RequireNoError(t, c.GetClient().Delete(ctx, &logging))

		loggingPatch := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "metrics-monitoring-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Metrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
					},
					BufferVolumeMetrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
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
					Metrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
					},
					BufferVolumeMetrics: &v1beta1.Metrics{
						Enabled:        utils.BoolPointer(true),
						ServiceMonitor: true,
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &loggingPatch))
		common.RequireNoError(t, common.WaitForPodReady(ctx, c.GetClient(), &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      loggingPatch.Name + "-" + string(fluentd) + "-0",
				Namespace: ns,
			},
		}, pollInterval, pollTimeout))
		serviceMonitorsFluentd := &v1.ServiceMonitorList{}
		common.RequireNoError(t, c.GetClient().List(ctx, serviceMonitorsFluentd))

		require.Eventually(t, func() bool {
			rawOut, err := mt.getMetrics(metricServices[fluentd], c, ns)
			if err != nil {
				t.Log(err)
				return false
			}
			if err := mt.validateMetrics(rawOut, fluentd); err != nil {
				t.Log(err)
				return false
			}
			return true
		}, pollTimeout, pollInterval)

		serviceMonitors := append(serviceMonitorsFluentd.Items, serviceMonitorsSyslogNG.Items...)
		common.RequireNoError(t, checkServiceMonitorAvailability(serviceMonitors))

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
		common.RequireNoError(t, v1.AddToScheme(o.Scheme))
	})
}

func installPrometheusOperator(c common.Cluster) error {
	manager := helm.New(c.KubeConfigFilePath())

	if err := manager.RunRepo(helm.WithArgs("add", "prometheus-community", "https://prometheus-community.github.io/helm-charts")); err != nil {
		return fmt.Errorf("failed to add prometheus-community repo: %v", err)
	}

	if err := manager.RunRepo(helm.WithArgs("update")); err != nil {
		return fmt.Errorf("failed to update helm repos: %v", err)
	}

	if err := manager.RunInstall(
		helm.WithName("prometheus"),
		helm.WithChart("prometheus-community/kube-prometheus-stack"),
		helm.WithArgs("--create-namespace"),
		helm.WithNamespace("monitoring"),
		helm.WithArgs("--set", "prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false"),
		helm.WithArgs("--set", "prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false"),
		helm.WithArgs("--set", "prometheus.prometheusSpec.maximumStartupDurationSeconds=900"),
		helm.WithWait(),
	); err != nil {
		return fmt.Errorf("failed to install prometheus: %v", err)
	}

	return nil
}

func setupMetricsTester(ctx context.Context, c common.Cluster, ns string) (metricsTester, error) {
	pod, err := common.SetupCurlPod(ctx, c.GetClient(), ns, "metrics-tester", pollInterval, pollTimeout)
	if err != nil {
		return metricsTester{}, err
	}

	return metricsTester{
		testPod: pod,
	}, nil
}

func checkServiceMonitorAvailability(serviceMonitors []v1.ServiceMonitor) error {
	if len(serviceMonitors) == 0 {
		return errors.New("no service monitors found")
	}

	var expectedServiceMonitors = map[string]bool{
		fluentbitServiceName:              false,
		fluentbitBufferMetricsServiceName: false,
		syslogNGServiceName:               false,
		syslogNGBufferMetricsServiceName:  false,
		fluentdServiceName:                false,
		fluentdBufferMetricsServiceName:   false,
	}

	for _, sm := range serviceMonitors {
		delete(expectedServiceMonitors, sm.Name)
	}

	if len(expectedServiceMonitors) > 0 {
		return fmt.Errorf("the following service monitors are missing: %v", expectedServiceMonitors)
	}

	return nil
}

func (mt *metricsTester) getMetrics(endpoint metricsEndpoint, c common.Cluster, ns string) ([]byte, error) {
	serviceURL := fmt.Sprintf("http://%s.%s.svc:%d%s",
		endpoint.serviceName,
		ns,
		endpoint.port,
		endpoint.path,
	)
	cmd := common.CmdEnv(exec.Command("kubectl", "exec", mt.testPod.Name, "-n", ns, "--", "curl", serviceURL), c)
	rawOut, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return rawOut, nil
}

func (mt *metricsTester) validateMetrics(rawOut []byte, subject loggingResourceName) error {
	var missingMetrics []string
	for _, metric := range getKeyMetricsFor(subject) {
		if !strings.Contains(string(rawOut), metric) {
			missingMetrics = append(missingMetrics, metric)
		}
	}

	if len(missingMetrics) > 0 {
		return fmt.Errorf("for %s metrics, the following key metrics were not found: %v\n"+
			"Total metrics missing: %d\n"+
			"Full metrics response: %s",
			subject,
			missingMetrics,
			len(missingMetrics),
			string(rawOut),
		)
	}

	return nil
}

func getKeyMetricsFor(subject loggingResourceName) []string {
	keyMetrics := map[loggingResourceName][]string{
		fluentbit: {
			"fluentbit_input_records_total",
			"fluentbit_input_bytes_total",
			"fluentbit_filter_add_records_total",
			"fluentbit_filter_bytes_total",
			"fluentbit_output_retried_records_total",
			"fluentbit_output_retried_records_total",
		},
		syslogNG: {
			"syslogng_events_allocated_bytes",
			"syslogng_scratch_buffers_count",
			"syslogng_scratch_buffers_bytes",
		},
		fluentd: {
			"fluentd_output_status_retry_count",
			"fluentd_output_status_num_errors",
			"fluentd_output_status_emit_count",
			"fluentd_output_status_emit_records",
			"fluentd_output_status_write_count",
		},
	}

	return keyMetrics[subject]
}
