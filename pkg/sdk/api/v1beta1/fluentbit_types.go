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
	"strconv"

	"github.com/banzaicloud/operator-tools/pkg/volume"
	corev1 "k8s.io/api/core/v1"
)

// +name:"FluentbitSpec"
// +weight:"200"
type _hugoFluentbitSpec interface{}

// +name:"FluentbitSpec"
// +version:"v1beta1"
// +description:"FluentbitSpec defines the desired state of Fluentbit"
type _metaFluentbitSpec interface{}

// +kubebuilder:object:generate=true

// FluentbitSpec defines the desired state of Fluentbit
type FluentbitSpec struct {
	Annotations map[string]string `json:"annotations,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Image       ImageSpec         `json:"image,omitempty"`
	TLS         FluentbitTLS      `json:"tls,omitempty"`
	TargetHost  string            `json:"targetHost,omitempty"`
	TargetPort  int32             `json:"targetPort,omitempty"`
	// Set the flush time in seconds.nanoseconds. The engine loop uses a Flush timeout to define when is required to flush the records ingested by input plugins through the defined output plugins. (default: 1)
	Flush int32 `json:"flush,omitempty"  plugin:"default:1"`
	// Set the grace time in seconds as Integer value. The engine loop uses a Grace timeout to define wait time on exit (default: 5)
	Grace int32 `json:"grace,omitempty" plugin:"default:5"`
	// Set the logging verbosity level. Allowed values are: error, warn, info, debug and trace. Values are accumulative, e.g: if 'debug' is set, it will include error, warning, info and debug.  Note that trace mode is only available if Fluent Bit was built with the WITH_TRACE option enabled. (default: info)
	LogLevel string `json:"logLevel,omitempty" plugin:"default:info"`
	// Set the coroutines stack size in bytes. The value must be greater than the page size of the running system. Don't set too small value (say 4096), or coroutine threads can overrun the stack buffer.
	//Do not change the default value of this parameter unless you know what you are doing. (default: 24576)
	CoroStackSize int32                       `json:"coroStackSize,omitempty" plugin:"default:24576"`
	Resources     corev1.ResourceRequirements `json:"resources,omitempty"`
	Tolerations   []corev1.Toleration         `json:"tolerations,omitempty"`
	NodeSelector  map[string]string           `json:"nodeSelector,omitempty"`
	Affinity      *corev1.Affinity            `json:"affinity,omitempty"`
	Metrics       *Metrics                    `json:"metrics,omitempty"`
	Security      *Security                   `json:"security,omitempty"`
	// +docLink:"volume.KubernetesVolume,https://github.com/banzaicloud/operator-tools/tree/master/docs/types"
	PositionDB volume.KubernetesVolume `json:"positiondb,omitempty"`
	// Deprecated, use positiondb
	PosisionDBLegacy  *volume.KubernetesVolume `json:"position_db,omitempty"`
	MountPath         string                   `json:"mountPath,omitempty"`
	ExtraVolumeMounts []VolumeMount            `json:"extraVolumeMounts,omitempty"`
	InputTail         InputTail                `json:"inputTail,omitempty"`
	FilterAws         *FilterAws               `json:"filterAws,omitempty"`
	// Deprecated, use inputTail.parser
	Parser string `json:"parser,omitempty"`
	// Parameters for Kubernetes metadata filter
	FilterKubernetes FilterKubernetes `json:"filterKubernetes,omitempty"`
	// Disable Kubernetes metadata filter
	DisableKubernetesFilter *bool         `json:"disableKubernetesFilter,omitempty"`
	BufferStorage           BufferStorage `json:"bufferStorage,omitempty"`
	// +docLink:"volume.KubernetesVolume,https://github.com/banzaicloud/operator-tools/tree/master/docs/types"
	BufferStorageVolume  volume.KubernetesVolume `json:"bufferStorageVolume,omitempty"`
	CustomConfigSecret   string                  `json:"customConfigSecret,omitempty"`
	PodPriorityClassName string                  `json:"podPriorityClassName,omitempty"`
	LivenessProbe        *corev1.Probe           `json:"livenessProbe,omitempty"`
	LivenessDefaultCheck bool                    `json:"livenessDefaultCheck,omitempty"`
	ReadinessProbe       *corev1.Probe           `json:"readinessProbe,omitempty"`
	Network              *FluentbitNetwork       `json:"network,omitempty"`
	ForwardOptions       *ForwardOptions         `json:"forwardOptions,omitempty"`
	EnableUpstream       bool                    `json:"enableUpstream,omitempty"`
}

// +kubebuilder:object:generate=true

// FluentbitTLS defines the TLS configs
type FluentbitTLS struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName,omitempty"`
	SharedKey  string `json:"sharedKey,omitempty"`
}

