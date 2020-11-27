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
)

// +name:"Common"
// +weight:"200"
type _hugoCommon interface{}

// +name:"Common"
// +version:"v1beta1"
// +description:"ImageSpec Metrics Security"
type _metaCommon interface{}

const (
	HostPath = "/opt/logging-operator/%s/%s"
)

// ImageSpec struct hold information about image specification
type ImageSpec struct {
	Repository       string                        `json:"repository,omitempty"`
	Tag              string                        `json:"tag,omitempty"`
	PullPolicy       string                        `json:"pullPolicy,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// Metrics defines the service monitor endpoints
type Metrics struct {
	Interval              string               `json:"interval,omitempty"`
	Timeout               string               `json:"timeout,omitempty"`
	Port                  int32                `json:"port,omitempty"`
	Path                  string               `json:"path,omitempty"`
	ServiceMonitor        bool                 `json:"serviceMonitor,omitempty"`
	ServiceMonitorConfig  ServiceMonitorConfig `json:"serviceMonitorConfig,omitempty"`
	PrometheusAnnotations bool                 `json:"prometheusAnnotations,omitempty"`
}

// ServiceMonitorConfig defines the ServiceMonitor properties
type ServiceMonitorConfig struct {
	AdditionalLabels   map[string]string   `json:"additionalLabels,omitempty"`
	HonorLabels        bool                `json:"honorLabels,omitempty"`
	Relabelings        []*v1.RelabelConfig `json:"relabelings,omitempty"`
	MetricsRelabelings []*v1.RelabelConfig `json:"metricRelabelings,omitempty"`
}

// Security defines Fluentd, Fluentbit deployment security properties
type Security struct {
	ServiceAccount               string                     `json:"serviceAccount,omitempty"`
	RoleBasedAccessControlCreate *bool                      `json:"roleBasedAccessControlCreate,omitempty"`
	PodSecurityPolicyCreate      bool                       `json:"podSecurityPolicyCreate,omitempty"`
	SecurityContext              *corev1.SecurityContext    `json:"securityContext,omitempty"`
	PodSecurityContext           *corev1.PodSecurityContext `json:"podSecurityContext,omitempty"`
}
