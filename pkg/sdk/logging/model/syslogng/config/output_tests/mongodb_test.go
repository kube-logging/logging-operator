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
	"strings"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
)

func TestMongoDBOutput(t *testing.T) {
	expectedConfig := config.Untab(`@version: current

@include "scl.conf"

source "main_input" {
	channel {
		source {
			network(flags("no-parse") port(601) transport("tcp"));
		};
		parser {
			json-parser(prefix("json."));
		};
	};
};

destination "output_default_test-mongodb-out" {
	mongodb(collection("messages") uri("mongodb://127.0.0.1:27017/syslog") value_pairs(scope("selected-macros" "nv-pairs" "sdata")) persist_name("output_default_test-mongodb-out"));
};
`)

	testCaseInput := config.Input{
		SyslogNGSpec:   &v1beta1.SyslogNGSpec{},
		Namespace:      "config-test",
		Name:           "test",
		ClusterOutputs: []v1beta1.SyslogNGClusterOutput{},
		ClusterFlows:   []v1beta1.SyslogNGClusterFlow{},
		Flows:          []v1beta1.SyslogNGFlow{},
		SourcePort:     601,
		SecretLoaderFactory: &config.TestSecretLoaderFactory{
			Reader: config.SecretReader{
				Secrets: []corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "mongodb-connectionstring",
						},
						Data: map[string][]byte{
							"connectionstring": []byte("mongodb://127.0.0.1:27017/syslog"),
						},
					},
				},
			},
			MountPath: "/etc/syslog-ng/secret",
		},
		Outputs: []v1beta1.SyslogNGOutput{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-mongodb-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					MongoDB: &output.MongoDB{
						Uri: &secret.Secret{
							ValueFrom: &secret.ValueFrom{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mongodb-connectionstring",
									},
									Key: "connectionstring",
								},
							},
						},
						Collection: "messages",
						ValuePairs: output.ValuePairs{
							Scope: `"selected-macros" "nv-pairs" "sdata"`,
						},
					},
				},
			},
		},
	}

	var buf strings.Builder
	err := config.RenderConfigInto(testCaseInput, &buf)
	config.CheckError(t, false, err)
	require.Equal(t, expectedConfig, buf.String())

}

func TestMongoDBOutputWithWriteConcernKeyword(t *testing.T) {
	expectedConfig := config.Untab(`@version: current

@include "scl.conf"

source "main_input" {
	channel {
		source {
			network(flags("no-parse") port(601) transport("tcp"));
		};
		parser {
			json-parser(prefix("json."));
		};
	};
};

destination "output_default_test-mongodb-out" {
	mongodb(collection("messages") uri("mongodb://127.0.0.1:27017/syslog") value_pairs(scope("selected-macros" "nv-pairs" "sdata")) persist_name("output_default_test-mongodb-out") write_concern(unacked));
};
`)

	testCaseInput := config.Input{
		SyslogNGSpec:   &v1beta1.SyslogNGSpec{},
		Namespace:      "config-test",
		Name:           "test",
		ClusterOutputs: []v1beta1.SyslogNGClusterOutput{},
		ClusterFlows:   []v1beta1.SyslogNGClusterFlow{},
		Flows:          []v1beta1.SyslogNGFlow{},
		SourcePort:     601,
		SecretLoaderFactory: &config.TestSecretLoaderFactory{
			Reader: config.SecretReader{
				Secrets: []corev1.Secret{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "mongodb-connectionstring",
						},
						Data: map[string][]byte{
							"connectionstring": []byte("mongodb://127.0.0.1:27017/syslog"),
						},
					},
				},
			},
			MountPath: "/etc/syslog-ng/secret",
		},
		Outputs: []v1beta1.SyslogNGOutput{
			{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "test-mongodb-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					MongoDB: &output.MongoDB{
						Uri: &secret.Secret{
							ValueFrom: &secret.ValueFrom{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "mongodb-connectionstring",
									},
									Key: "connectionstring",
								},
							},
						},
						Collection: "messages",
						ValuePairs: output.ValuePairs{
							Scope: `"selected-macros" "nv-pairs" "sdata"`,
						},
						WriteConcern: output.Unacked,
					},
				},
			},
		},
	}

	var buf strings.Builder
	err := config.RenderConfigInto(testCaseInput, &buf)
	config.CheckError(t, false, err)
	require.Equal(t, expectedConfig, buf.String())

}
