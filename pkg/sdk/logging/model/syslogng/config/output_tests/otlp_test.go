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

package test

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/cisco-open/operator-tools/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
)

func TestOpenTelemetryOutputTable(t *testing.T) {
	var tests = []struct {
		name   string
		output v1beta1.SyslogNGOutput
		config string
	}{
		{
			name: "test_minimal",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-otlp-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					OpenTelemetry: &output.OpenTelemetryOutput{
						URL: "otlp-server",
					},
				},
			},
			config: heredoc.Doc(`
				destination "output_default_test-otlp-out" {
					opentelemetry(url("otlp-server"));
				};
			`),
		},
		{
			name: "test_full",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-otlp-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					OpenTelemetry: &output.OpenTelemetryOutput{
						URL: "otlp-server",
						Auth: &output.Auth{
							Insecure: &output.Insecure{},
						},
						DiskBuffer: &output.DiskBuffer{
							Reliable: true,
						},
						Batch: output.Batch{
							BatchLines: 1,
						},
						Compression: utils.BoolPointer(true),
						ChannelArgs: filter.ArrowMap{
							"a": "b",
						},
					},
				},
			},
			config: heredoc.Doc(`
				destination "output_default_test-otlp-out" {
					opentelemetry(url("otlp-server") auth(insecure()) disk_buffer(disk_buf_size(0) reliable(yes)) batch-lines(1) compression(yes) channel_args(
						"a" => "b"
					));
				};
			`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CheckConfigForOutput(t, tt.output, tt.config)
		})
	}
}
