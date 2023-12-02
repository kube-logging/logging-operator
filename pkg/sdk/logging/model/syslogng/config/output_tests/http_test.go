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

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHTTPOutputTable(t *testing.T) {
	var tests = []struct {
		name   string
		output v1beta1.SyslogNGOutput
		config string
	}{
		{
			name: "test_peer_verify_true",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
						TLS: &output.TLS{
							PeerVerify: config.NewTrue(),
						},
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") tls(peer_verify(yes)) batch-lines(2000) workers(3) persist_name("output_default_test-http-out"));
};
`,
		},
		{
			name: "test_peer_verify_false",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
						TLS: &output.TLS{
							PeerVerify: config.NewFalse(),
						},
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") tls(peer_verify(no)) batch-lines(2000) workers(3) persist_name("output_default_test-http-out"));
};
`,
		},
		{
			name: "test_peer_verify_omitted",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-http-out"));
};
`,
		},
		{
			name: "test_tls_version",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
						TLS: &output.TLS{
							SslVersion: "tlsv1_3",
						},
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") tls(ssl_version("tlsv1_3")) batch-lines(2000) workers(3) persist_name("output_default_test-http-out"));
};
`,
		},
		{
			name: "test_fifo_size",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers:     3,
						LogFIFOSize: 1000,
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-http-out") log-fifo-size(1000));
};
`,
		},
		{
			name: "test_timeout",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
						Timeout: 100,
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-http-out") timeout(100));
};
		`,
		},
		{
			name: "test_response_action",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-http-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					HTTP: &output.HTTPOutput{
						URL:     "test.local",
						Headers: []string{"a:b", "c:d"},
						Batch: output.Batch{
							BatchLines: 2000,
						},
						Workers: 3,
						ResponseAction: filter.RawArrowMap{
							"418": "drop",
						},
					},
				},
			},
			config: `destination "output_default_test-http-out" {
	http(url("test.local") headers("a:b" "c:d") batch-lines(2000) workers(3) persist_name("output_default_test-http-out") response-action(
		418 => drop
	));
};
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.CheckConfigForOutput(t, tt.output, tt.config)
		})
	}
}
