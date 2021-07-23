// Copyright © 2021 Banzai Cloud
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
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/banzaicloud/operator-tools/pkg/secret"
)

// +name:"SQS"
// +weight:"200"
type _hugoSQS interface{}

// +kubebuilder:object:generate=true
// +docName:"[SQS Output](https://github.com/ixixi/fluent-plugin-sqs)"
// Fluentd output plugin for SQS.
type _docSQS interface{}

// +name:"SQS"
// +url:"https://github.com/ixixi/fluent-plugin-sqs"
// +version:"v2.1.0"
// +description:"Output plugin writes fluent-events as queue messages to Amazon SQS"
// +status:"Testing"
type _metaSQS interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type SQSOutputConfig struct {
	// SQS queue name
	QueueName string `json:"queue_name"`
	// Create SQS queue (default: false)
	CreateQueue *bool `json:"create_queue,omitempty"`
	// AWS region (default: us-east-1)
	Region string `json:"region,omitempty"`
	// +docLink:"Buffer,../buffer/"
	Buffer *Buffer `json:"buffer,omitempty"`
}

//
// #### Example `SQS` output configurations
// ```
//apiVersion: logging.banzaicloud.io/v1beta1
//kind: Output
//metadata:
//  name: sqs-output-sample
//spec:
//  sqs:
//    queue_name: some-aws-sqs-queue
//    create_queue: false
//    region: us-east-2
//    buffer:
//      flush_thread_count: 8
//      flush_interval: 5s
//      chunk_limit_size: 8M
//      queue_limit_length: 512
//      retry_max_interval: 30
//      retry_forever: true
// ```
//
// #### Fluentd Config Result
// ```
//  <match **>
//      @type sqs
//      @id test_sqs
//      queue_name some-aws-sqs-queue
//      create_queue false
//      region us-east-2
//      <buffer tag,time>
//        @type file
//        path /buffers/test_file.*.buffer
//    flush_thread_count 8
//    flush_interval 5s
//    chunk_limit_size 8M
//    queue_limit_length 512
//    retry_max_interval 30
//    retry_forever true
//      </buffer>
//  </match>
// ```
type _expSQS interface{}

func (s *SQSOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "sqs"
	sqs := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        id,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(s); err != nil {
		return nil, err
	} else {
		sqs.Params = params
	}
	if s.Buffer == nil {
		s.Buffer = &Buffer{}
	}
	if buffer, err := s.Buffer.ToDirective(secretLoader, id); err != nil {
		return nil, err
	} else {
		sqs.SubDirectives = append(sqs.SubDirectives, buffer)
	}

	return sqs, nil
}
