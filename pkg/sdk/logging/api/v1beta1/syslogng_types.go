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
	MaxConnections                      int                          `json:"maxConnections,omitempty"`
	LogIWSize                           int                          `json:"logIWSize,omitempty"`

	// TODO: option to turn on/off buffer volume PVC
}

// +kubebuilder:object:generate=true

// SyslogNGTLS defines the TLS configs
type SyslogNGTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName,omitempty"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

type GlobalOptions struct {
	StatsLevel *int `json:"stats_level,omitempty"`
	StatsFreq  *int `json:"stats_freq,omitempty"`
}

type FileSource struct {
	//  This parameter assigns a facility value to the messages received from the file source if the message does not specify one. (default:kern)
	DefaultFacility string `json:"default-facility,omitempty"`
	//  This parameter assigns an emergency level to the messages received from the file source if the message does not specify one. For example, default-priority(warning).
	DefaultPriority string `json:"default-priority,omitempty"`
	// Specifies the character set (encoding, for example UTF-8) of messages using the legacy BSD-syslog protocol. To list the available character sets on a host, execute the iconv -l command. For details on how encoding affects the size of the message
	Encoding string `json:"encoding,omitempty"`
	// Indicates that the source should be checked periodically. This is useful for files which always indicate readability, even though no new lines were appended. If this value is higher 	than zero, syslog-ng will not attempt to use poll() on the file, but checks whether the file changed every time the follow-freq() interval (in seconds) has elapsed. Floating-point numbers (for example 1.5) can be used as well. (default:1)
	FollowFreq int `json:"follow-freq,omitempty"`
	// Specifies whether syslog-ng should accept the timestamp received from the sending application or client. If disabled, the time of reception will be used instead. This option can be specified globally, and per-source as well. The local setting of the source overrides the global option if available. (default: true)
	KeepTimestamp bool `json:"keep-timestamp,omitempty"`
	// The maximum number of messages fetched from a source during a single poll loop. The destination queues might fill up before flow-control could stop reading if log-fetch-limit() is too high.(default: 100)
	LogFetchLimit int `json:"log-fetch-limit,omitempty"`
	//  The size of the initial window, this value is used during flow control. Make sure that log-iw-size() is larger than the value of log-fetch-limit().(default: 10000)
	LogIwSize int `json:"log-iw-size,omitempty"`
	//  Maximum length of a message in bytes. This length includes the entire message (the data structure and individual fields). The maximal value that can be set is 268435456 bytes (256MB). For messages using the IETF-syslog message format (RFC5424), the maximal size of the value of an SDATA field is 64kB. (default: 65536)
	LogMsgSize int `json:"log-msg-size,omitempty"`
	// Specifies input padding. Some operating systems (such as HP-UX) pad all messages to block boundary. This option can be used to specify the block size. (HP-UX uses 2048 bytes). The syslog-ng OSE application will pad reads from the associated device to the number of bytes set in pad-size(). Mostly used on HP-UX where /dev/log is a named pipe and every write is padded to 2048 bytes. If pad-size() was given and the incoming message does not fit into pad-size(), syslog-ng will not read anymore from this pipe and displays the following error message: (default: 0)
	PadSize int `json:"pad-size,omitempty"`
	//  Label the messages received from the source with custom tags. Tags must be unique, and enclosed between double quotes. When adding multiple tags, separate them with comma, for example tags("dmz", "router"). This option is available only in syslog-ng 3.1 and later.
	Tags string `json:"tags,omitempty"`
	//  The default timezone for messages read from the source. Applies only if no timezone is specified within the message itself.
	//
	//The timezone can be specified by using the name, for example, time-zone("Europe/Budapest")), or as the timezone offset in +/-HH:MM format, for example, +01:00). On Linux and UNIX platforms, the valid timezone names are listed under the /usr/share/zoneinfo directory.
	TimeZone string `json:"time-zone,omitempty"`
	//  TODO
	WildcardFile string `json:"wildcard-file,omitempty"`
}
