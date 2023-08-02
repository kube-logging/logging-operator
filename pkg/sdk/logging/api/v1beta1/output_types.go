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
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +name:"OutputSpec"
// +weight:"200"
type _hugoOutputSpec interface{} //nolint:deadcode,unused

// +name:"OutputSpec"
// +version:"v1beta1"
// +description:"OutputSpec defines the desired state of Output"
type _metaOutputSpec interface{} //nolint:deadcode,unused

// OutputSpec defines the desired state of Output
type OutputSpec struct {
	AwsElasticsearchOutputConfig *output.AwsElasticsearchOutputConfig `json:"awsElasticsearch,omitempty"`
	AzureStorage                 *output.AzureStorage                 `json:"azurestorage,omitempty"`
	CloudWatchOutput             *output.CloudWatchOutput             `json:"cloudwatch,omitempty"`
	DatadogOutput                *output.DatadogOutput                `json:"datadog,omitempty"`
	ElasticsearchOutput          *output.ElasticsearchOutput          `json:"elasticsearch,omitempty"`
	FileOutput                   *output.FileOutputConfig             `json:"file,omitempty"`
	ForwardOutput                *output.ForwardOutput                `json:"forward,omitempty"`
	GELFOutputConfig             *output.GELFOutputConfig             `json:"gelf,omitempty"`
	GCSOutput                    *output.GCSOutput                    `json:"gcs,omitempty"`
	HTTPOutput                   *output.HTTPOutputConfig             `json:"http,omitempty"`
	KafkaOutputConfig            *output.KafkaOutputConfig            `json:"kafka,omitempty"`
	KinesisStreamOutputConfig    *output.KinesisStreamOutputConfig    `json:"kinesisStream,omitempty"`
	LogDNAOutput                 *output.LogDNAOutput                 `json:"logdna,omitempty"`
	LoggingRef                   string                               `json:"loggingRef,omitempty"`
	LogZOutput                   *output.LogZOutput                   `json:"logz,omitempty"`
	LokiOutput                   *output.LokiOutput                   `json:"loki,omitempty"`
	MattermostOutputConfig       *output.MattermostOutputConfig       `json:"mattermost,omitempty"`
	NewRelicOutputConfig         *output.NewRelicOutputConfig         `json:"newrelic,omitempty"`
	NullOutputConfig             *output.NullOutputConfig             `json:"nullout,omitempty"`
	OSSOutput                    *output.OSSOutput                    `json:"oss,omitempty"`
	OpenSearchOutput             *output.OpenSearchOutput             `json:"opensearch,omitempty"`
	RedisOutputConfig            *output.RedisOutputConfig            `json:"redis,omitempty"`
	RelabelOutputConfig          *output.RelabelOutputConfig          `json:"relabel,omitempty"`
	S3OutputConfig               *output.S3OutputConfig               `json:"s3,omitempty"`
	SplunkHecOutput              *output.SplunkHecOutput              `json:"splunkHec,omitempty"`
	SQSOutputConfig              *output.SQSOutputConfig              `json:"sqs,omitempty"`
	SumologicOutput              *output.SumologicOutput              `json:"sumologic,omitempty"`
	SyslogOutputConfig           *output.SyslogOutputConfig           `json:"syslog,omitempty"`
}

// OutputStatus defines the observed state of Output
type OutputStatus struct {
	Active        *bool    `json:"active,omitempty"`
	Problems      []string `json:"problems,omitempty"`
	ProblemsCount int      `json:"problemsCount,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories=logging-all
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Active",type="boolean",JSONPath=".status.active",description="Is the output active?"
// +kubebuilder:printcolumn:name="Problems",type="integer",JSONPath=".status.problemsCount",description="Number of problems"
// +kubebuilder:storageversion

// Output is the Schema for the outputs API
type Output struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OutputSpec   `json:"spec,omitempty"`
	Status            OutputStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OutputList contains a list of Output
type OutputList struct {
	Items           []Output `json:"items"`
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Output{}, &OutputList{})
}
