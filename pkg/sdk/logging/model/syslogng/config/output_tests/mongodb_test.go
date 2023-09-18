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
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMongoDBOutput(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-mongodb-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				MongoDB: &output.MongoDB{
					Uri:        "mongodb://127.0.0.1:27017/syslog",
					Collection: "messages",
					Compaction: *config.NewFalse(),
					ValuePairs: output.ValuePairs{
						Scope: output.RawString{
							String: `"selected-macros" "nv-pairs" "sdata"`,
						},
					},
				},
			},
		},
		`
destination "output_default_test-mongodb-out" {
	mongodb(collection("messages") compaction(no) uri("mongodb://127.0.0.1:27017/syslog") value_pairs(scope("selected-macros" "nv-pairs" "sdata")) persist_name("output_default_test-mongodb-out"));
};
`,
	)
}

func TestMongoDBOutputWithWriteConcernKeyword(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-mongodb-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				MongoDB: &output.MongoDB{
					Uri:        "mongodb://127.0.0.1:27017/syslog",
					Collection: "messages",
					Compaction: *config.NewFalse(),
					ValuePairs: output.ValuePairs{
						Scope: output.RawString{
							String: `"selected-macros" "nv-pairs" "sdata"`,
						},
					},
					WriteConcern: output.RawString{String: output.Unacked},
				},
			},
		},
		`
destination "output_default_test-mongodb-out" {
	mongodb(collection("messages") compaction(no) uri("mongodb://127.0.0.1:27017/syslog") value_pairs(scope("selected-macros" "nv-pairs" "sdata")) persist_name("output_default_test-mongodb-out") write_concern(unacked));
};
`,
	)
}

func TestMongoDBOutputWithWriteConcernNonNegativeInteger(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-mongodb-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				MongoDB: &output.MongoDB{
					Uri:        "mongodb://127.0.0.1:27017/syslog",
					Collection: "messages",
					Compaction: *config.NewFalse(),
					ValuePairs: output.ValuePairs{
						Scope: output.RawString{
							String: `"selected-macros" "nv-pairs" "sdata"`,
						},
					},
					WriteConcern: output.RawString{String: `5`},
				},
			},
		},
		`
destination "output_default_test-mongodb-out" {
	mongodb(collection("messages") compaction(no) uri("mongodb://127.0.0.1:27017/syslog") value_pairs(scope("selected-macros" "nv-pairs" "sdata")) persist_name("output_default_test-mongodb-out") write_concern(5));
};
`,
	)
}
