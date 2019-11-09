// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/secret"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
)

// +docName:"Google Cloud Logging for Fluentd"
//  More info at https://cloud.google.com/logging/docs/agent/configuration#cloud-fluentd-config
//>Example: [Google Cloud Logging Output Deployment](../../../docs/example-gcl.md)
//
// #### Example output configurations
// ```
// spec:
//  google_cloud:
//    num_threads: 8
//    use_grpc: true
//    partial_success: true
//    autoformat_stackdriver_trace: true
//    buffer:
//      timekey: 10m
//      timekey_wait: 30s
//      timekey_use_utc: true*/
// ```
type _docgcl interface{}

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type GclOutputConfig struct {
	// The number of simultaneous log flushes that can be processed by the output plugin.
	NumThreads int `json:"num_threads,omitempty"`
	// Whether to use gRPC instead of REST/JSON to communicate to the Logging API. With gRPC enabled, CPU usage will typically be lower. (default: true)
	UseGrpc bool `json:"use_grpc,omitempty"`
	// Whether to support partial success for logs ingestion. If true, invalid log entries in a full set are dropped, and valid log entries are successfully ingested into the Logging API. If false, the full set would be dropped if it contained any invalid log entries.  (default: true)
	PartialSuccess bool `json:"partial_success,omitempty"`
	// When set to true, the trace will be reformatted if the value of structured payload field logging.googleapis.com/trace matches ResourceTrace traceId format. Details of the autoformatting can be found under Special fields in structured payloads. (default: true)
	AutoformatStackdriverTrace bool `json:"autoformat_stackdriver_trace,omitempty"`
	// +docLink:"Buffer,./buffer.md"
	Buffer *Buffer `json:"buffer,omitempty"`
	// +docLink:"Format,./format.md"
	Format *Format `json:"format,omitempty"`
}

func (c *GclOutputConfig) ToDirective(secretLoader secret.SecretLoader, id string) (types.Directive, error) {
	pluginType := "google_cloud"
	pluginID := id + "_" + pluginType
	gcl := &types.OutputPlugin{
		PluginMeta: types.PluginMeta{
			Type:      pluginType,
			Directive: "match",
			Tag:       "**",
			Id:        pluginID,
		},
	}
	if params, err := types.NewStructToStringMapper(secretLoader).StringsMap(c); err != nil {
		return nil, err
	} else {
		gcl.Params = params
	}
	if c.Buffer != nil {
		if buffer, err := c.Buffer.ToDirective(secretLoader, pluginID); err != nil {
			return nil, err
		} else {
			gcl.SubDirectives = append(gcl.SubDirectives, buffer)
		}
	}
	if c.Format != nil {
		if format, err := c.Format.ToDirective(secretLoader, ""); err != nil {
			return nil, err
		} else {
			gcl.SubDirectives = append(gcl.SubDirectives, format)
		}
	}
	return gcl, nil
}