// FluentbitNetwork defines network configuration for fluentbit
type FluentbitNetwork struct {
	// Sets the timeout for connecting to an upstream (default: 10)
	ConnectTimeout *uint32 `json:"connectTimeout,omitempty"`
	// Whether or not TCP keepalive is used for the upstream connection (default: true)
	Keepalive *bool `json:"keepalive,omitempty"`
	// How long in seconds a TCP keepalive connection can be idle before being recycled (default: 30)
	KeepaliveIdleTimeout *uint32 `json:"keepaliveIdleTimeout,omitempty"`
	// How many times a TCP keepalive connection can be used before being recycled (default: 0, disabled)
	KeepaliveMaxRecycle *uint32 `json:"keepaliveMaxRecycle,omitempty"`
}

// GetPrometheusPortFromAnnotation gets the port value from annotation
func (spec FluentbitSpec) GetPrometheusPortFromAnnotation() int32 {
	var err error
	var port int64
	if spec.Annotations != nil {
		port, err = strconv.ParseInt(spec.Annotations["prometheus.io/port"], 10, 32)
		if err != nil {
			panic(err)
		}
	}
	return int32(port)
}

// BufferStorage is the Service Section Configuration of fluent-bit
type BufferStorage struct {
	// Set an optional location in the file system to store streams and chunks of data. If this parameter is not set, Input plugins can only use in-memory buffering.
	StoragePath string `json:"storage.path,omitempty"`
	// Configure the synchronization mode used to store the data into the file system. It can take the values normal or full. (default:normal)
	StorageSync string `json:"storage.sync,omitempty"`
	// Enable the data integrity check when writing and reading data from the filesystem. The storage layer uses the CRC32 algorithm. (default:Off)
	StorageChecksum string `json:"storage.checksum,omitempty"`
	// If storage.path is set, Fluent Bit will look for data chunks that were not delivered and are still in the storage layer, these are called backlog data. This option configure a hint of maximum value of memory to use when processing these records. (default:5M)
	StorageBacklogMemLimit string `json:"storage.backlog.mem_limit,omitempty"`
}

