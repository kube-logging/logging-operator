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
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/e2e-framework/third_party/helm"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/image"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
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
	ns          = "test"
	release     = "e2e"
	loggingName = "metrics-monitoring-test"

	pollInterval = 5 * time.Second
	pollTimeout  = 5 * time.Minute

	fluentbit loggingResourceName = "fluentbit"
	syslogNG  loggingResourceName = "syslog-ng"
	fluentd   loggingResourceName = "fluentd"

	fluentbitServiceName              = loggingName + "-" + string(fluentbit) + "-metrics"
	fluentbitBufferMetricsServiceName = loggingName + "-" + string(fluentbit) + "-buffer-metrics"
	syslogNGServiceName               = loggingName + "-" + string(syslogNG) + "-metrics"
	syslogNGBufferMetricsServiceName  = loggingName + "-" + string(syslogNG) + "-buffer-metrics"
	fluentdServiceName                = loggingName + "-" + string(fluentd) + "-metrics"
	fluentdBufferMetricsServiceName   = loggingName + "-" + string(fluentd) + "-buffer-metrics"
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

func metricsEnabled() *v1beta1.Metrics {
	return &v1beta1.Metrics{Enabled: new(true), ServiceMonitor: true}
}

func fluentbitWithMetrics() *v1beta1.FluentbitSpec {
	return &v1beta1.FluentbitSpec{
		Metrics:             metricsEnabled(),
		BufferVolumeMetrics: metricsEnabled(),
		ConfigHotReload:     &v1beta1.HotReload{Image: image.ConfigReloader().Spec()},
		BufferVolumeImage:   image.NodeExporter().Spec(),
	}
}

// The suite runs syslog-ng and fluentd in turn rather than together, since a
// Logging carries one aggregator at a time. Fluent Bit is in both halves
// because its metrics have to survive the aggregator being replaced.
func syslogNGLogging() *v1beta1.Logging {
	return &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: loggingName, Namespace: ns},
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: ns,
			FluentbitSpec:    fluentbitWithMetrics(),
			SyslogNGSpec: &v1beta1.SyslogNGSpec{
				ConfigReloadImage:        image.SyslogNGReloader().Basic(),
				BufferVolumeMetricsImage: image.NodeExporter().Basic(),
				Metrics:                  metricsEnabled(),
				BufferVolumeMetrics:      &v1beta1.BufferMetrics{Metrics: *metricsEnabled()},
			},
		},
	}
}

func fluentdLogging() *v1beta1.Logging {
	return &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{Name: loggingName, Namespace: ns},
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: ns,
			FluentbitSpec:    fluentbitWithMetrics(),
			FluentdSpec: &v1beta1.FluentdSpec{
				Image:               image.Fluentd().Spec(),
				ConfigReloaderImage: image.ConfigReloader().Spec(),
				BufferVolumeImage:   image.NodeExporter().Spec(),
				Metrics:             metricsEnabled(),
				BufferVolumeMetrics: metricsEnabled(),
			},
		},
	}
}

func TestLoggingMetrics_Monitoring(t *testing.T) {
	env := harness.New(t).
		WithCluster("logging-metrics-monitoring").
		WithRelease(release).
		WithControlNamespace(ns).
		// Only this suite installs the prometheus-operator CRDs, so only it
		// registers their types.
		WithScheme(v1.AddToScheme).
		Start()

	require.NoError(t, installPrometheusOperator(env))

	logging := syslogNGLogging()
	env.Create(logging)
	env.WaitFor(wait.Pod(ns, loggingName+"-"+string(syslogNG)+"-0"))

	serviceMonitorsSyslogNG := &v1.ServiceMonitorList{}
	require.NoError(t, env.Client.List(env.Ctx, serviceMonitorsSyslogNG))

	mt, err := setupMetricsTester(env)
	require.NoError(t, err)

	mt.mustServe(env, fluentbit)
	mt.mustServe(env, syslogNG)

	require.NoError(t, env.Client.Delete(env.Ctx, logging))

	env.Create(fluentdLogging())
	env.WaitFor(wait.Pod(ns, loggingName+"-"+string(fluentd)+"-0"))

	serviceMonitorsFluentd := &v1.ServiceMonitorList{}
	require.NoError(t, env.Client.List(env.Ctx, serviceMonitorsFluentd))

	mt.mustServe(env, fluentd)

	serviceMonitors := append(serviceMonitorsFluentd.Items, serviceMonitorsSyslogNG.Items...)
	require.NoError(t, checkServiceMonitorAvailability(serviceMonitors))
}

// installPrometheusOperator keeps the chart install rather than a manifest: the
// stack is what the ServiceMonitors are read by, and pinning our own copy of it
// would be a second thing to keep current.
func installPrometheusOperator(env *harness.Env) error {
	manager := helm.New(env.Cluster.KubeConfigFilePath())

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

func setupMetricsTester(env *harness.Env) (metricsTester, error) {
	pod, err := common.SetupCurlPod(env.Ctx, env.Client, ns, "metrics-tester", pollInterval, pollTimeout)
	if err != nil {
		return metricsTester{}, err
	}

	return metricsTester{testPod: pod}, nil
}

func checkServiceMonitorAvailability(serviceMonitors []v1.ServiceMonitor) error {
	if len(serviceMonitors) == 0 {
		return errors.New("no service monitors found")
	}

	expectedServiceMonitors := map[string]bool{
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

// mustServe polls rather than waits once: the endpoint answers before the first
// scrape has populated it, so an empty body is a retry and not a failure. It
// stays a suite-local poll because reading it means curling from inside a pod,
// which a read-only Condition cannot do.
func (mt *metricsTester) mustServe(env *harness.Env, subject loggingResourceName) {
	env.T.Helper()

	require.Eventuallyf(env.T, func() bool {
		rawOut, err := mt.getMetrics(metricServices[subject], env)
		if err != nil {
			env.T.Log(err)
			return false
		}
		if err := mt.validateMetrics(rawOut, subject); err != nil {
			env.T.Log(err)
			return false
		}
		return true
	}, pollTimeout, pollInterval, "%s never served its key metrics", subject)
}

func (mt *metricsTester) getMetrics(endpoint metricsEndpoint, env *harness.Env) ([]byte, error) {
	serviceURL := fmt.Sprintf("http://%s.%s.svc:%d%s",
		endpoint.serviceName,
		ns,
		endpoint.port,
		endpoint.path,
	)
	cmd := common.CmdEnv(exec.Command("kubectl", "exec", mt.testPod.Name, "-n", ns, "--", "curl", serviceURL), env.Cluster)
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
