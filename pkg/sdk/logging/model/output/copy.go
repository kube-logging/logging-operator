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

package output

import (
	"reflect"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/filter"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +name:"Copy"
// +weight:"200"
type _hugoCopy interface{} //nolint:deadcode,unused

// +name:"Copy"
// +url:"https://docs.fluentd.org/output/copy"
// +version:"more info"
// +description:"The copy output plugin copies events to multiple outputs"
// +status:"GA"
type _metaCopy interface{} //nolint:deadcode,unused

type Store struct {
	S3OutputConfig               *S3OutputConfig               `json:"s3,omitempty"`
	AzureStorage                 *AzureStorage                 `json:"azurestorage,omitempty"`
	GCSOutput                    *GCSOutput                    `json:"gcs,omitempty"`
	OSSOutput                    *OSSOutput                    `json:"oss,omitempty"`
	ElasticsearchOutput          *ElasticsearchOutput          `json:"elasticsearch,omitempty"`
	OpenSearchOutput             *OpenSearchOutput             `json:"opensearch,omitempty"`
	LogZOutput                   *LogZOutput                   `json:"logz,omitempty"`
	LokiOutput                   *LokiOutput                   `json:"loki,omitempty"`
	SumologicOutput              *SumologicOutput              `json:"sumologic,omitempty"`
	DatadogOutput                *DatadogOutput                `json:"datadog,omitempty"`
	ForwardOutput                *ForwardOutput                `json:"forward,omitempty"`
	FileOutput                   *FileOutputConfig             `json:"file,omitempty"`
	NullOutputConfig             *NullOutputConfig             `json:"nullout,omitempty"`
	KafkaOutputConfig            *KafkaOutputConfig            `json:"kafka,omitempty"`
	CloudWatchOutput             *CloudWatchOutput             `json:"cloudwatch,omitempty"`
	KinesisStreamOutputConfig    *KinesisStreamOutputConfig    `json:"kinesisStream,omitempty"`
	LogDNAOutput                 *LogDNAOutput                 `json:"logdna,omitempty"`
	NewRelicOutputConfig         *NewRelicOutputConfig         `json:"newrelic,omitempty"`
	SplunkHecOutput              *SplunkHecOutput              `json:"splunkHec,omitempty"`
	HTTPOutput                   *HTTPOutputConfig             `json:"http,omitempty"`
	AwsElasticsearchOutputConfig *AwsElasticsearchOutputConfig `json:"awsElasticsearch,omitempty"`
	RedisOutputConfig            *RedisOutputConfig            `json:"redis,omitempty"`
	SyslogOutputConfig           *SyslogOutputConfig           `json:"syslog,omitempty"`
	GELFOutputConfig             *GELFOutputConfig             `json:"gelf,omitempty"`
	SQSOutputConfig              *SQSOutputConfig              `json:"sqs,omitempty"`
	PrometheusConfig             *filter.PrometheusConfig      `json:"prometheus,omitempty"`
}

func (in *Store) DeepCopyInto(out *Store) {
	*out = *in
}

type StoreDirectiveConverter interface {
	ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error)
}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type CopyOutputConfig struct {
	// Chooses how to pass the events to <store> plugins. (default:no_copy)
	CopyMode string `json:"copy_mode,omitempty"`

	// Specifies the storage destinations.
	Stores []Store `json:"stores"`
}

func (c *CopyOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	const pluginType = "copy"
	copy := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        id,
		},
	}
	if _, err := types.NewStructToStringMapper(secretLoader).StringsMap(c); err != nil {
		return nil, err
	}
	for _, store := range c.Stores {
		v := reflect.ValueOf(store)
		var converters []StoreDirectiveConverter
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
				if converter, ok := v.Field(i).Interface().(StoreDirectiveConverter); ok {
					converters = append(converters, converter)
				}
			}
		}
		switch len(converters) {
		case 0:
			return nil, errors.New("no plugin config available for store")
		case 1:
			if meta, err := converters[0].ToDirective(secretLoader, ""); err != nil {
				return nil, err
			} else {
				meta.GetPluginMeta().Directive = "store"
				meta.GetPluginMeta().Tag = ""
				copy.SubDirectives = append(copy.SubDirectives, meta)
			}
		default:
			return nil, errors.Errorf("more then one plugin config is not allowed for a store")
		}
	}
	return copy, nil
}
