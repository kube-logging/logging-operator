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

package v1beta1

import (
	v1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +name:"Common"
// +weight:"200"
type _hugoCommon interface{} //nolint:deadcode,unused

// +name:"Common"
// +version:"v1beta1"
// +description:"ImageSpec Metrics Security"
type _metaCommon interface{} //nolint:deadcode,unused

const (
	HostPath = "/opt/logging-operator/%s/%s"
)

// BasicImageSpec struct hold basic information about image specification
type BasicImageSpec struct {
	Repository string `json:"repository,omitempty"`
	Tag        string `json:"tag,omitempty"`
}

func (s BasicImageSpec) RepositoryWithTag() string {
	return RepositoryWithTag(s.Repository, s.Tag)
}

// ImageSpec struct hold information about image specification
type ImageSpec struct {
	Repository       string                        `json:"repository,omitempty"`
	Tag              string                        `json:"tag,omitempty"`
	PullPolicy       string                        `json:"pullPolicy,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

func (s ImageSpec) RepositoryWithTag() string {
	return RepositoryWithTag(s.Repository, s.Tag)
}

func RepositoryWithTag(repository, tag string) string {
	res := repository
	if tag != "" {
		res += ":" + tag
	}
	return res
}

// Metrics defines the service monitor endpoints
type Metrics struct {
	Interval                string                    `json:"interval,omitempty"`
	Timeout                 string                    `json:"timeout,omitempty"`
	Port                    int32                     `json:"port,omitempty"`
	Path                    string                    `json:"path,omitempty"`
	ServiceMonitor          bool                      `json:"serviceMonitor,omitempty"`
	ServiceMonitorConfig    ServiceMonitorConfig      `json:"serviceMonitorConfig,omitempty"`
	PrometheusAnnotations   bool                      `json:"prometheusAnnotations,omitempty"`
	PrometheusRules         bool                      `json:"prometheusRules,omitempty"`
	PrometheusRulesOverride []PrometheusRulesOverride `json:"prometheusRulesOverride,omitempty"`
}

type PrometheusRulesOverride struct {
	// Name of the time series to output to. Must be a valid metric name.
	// Only one of `record` and `alert` must be set.
	Record string `json:"record,omitempty"`
	// Name of the alert. Must be a valid label value.
	// Only one of `record` and `alert` must be set.
	Alert string `json:"alert,omitempty"`
	// PromQL expression to evaluate.
	Expr *intstr.IntOrString `json:"expr,omitempty"`
	// Alerts are considered firing once they have been returned for this long.
	// +optional
	For *v1.Duration `json:"for,omitempty"`
	// KeepFiringFor defines how long an alert will continue firing after the condition that triggered it has cleared.
	// +optional
	KeepFiringFor *v1.NonEmptyDuration `json:"keep_firing_for,omitempty"`
	// Labels to add or overwrite.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations to add to each alert.
	// Only valid for alerting rules.
	Annotations map[string]string `json:"annotations,omitempty"`
}

func (o PrometheusRulesOverride) ListOverride(listOfRules []v1.Rule) []v1.Rule {
	var updatedRules []v1.Rule
	for _, rule := range listOfRules {
		if (o.Record != "" && o.Record == rule.Record) || (o.Alert != "" && o.Alert == rule.Alert) {
			updatedRule := o.Override(&rule)
			updatedRules = append(updatedRules, *updatedRule)
		} else {
			updatedRules = append(updatedRules, rule)
		}
	}

	return updatedRules
}

func (o PrometheusRulesOverride) Override(rule *v1.Rule) *v1.Rule {
	updatedRule := rule.DeepCopy()
	if o.Expr != nil {
		updatedRule.Expr = *o.Expr
	}
	if o.For != nil {
		updatedRule.For = o.For
	}
	if o.KeepFiringFor != nil {
		updatedRule.KeepFiringFor = o.KeepFiringFor
	}
	if o.Labels != nil {
		updatedRule.Labels = o.Labels
	}
	if o.Annotations != nil {
		updatedRule.Annotations = o.Annotations
	}
	return updatedRule
}

// BufferMetrics defines the service monitor endpoints
type BufferMetrics struct {
	Metrics   `json:",inline"`
	MountName string `json:"mount_name,omitempty"`
}

// ServiceMonitorConfig defines the ServiceMonitor properties
type ServiceMonitorConfig struct {
	AdditionalLabels   map[string]string  `json:"additionalLabels,omitempty"`
	HonorLabels        bool               `json:"honorLabels,omitempty"`
	Relabelings        []v1.RelabelConfig `json:"relabelings,omitempty"`
	MetricsRelabelings []v1.RelabelConfig `json:"metricRelabelings,omitempty"`
	Scheme             string             `json:"scheme,omitempty"`
	TLSConfig          *v1.TLSConfig      `json:"tlsConfig,omitempty"`
}

// Security defines Fluentd, FluentbitAgent deployment security properties
type Security struct {
	ServiceAccount               string `json:"serviceAccount,omitempty"`
	RoleBasedAccessControlCreate *bool  `json:"roleBasedAccessControlCreate,omitempty"`
	// Warning: this is not supported anymore and does nothing
	PodSecurityPolicyCreate bool                       `json:"podSecurityPolicyCreate,omitempty"`
	SecurityContext         *corev1.SecurityContext    `json:"securityContext,omitempty"`
	PodSecurityContext      *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
	CreateOpenShiftSCC      *bool                      `json:"createOpenShiftSCC,omitempty"`
}

// ReadinessDefaultCheck Enable default readiness checks
type ReadinessDefaultCheck struct {
	// Enable default Readiness check it'll fail if the buffer volume free space exceeds the `readinessDefaultThreshold` percentage (90%).
	BufferFreeSpace          bool  `json:"bufferFreeSpace,omitempty"`
	BufferFreeSpaceThreshold int32 `json:"bufferFreeSpaceThreshold,omitempty"`
	BufferFileNumber         bool  `json:"bufferFileNumber,omitempty"`
	BufferFileNumberMax      int32 `json:"bufferFileNumberMax,omitempty"`
	InitialDelaySeconds      int32 `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds           int32 `json:"timeoutSeconds,omitempty"`
	PeriodSeconds            int32 `json:"periodSeconds,omitempty"`
	SuccessThreshold         int32 `json:"successThreshold,omitempty"`
	FailureThreshold         int32 `json:"failureThreshold,omitempty"`
}
