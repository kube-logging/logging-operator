// Copyright © 2023 Kube logging authors
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
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/config"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestOpenobserveOutputWithAuthentication(t *testing.T) {
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

destination "output_default_test-openobserve-out" {
	openobserve-log(url("http://localhost") user("root@example.com") password("V2tsn88GhdNTKxaS") persist_name("output_default_test-openobserve-out") port(5080) organization("default") stream("default"));
};
`)

	testCaseInput := config.Input{
		Namespace:      "config-test",
		Name:           "test",
		SyslogNGSpec:   &v1beta1.SyslogNGSpec{},
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
							Name:      "openobserve",
						},
						Data: map[string][]byte{
							"password": []byte("V2tsn88GhdNTKxaS"),
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
					Name:      "test-openobserve-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Openobserve: &output.OpenobserveOutput{
						HTTPOutput: output.HTTPOutput{
							URL:  "http://localhost",
							User: "root@example.com",
							Password: secret.Secret{
								ValueFrom: &secret.ValueFrom{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "openobserve",
										},
										Key: "password",
									},
								}},
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

func TestOpenobserveOutputWithOtherPort(t *testing.T) {
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

destination "output_default_test-openobserve-out" {
	openobserve-log(url("http://localhost") user("root@example.com") password("V2tsn88GhdNTKxaS") persist_name("output_default_test-openobserve-out") port(5081) organization("default") stream("default"));
};
`)

	testCaseInput := config.Input{
		Namespace:      "config-test",
		Name:           "test",
		SyslogNGSpec:   &v1beta1.SyslogNGSpec{},
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
							Name:      "openobserve",
						},
						Data: map[string][]byte{
							"password": []byte("V2tsn88GhdNTKxaS"),
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
					Name:      "test-openobserve-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Openobserve: &output.OpenobserveOutput{
						Port: 5081,
						HTTPOutput: output.HTTPOutput{
							URL:  "http://localhost",
							User: "root@example.com",
							Password: secret.Secret{
								ValueFrom: &secret.ValueFrom{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "openobserve",
										},
										Key: "password",
									},
								}},
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