// InputTail defines Fluentbit tail input configuration The tail input plugin allows to monitor one or several text files. It has a similar behavior like tail -f shell command.
type InputTail struct {
	// Specify the buffering mechanism to use. It can be memory or filesystem. (default:memory)
	StorageType string `json:"storage.type,omitempty"`
	// Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the Unit Size specification. (default:32k)
	BufferChunkSize string `json:"Buffer_Chunk_Size,omitempty"`
	// Set the limit of the buffer size per monitored file. When a buffer needs to be increased (e.g: very long lines), this value is used to restrict how much the memory buffer can grow. If reading a file exceed this limit, the file is removed from the monitored file list. The value must be according to the Unit Size specification. (default:Buffer_Chunk_Size)
	BufferMaxSize string `json:"Buffer_Max_Size,omitempty"`
	// Pattern specifying a specific log files or multiple ones through the use of common wildcards.
	Path string `json:"Path,omitempty"`
	// If enabled, it appends the name of the monitored file as part of the record. The value assigned becomes the key in the map.
	PathKey string `json:"Path_Key,omitempty"`
	// Set one or multiple shell patterns separated by commas to exclude files matching a certain criteria, e.g: exclude_path=*.gz,*.zip
	ExcludePath string `json:"Exclude_Path,omitempty"`
	// The interval of refreshing the list of watched files in seconds. (default:60)
	RefreshInterval string `json:"Refresh_Interval,omitempty"`
	// Specify the number of extra time in seconds to monitor a file once is rotated in case some pending data is flushed. (default:5)
	RotateWait string `json:"Rotate_Wait,omitempty"`
	// Ignores files that have been last modified before this time in seconds. Supports m,h,d (minutes, hours,days) syntax. Default behavior is to read all specified files.
	IgnoreOlder string `json:"Ignore_Older,omitempty"`
	// When a monitored file reach it buffer capacity due to a very long line (Buffer_Max_Size), the default behavior is to stop monitoring that file. Skip_Long_Lines alter that behavior and instruct Fluent Bit to skip long lines and continue processing other lines that fits into the buffer size. (default:Off)
	SkipLongLines string `json:"Skip_Long_Lines,omitempty"`
	// Specify the database file to keep track of monitored files and offsets.
	DB *string `json:"DB,omitempty"`
	// Set a default synchronization (I/O) method. Values: Extra, Full, Normal, Off. This flag affects how the internal SQLite engine do synchronization to disk, for more details about each option please refer to this section. (default:Full)
	DBSync string `json:"DB_Sync,omitempty"`
	// Set a limit of memory that Tail plugin can use when appending data to the Engine. If the limit is reach, it will be paused; when the data is flushed it resumes.
	MemBufLimit string `json:"Mem_Buf_Limit,omitempty"`
	// Specify the name of a parser to interpret the entry as a structured message.
	Parser string `json:"Parser,omitempty"`
	// When a message is unstructured (no parser applied), it's appended as a string under the key name log. This option allows to define an alternative name for that key. (default:log)
	Key string `json:"Key,omitempty"`
	// Set a tag (with regex-extract fields) that will be placed on lines read.
	Tag string `json:"Tag,omitempty"`
	// Set a regex to extract fields from the file.
	TagRegex string `json:"Tag_Regex,omitempty"`
	// If enabled, the plugin will try to discover multiline messages and use the proper parsers to compose the outgoing messages. Note that when this option is enabled the Parser option is not used. (default:Off)
	Multiline string `json:"Multiline,omitempty"`
	// Wait period time in seconds to process queued multiline messages (default:4)
	MultilineFlush string `json:"Multiline_Flush,omitempty"`
	// Name of the parser that machs the beginning of a multiline message. Note that the regular expression defined in the parser must include a group name (named capture)
	ParserFirstline string `json:"Parser_Firstline,omitempty"`
	// Optional-extra parser to interpret and structure multiline entries. This option can be used to define multiple parsers, e.g: Parser_1 ab1,  Parser_2 ab2, Parser_N abN.
	ParserN []string `json:"Parser_N,omitempty"`
	// If enabled, the plugin will recombine split Docker log lines before passing them to any parser as configured above. This mode cannot be used at the same time as Multiline. (default:Off)
	DockerMode string `json:"Docker_Mode,omitempty"`
	//Wait period time in seconds to flush queued unfinished split lines. (default:4)
	DockerModeFlush string `json:"Docker_Mode_Flush,omitempty"`
}

// FilterKubernetes Fluent Bit Kubernetes Filter allows to enrich your log files with Kubernetes metadata.
type FilterKubernetes struct {
	// Match filtered records (default:kube.*)
	Match string `json:"Match,omitempty" plugin:"default:kubernetes.*"`
	// Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the Unit Size specification. (default:32k)
	BufferSize string `json:"Buffer_Size,omitempty"`
	// API Server end-point (default:https://kubernetes.default.svc:443)
	KubeURL string `json:"Kube_URL,omitempty" plugin:"default:https://kubernetes.default.svc:443"`
	//	CA certificate file (default:/var/run/secrets/kubernetes.io/serviceaccount/ca.crt)
	KubeCAFile string `json:"Kube_CA_File,omitempty" plugin:"default:/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"`
	// Absolute path to scan for certificate files
	KubeCAPath string `json:"Kube_CA_Path,omitempty"`
	// Token file  (default:/var/run/secrets/kubernetes.io/serviceaccount/token)
	KubeTokenFile string `json:"Kube_Token_File,omitempty" plugin:"default:/var/run/secrets/kubernetes.io/serviceaccount/token"`
	// When the source records comes from Tail input plugin, this option allows to specify what's the prefix used in Tail configuration. (default:kube.var.log.containers.)
	KubeTagPrefix string `json:"Kube_Tag_Prefix,omitempty" plugin:"default:kubernetes.var.log.containers"`
	// When enabled, it checks if the log field content is a JSON string map, if so, it append the map fields as part of the log structure. (default:Off)
	MergeLog string `json:"Merge_Log,omitempty" plugin:"default:On"`
	// When Merge_Log is enabled, the filter tries to assume the log field from the incoming message is a JSON string message and make a structured representation of it at the same level of the log field in the map. Now if Merge_Log_Key is set (a string name), all the new structured fields taken from the original log content are inserted under the new key.
	MergeLogKey string `json:"Merge_Log_Key,omitempty"`
	// When Merge_Log is enabled, trim (remove possible \n or \r) field values.  (default:On)
	MergeLogTrim string `json:"Merge_Log_Trim,omitempty"`
	// Optional parser name to specify how to parse the data contained in the log key. Recommended use is for developers or testing only.
	MergeParser string `json:"Merge_Parser,omitempty"`
	// When Keep_Log is disabled, the log field is removed from the incoming message once it has been successfully merged (Merge_Log must be enabled as well). (default:On)
	KeepLog string `json:"Keep_Log,omitempty"`
	// Debug level between 0 (nothing) and 4 (every detail). (default:-1)
	TLSDebug string `json:"tls.debug,omitempty"`
	// When enabled, turns on certificate validation when connecting to the Kubernetes API server. (default:On)
	TLSVerify string `json:"tls.verify,omitempty"`
	// When enabled, the filter reads logs coming in Journald format. (default:Off)
	UseJournal string `json:"Use_Journal,omitempty"`
	// Set an alternative Parser to process record Tag and extract pod_name, namespace_name, container_name and docker_id. The parser must be registered in a parsers file (refer to parser filter-kube-test as an example).
	RegexParser string `json:"Regex_Parser,omitempty"`
	// Allow Kubernetes Pods to suggest a pre-defined Parser (read more about it in Kubernetes Annotations section) (default:Off)
	K8SLoggingParser string `json:"K8S-Logging.Parser,omitempty"`
	// Allow Kubernetes Pods to exclude their logs from the log processor (read more about it in Kubernetes Annotations section). (default:Off)
	K8SLoggingExclude string `json:"K8S-Logging.Exclude,omitempty"`
	// Include Kubernetes resource labels in the extra metadata. (default:On)
	Labels string `json:"Labels,omitempty"`
	// Include Kubernetes resource annotations in the extra metadata. (default:On)
	Annotations string `json:"Annotations,omitempty"`
	// If set, Kubernetes meta-data can be cached/pre-loaded from files in JSON format in this directory, named as namespace-pod.meta
	KubeMetaPreloadCacheDir string `json:"Kube_meta_preload_cache_dir,omitempty"`
	// If set, use dummy-meta data (for test/dev purposes) (default:Off)
	DummyMeta string `json:"Dummy_Meta,omitempty"`
}

