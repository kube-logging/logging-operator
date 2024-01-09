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
	"github.com/cisco-open/operator-tools/pkg/typeoverride"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
)

// +name:"SyslogNGSpec"
// +weight:"200"
type _hugoSyslogNGSpec interface{} //nolint:deadcode,unused

// +name:"SyslogNGSpec"
// +version:"v1beta1"
// +description:"SyslogNGSpec defines the desired state of SyslogNG"
type _metaSyslogNGSpec interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true

// SyslogNGSpec defines the desired state of SyslogNG
type SyslogNGSpec struct {
	TLS                                 SyslogNGTLS                  `json:"tls,omitempty"`
	ReadinessDefaultCheck               ReadinessDefaultCheck        `json:"readinessDefaultCheck,omitempty"`
	SkipRBACCreate                      bool                         `json:"skipRBACCreate,omitempty"`
	StatefulSetOverrides                *typeoverride.StatefulSet    `json:"statefulSet,omitempty"`
	ServiceOverrides                    *typeoverride.Service        `json:"service,omitempty"`
	ServiceAccountOverrides             *typeoverride.ServiceAccount `json:"serviceAccount,omitempty"`
	ConfigCheckPodOverrides             *typeoverride.PodSpec        `json:"configCheckPod,omitempty"`
	Metrics                             *Metrics                     `json:"metrics,omitempty"`
	MetricsServiceOverrides             *typeoverride.Service        `json:"metricsService,omitempty"`
	BufferVolumeMetrics                 *BufferMetrics               `json:"bufferVolumeMetrics,omitempty"`
	BufferVolumeMetricsServiceOverrides *typeoverride.Service        `json:"bufferVolumeMetricsService,omitempty"`
	GlobalOptions                       *GlobalOptions               `json:"globalOptions,omitempty"`
	JSONKeyPrefix                       string                       `json:"jsonKeyPrefix,omitempty"`
	JSONKeyDelimiter                    string                       `json:"jsonKeyDelim,omitempty"`
	// Available in Logging operator version 4.5 and later.
	// Parses date automatically from the timestamp registered by the container runtime.
	// Note: `jsonKeyPrefix` and `jsonKeyDelim` are respected.
	SourceDateParser *SourceDateParser     `json:"sourceDateParser,omitempty"`
	// Available in Logging operator version 4.5 and later.
	// Set the maximum number of connections for the source. For details, see [documentation of the AxoSyslog syslog-ng distribution](https://axoflow.com/docs/axosyslog-core/chapter-routing-filters/concepts-flow-control/configuring-flow-control/).
	MaxConnections   int                   `json:"maxConnections,omitempty"`
	LogIWSize        int                   `json:"logIWSize,omitempty"`
	// Available in Logging operator version 4.5 and later.
	// Create [custom log metrics for sources and outputs]({{< relref "/docs/examples/custom-syslog-ng-metrics.md" >}}).
	SourceMetrics    []filter.MetricsProbe `json:"sourceMetrics,omitempty"`
	// TODO: option to turn on/off buffer volume PVC
}

//
/* 
Available in Logging operator version 4.5 and later.

Parses date automatically from the timestamp registered by the container runtime.
Note: `jsonKeyPrefix` and `jsonKeyDelim` are respected.
It is disabled by default, but if enabled, then the default settings parse the timestamp written by the container runtime and parsed by Fluent Bit using the `cri` or the `docker` parser.
*/
type SourceDateParser struct {
	// Default: "%FT%T.%f%z"
	Format *string `json:"format,omitempty"`
	// Default(depending on JSONKeyPrefix): "${json.time}"
	Template *string `json:"template,omitempty"`
}

// +name:"SyslogNGConfig"
// +weight:"200"
type _hugoSyslogNGConfig interface{} //nolint:deadcode,unused

// +name:"SyslogNG"
// +version:"v1beta1"
// +description:"SyslogNG is a reference to the desired SyslogNG state"
type _metaSyslogNGConfig interface{} //nolint:deadcode,unused

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

type SyslogNGConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyslogNGSpec         `json:"spec,omitempty"`
	Status SyslogNGConfigStatus `json:"status,omitempty"`
}

// SyslogNGConfigStatus
type SyslogNGConfigStatus struct {
	Logging       string   `json:"logging,omitempty"`
	Active        *bool    `json:"active,omitempty"`
	Problems      []string `json:"problems,omitempty"`
	ProblemsCount int      `json:"problemsCount,omitempty"`
}

// +kubebuilder:object:root=true

// SyslogNGConfigList
type SyslogNGConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyslogNGConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyslogNGConfig{}, &SyslogNGConfigList{})
}

// +kubebuilder:object:generate=true

// SyslogNGTLS defines the TLS configs
type SyslogNGTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName,omitempty"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

type GlobalOptions struct {
	// Deprecated. Use stats/level from 4.1+
	StatsLevel *int `json:"stats_level,omitempty"`
	// Deprecated. Use stats/freq from 4.1+
	StatsFreq *int `json:"stats_freq,omitempty"`
	// See the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-global-options/reference-options/#global-option-stats).
	Stats *Stats `json:"stats,omitempty"`
	// See the [AxoSyslog Core documentation](https://axoflow.com/docs/axosyslog-core/chapter-global-options/reference-options/#global-options-log-level).
	LogLevel *string `json:"log_level,omitempty"`
}

type Stats struct {
	Level *int `json:"level,omitempty"`
	Freq  *int `json:"freq,omitempty"`
}

func (s *SyslogNGSpec) SetDefaults() {
	if s != nil {
		// if s.MaxConnections == 0 {
		// 	max connections is now configured dynamically if not set
		// }
		if s.Metrics != nil {
			if s.Metrics.Path == "" {
				s.Metrics.Path = "/metrics"
			}
			if s.Metrics.Port == 0 {
				s.Metrics.Port = 9577
			}
			if s.Metrics.Timeout == "" {
				s.Metrics.Timeout = "5s"
			}
			if s.Metrics.Interval == "" {
				s.Metrics.Interval = "15s"
			}
		}
	}
}
