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

package metrics

import (
	"flag"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "sidecar_reloader"

var (
	listenAddress = flag.String("web.listen-address", ":9533", "Address to listen on for web interface and telemetry.")
	metricPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")

	LastReloadError = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "config_reloader_last_reload_error",
		Help:      "Whether the last reload resulted in an error (1 for error, 0 for success)",
	}, []string{"webhook"})
	RequestDuration = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "config_reloader_last_request_duration_seconds",
		Help:      "Duration of last webhook request",
	}, []string{"webhook"})
	SuccessReloads = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "config_reloader_success_reloads_total",
		Help:      "Total success reload calls",
	}, []string{"webhook"})
	RequestErrorsByReason = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "config_reloader_request_errors_total",
		Help:      "Total request errors by reason",
	}, []string{"webhook", "reason"})
	WatcherErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "config_reloader_watcher_errors_total",
		Help:      "Total filesystem watcher errors",
	})
	RequestsByStatusCode = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "config_reloader_requests_total",
		Help:      "Total requests by response status code",
	}, []string{"webhook", "status_code"})
)

func init() {
	prometheus.MustRegister(LastReloadError)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(SuccessReloads)
	prometheus.MustRegister(RequestErrorsByReason)
	prometheus.MustRegister(WatcherErrors)
	prometheus.MustRegister(RequestsByStatusCode)
}

func Run() error {
	return serverMetrics(*listenAddress, *metricPath)
}

func serverMetrics(listenAddress, metricsPath string) error {
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>ConfigMap Reload Metrics</title></head>
			<body>
			<h1>ConfigMap Reload</h1>
			<p><a href='` + metricsPath + `'>Metrics</a></p>
			</body>
			</html>
		`))
	})
	return http.ListenAndServe(listenAddress, nil)
}

func SetFailureMetrics(h, reason string) {
	RequestErrorsByReason.WithLabelValues(h, reason).Inc()
	LastReloadError.WithLabelValues(h).Set(1.0)
}

func SetSuccessMetrics(h string, begun time.Time) {
	RequestDuration.WithLabelValues(h).Set(time.Since(begun).Seconds())
	SuccessReloads.WithLabelValues(h).Inc()
	LastReloadError.WithLabelValues(h).Set(0.0)
}
