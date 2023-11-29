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

func TestS3OutputMinimal(t *testing.T) {
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

destination "output_default_test-s3-out" {
	s3(url("http://localhost:9000") bucket("s3bucket") access_key("access-key-secret-value") secret_key("secret-key-value") object_key("${HOST}/my-logs") persist_name("output_default_test-s3-out"));
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
							Name:      "s3-secrets",
						},
						Data: map[string][]byte{
							"secret_key": []byte("secret-key-value"),
							"object_key": []byte("${HOST}/my-logs"),
							"access_key": []byte("access-key-secret-value"),
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
					Name:      "test-s3-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{
					S3: &output.S3Output{
						Url:    "http://localhost:9000",
						Bucket: "s3bucket",
						AccessKey: &secret.Secret{
							ValueFrom: &secret.ValueFrom{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "s3-secrets",
									},
									Key: "access_key",
								},
							},
						},
						SecretKey: &secret.Secret{
							ValueFrom: &secret.ValueFrom{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "s3-secrets",
									},
									Key: "secret_key",
								},
							},
						},
						ObjectKey: "${HOST}/my-logs",
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

func TestS3OutputAllOptions(t *testing.T) {
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

destination "output_default_test-s3-out" {
	s3(url("http://localhost:9000") bucket("s3bucket") access_key("access-key-secret-value") secret_key("secret-key-value") object_key("${HOST}/my-logs") object_key_timestamp("timestamp") template("${MESSAGE}\n") compression(no) compresslevel(9) chunk_size(5) max_object_size(5000) upload_threads(8) max_pending_uploads(32) flush_grace_period(60) region("s3region") storage_class("STANDARD") canned_acl("s3-canned-acl") persist_name("output_default_test-s3-out"));
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
							Name:      "s3-secrets",
						},
						Data: map[string][]byte{
							"secret_key": []byte("secret-key-value"),
							"object_key": []byte("${HOST}/my-logs"),
							"access_key": []byte("access-key-secret-value"),
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
					Name:      "test-s3-out",
				},
				Spec: v1beta1.SyslogNGOutputSpec{S3: &output.S3Output{
					Url:    "http://localhost:9000",
					Bucket: "s3bucket",
					AccessKey: &secret.Secret{
						ValueFrom: &secret.ValueFrom{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "s3-secrets",
								},
								Key: "access_key",
							},
						},
					},
					SecretKey: &secret.Secret{
						ValueFrom: &secret.ValueFrom{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "s3-secrets",
								},
								Key: "secret_key",
							},
						},
					},
					ObjectKey:          "${HOST}/my-logs",
					ObjectKeyTimestamp: `"timestamp"`,
					Template:           `"${MESSAGE}\n"`,
					Compression:        config.NewFalse(),
					CompressLevel:      9,
					ChunkSize:          5,
					MaxObjectSize:      5000,
					UploadThreads:      8,
					MaxPendingUploads:  32,
					FlushGracePeriod:   60,
					Region:             "s3region",
					StorageClass:       "STANDARD",
					CannedAcl:          "s3-canned-acl",
				}},
			},
		},
	}
	var buf strings.Builder
	err := config.RenderConfigInto(testCaseInput, &buf)
	config.CheckError(t, false, err)
	require.Equal(t, expectedConfig, buf.String())

}