// FilterAws The AWS Filter Enriches logs with AWS Metadata.
type FilterAws struct {
	// Specify which version of the instance metadata service to use. Valid values are 'v1' or 'v2' (default).
	ImdsVersion string `json:"imds_version,omitempty" plugin:"default:v2"`
	// The availability zone (default:true).
	AZ *bool `json:"az,omitempty" plugin:"default:true"`
	// The EC2 instance ID. (default:true)
	Ec2InstanceID *bool `json:"ec2_instance_id,omitempty" plugin:"default:true"`
	// The EC2 instance type. (default:false)
	Ec2InstanceType *bool `json:"ec2_instance_type,omitempty" plugin:"default:false"`
	// The EC2 instance private ip. (default:false)
	PrivateIP *bool `json:"private_ip,omitempty" plugin:"default:false"`
	// The EC2 instance image id. (default:false)
	AmiID *bool `json:"ami_id,omitempty" plugin:"default:false"`
	// The account ID for current EC2 instance. (default:false)
	AccountID *bool `json:"account_id,omitempty" plugin:"default:false"`
	// The hostname for current EC2 instance. (default:false)
	Hostname *bool `json:"hostname,omitempty" plugin:"default:false"`
	// The VPC ID for current EC2 instance. (default:false)
	VpcID *bool `json:"vpc_id,omitempty" plugin:"default:false"`
	// Match filtered records (default:*)
	Match string `json:"Match,omitempty" plugin:"default:*"`
}

// VolumeMount defines source and destination folders of a hostPath type pod mount
type VolumeMount struct {
	// Source folder
	// +kubebuilder:validation:Pattern=^/.+$
	Source string `json:"source"`
	// Destination Folder
	// +kubebuilder:validation:Pattern=^/.+$
	Destination string `json:"destination"`
	// Mount Mode
	ReadOnly bool `json:"readOnly,omitempty"`
}

// ForwardOptions defines custom forward output plugin options, see https://docs.fluentbit.io/manual/pipeline/outputs/forward
type ForwardOptions struct {
	TimeAsInteger      bool   `json:"Time_as_Integer,omitempty"`
	SendOptions        bool   `json:"Send_options,omitempty"`
	RequireAckResponse bool   `json:"Require_ack_response,omitempty"`
	Tag                string `json:"Tag,omitempty"`
	RetryLimit         string `json:"Retry_Limit,omitempty"`
}
