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

func TestLokiOutputTable(t *testing.T) {
	var tests = []struct {
		name   string
		output v1beta1.SyslogNGOutput
		config string
	}{
		{
			name: "test_general",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-loki-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Loki: &output.LokiOutput{
						URL:          "test.local",
						BatchLines:   2000,
						BatchTimeout: 10,
						Workers:      3,
						LogFIFOSize:  1000,
					},
				},
			},
			config: `destination "output_default_test-loki-out" {
	loki(url("test.local") batch-lines(2000) batch-timeout(10) workers(3) persist_name("output_default_test-loki-out") log-fifo-size(1000));
};
`,
		},
		{
			name: "test_labels",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-loki-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Loki: &output.LokiOutput{
						URL:          "test.local",
						BatchLines:   2000,
						BatchTimeout: 10,
						Workers:      3,
						LogFIFOSize:  1000,
						Labels: filter.ArrowMap{
							"app":  "$PROGRAM",
							"host": "$HOST",
						},
					},
				},
			},
			config: `destination "output_default_test-loki-out" {
	loki(labels(
		"app" => "$PROGRAM"
		"host" => "$HOST"
	) url("test.local") batch-lines(2000) batch-timeout(10) workers(3) persist_name("output_default_test-loki-out") log-fifo-size(1000));
};
`,
		},
		{
			name: "test_timestamp",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-loki-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Loki: &output.LokiOutput{
						URL:          "test.local",
						BatchLines:   2000,
						BatchTimeout: 10,
						Workers:      3,
						LogFIFOSize:  1000,
						Timestamp:    "msg",
					},
				},
			},
			config: `destination "output_default_test-loki-out" {
	loki(url("test.local") batch-lines(2000) batch-timeout(10) workers(3) persist_name("output_default_test-loki-out") log-fifo-size(1000) timestamp("msg"));
};
`,
		},
		{
			name: "test_template",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-loki-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Loki: &output.LokiOutput{
						URL:          "test.local",
						BatchLines:   2000,
						BatchTimeout: 10,
						Workers:      3,
						LogFIFOSize:  1000,
						Template:     "$ISODATE $HOST $MSGHDR$MSG",
					},
				},
			},
			config: `destination "output_default_test-loki-out" {
	loki(url("test.local") batch-lines(2000) batch-timeout(10) workers(3) persist_name("output_default_test-loki-out") log-fifo-size(1000) template("$ISODATE $HOST $MSGHDR$MSG"));
};
`,
		},
		{
			name: "test_auth",
			output: v1beta1.SyslogNGOutput{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-loki-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Loki: &output.LokiOutput{
						URL:          "test.local",
						BatchLines:   2000,
						BatchTimeout: 10,
						Workers:      3,
						LogFIFOSize:  1000,
						Auth: &output.Auth{
							Insecure: &output.Insecure{},
						},
					},
				},
			},
			config: `destination "output_default_test-loki-out" {
	loki(auth(insecure()) url("test.local") batch-lines(2000) batch-timeout(10) workers(3) persist_name("output_default_test-loki-out") log-fifo-size(1000));
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
