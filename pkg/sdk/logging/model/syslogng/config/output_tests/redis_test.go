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

func TestRedisOutput(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-redis-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				Redis: &output.RedisOutput{
					Host:                "127.0.0.1",
					Port:                6379,
					CommandAndArguments: []string{"HINCRBY", "hosts", "$HOST", "1"},
				},
			},
		},
		`
destination "output_default_test-redis-out" {
	redis(host("127.0.0.1") port(6379) command("HINCRBY", "hosts", "$HOST", "1") persist_name("output_default_test-redis-out"));
};
`,
	)
}

func TestRedisOutputWithBatching(t *testing.T) {
	config.CheckConfigForOutput(t,
		v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "default",
				Name:      "test-redis-out",
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				Redis: &output.RedisOutput{
					Host:                "127.0.0.1",
					Port:                6379,
					CommandAndArguments: []string{"HINCRBY", "hosts", "$HOST", "1"},
					Batch:               output.Batch{BatchLines: 100, BatchTimeout: 10000},
					LogFIFOSize:         100000,
				},
			},
		},
		`
destination "output_default_test-redis-out" {
	redis(host("127.0.0.1") port(6379) command("HINCRBY", "hosts", "$HOST", "1") batch-lines(100) batch-timeout(10000) log-fifo-size(100000) persist_name("output_default_test-redis-out"));
};
`,
	)
}

func TestRedisOutputWithAuthentication(t *testing.T) {
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

destination "output_default_test-redis-out" {
	redis(host("127.0.0.1") auth("secret-redis-pwd") port(6379) command("HINCRBY", "hosts", "$HOST", "1") persist_name("output_default_test-redis-out"));
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
							Name:      "redis-password",
						},
						Data: map[string][]byte{
							"redis-pwd": []byte("secret-redis-pwd"),
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
					Name:      "test-redis-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					Redis: &output.RedisOutput{
						Host: "127.0.0.1",
						Auth: &secret.Secret{
							ValueFrom: &secret.ValueFrom{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "redis-password",
									},
									Key: "redis-pwd",
								},
							},
						},
						Port:                6379,
						CommandAndArguments: []string{"HINCRBY", "hosts", "$HOST", "1"},
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
