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

package config

import (
	"strings"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/secret"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"
)

func TestRenderConfigInto(t *testing.T) {
	testCases := map[string]struct {
		input   Input
		wantOut string
		wantErr any
	}{
		"empty input": {
			input: Input{
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantErr: true,
		},
		"no syslog-ng spec": {
			input: Input{
				SyslogNGSpec:        nil,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantErr: true,
		},
		"single flow with single output": {
			input: Input{
				SourcePort:   601,
				Name:         "test",
				Namespace:    "config-test",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Outputs: []v1beta1.SyslogNGOutput{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-syslog-out",
						},
						Spec: v1beta1.SyslogNGOutputSpec{
							Syslog: &output.SyslogOutput{
								Host:      "test.local",
								Transport: "tcp",
							},
						},
					},
				},
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Match: &v1beta1.SyslogNGMatch{
								Regexp: &filter.RegexpMatchExpr{
									Pattern: "nginx",
									Value:   "kubernetes.labels.app",
								},
							},
							Filters: []v1beta1.SyslogNGFilter{
								{
									Rewrite: []filter.RewriteConfig{
										{
											Set: &filter.SetConfig{
												FieldName: "cluster",
												Value:     "test-cluster",
											},
										},
									},
								},
							},
							LocalOutputRefs: []string{"test-syslog-out"},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantOut: Untab(`@version: current

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

destination "output_default_test-syslog-out" {
	syslog("test.local" transport("tcp") persist_name("output_default_test-syslog-out"));
};

filter "flow_default_test-flow_ns_filter" {
	match("default" value("json.kubernetes.namespace_name") type("string"));
};
filter "flow_default_test-flow_match" {
	match("nginx" value("kubernetes.labels.app"));
};
rewrite "flow_default_test-flow_filters_0" {
	set("test-cluster" value("cluster"));
};
log {
	source("main_input");
	filter("flow_default_test-flow_ns_filter");
	filter("flow_default_test-flow_match");
	rewrite("flow_default_test-flow_filters_0");
	log {
        destination("output_default_test-syslog-out");
    };
};
`),
		},
		"single flow with output metrics": {
			input: Input{
				SourcePort:   601,
				Name:         "test",
				Namespace:    "config-test",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Outputs: []v1beta1.SyslogNGOutput{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-syslog-out",
						},
						Spec: v1beta1.SyslogNGOutputSpec{
							Syslog: &output.SyslogOutput{
								Host:      "test.local",
								Transport: "tcp",
							},
						},
					},
				},
				ClusterOutputs: []v1beta1.SyslogNGClusterOutput{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "config-test",
							Name:      "test-syslog-out-global",
						},
						Spec: v1beta1.SyslogNGClusterOutputSpec{
							SyslogNGOutputSpec: v1beta1.SyslogNGOutputSpec{
								Syslog: &output.SyslogOutput{
									Host:      "test.local",
									Transport: "tcp",
								},
							},
						},
					},
				},
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Match: &v1beta1.SyslogNGMatch{
								Regexp: &filter.RegexpMatchExpr{
									Pattern: "nginx",
									Value:   "kubernetes.labels.app",
								},
							},
							Filters: []v1beta1.SyslogNGFilter{
								{
									Rewrite: []filter.RewriteConfig{
										{
											Set: &filter.SetConfig{
												FieldName: "cluster",
												Value:     "test-cluster",
											},
										},
									},
								},
							},
							LocalOutputRefs:  []string{"test-syslog-out"},
							GlobalOutputRefs: []string{"test-syslog-out-global"},
							OutputMetrics: []filter.MetricsProbe{
								{
									Key: "example",
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantOut: Untab(`@version: current

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

destination "clusteroutput_config-test_test-syslog-out-global" {
	syslog("test.local" transport("tcp") persist_name("clusteroutput_config-test_test-syslog-out-global"));
};

destination "output_default_test-syslog-out" {
	syslog("test.local" transport("tcp") persist_name("output_default_test-syslog-out"));
};

filter "flow_default_test-flow_ns_filter" {
	match("default" value("json.kubernetes.namespace_name") type("string"));
};
filter "flow_default_test-flow_match" {
	match("nginx" value("kubernetes.labels.app"));
};
rewrite "flow_default_test-flow_filters_0" {
	set("test-cluster" value("cluster"));
};
log {
	source("main_input");
	filter("flow_default_test-flow_ns_filter");
	filter("flow_default_test-flow_match");
	rewrite("flow_default_test-flow_filters_0");
    log {
        parser {
            metrics-probe(key("example") labels(
                "logging" => "test"
				"output_name" => "test-syslog-out-global"
				"output_namespace" => "config-test"
				"output_scope" => "global"
            ));
        };
        destination("clusteroutput_config-test_test-syslog-out-global");
    };
    log {
        parser {
            metrics-probe(key("example") labels(
				"logging" => "test"
                "output_name" => "test-syslog-out"
				"output_namespace" => "default"
				"output_scope" => "local"
            ));
        };
        destination("output_default_test-syslog-out");
    };
};
`),
		},
		"global options": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					GlobalOptions: &v1beta1.GlobalOptions{
						StatsLevel: amp(3),
						StatsFreq:  amp(0),
					},
				},
				SourcePort:          601,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantOut: `@version: current

@include "scl.conf"

options {
    stats_level(3);
    stats_freq(0);
};

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
`,
		},
		"global options default": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					Metrics: &v1beta1.Metrics{
						Path: "/metrics",
					},
					GlobalOptions: &v1beta1.GlobalOptions{},
				},
				SourcePort:          601,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantOut: `@version: current

@include "scl.conf"

options {
    stats(level(2) freq(0));
};

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
`,
		},
		"global options_new_stats": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					GlobalOptions: &v1beta1.GlobalOptions{
						Stats: &v1beta1.Stats{
							Level: amp(3),
							Freq:  amp(0),
						},
					},
				},
				SourcePort:          601,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
			},
			wantOut: `@version: current

@include "scl.conf"

options {
    stats(level(3) freq(0));
};

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
`,
		},
		"rewrite condition": {
			input: Input{
				SyslogNGSpec:        &v1beta1.SyslogNGSpec{},
				SourcePort:          601,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Filters: []v1beta1.SyslogNGFilter{
								{
									Rewrite: []filter.RewriteConfig{
										{
											Unset: &filter.UnsetConfig{
												FieldName: "MESSAGE",
												Condition: &filter.MatchExpr{
													Not: &filter.MatchExpr{
														Regexp: &filter.RegexpMatchExpr{
															Pattern: "foo",
															Value:   "MESSAGE",
															Type:    "string",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantOut: Untab(`@version: current

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

filter "flow_default_test-flow_ns_filter" {
	match("default" value("json.kubernetes.namespace_name") type("string"));
};
rewrite "flow_default_test-flow_filters_0" {
    unset(value("MESSAGE") condition((not match("foo" value("MESSAGE") type("string")))));
};
log {
    source("main_input");
	filter("flow_default_test-flow_ns_filter");
    rewrite("flow_default_test-flow_filters_0");
};
`),
		},
		"output with secret": {
			input: Input{
				Name:         "test",
				Namespace:    "logging",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Outputs: []v1beta1.SyslogNGOutput{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "my-output",
						},
						Spec: v1beta1.SyslogNGOutputSpec{
							Syslog: &output.SyslogOutput{
								Host: "127.0.0.1",
								TLS: &output.TLS{
									CaFile: &secret.Secret{
										MountFrom: &secret.ValueFrom{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "my-secret",
												},
												Key: "tls.crt",
											},
										},
									},
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{
					Reader: SecretReader{
						Secrets: []corev1.Secret{
							{
								ObjectMeta: metav1.ObjectMeta{
									Namespace: "default",
									Name:      "my-secret",
								},
								Data: map[string][]byte{
									"tls.crt": []byte("asdf"),
								},
							},
						},
					},
					MountPath: "/etc/syslog-ng/secret",
				},
				SourcePort: 601,
			},
			wantOut: `@version: current

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

destination "output_default_my-output" {
    syslog("127.0.0.1" tls(ca_file("/etc/syslog-ng/secret/default-my-secret-tls.crt")) persist_name("output_default_my-output"));
};
`,
		},
		"clusteroutput with flow ref": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Name:         "test",
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							GlobalOutputRefs: []string{
								"clusterout",
							},
						},
					},
				},
				ClusterOutputs: []v1beta1.SyslogNGClusterOutput{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "clusterout",
							Namespace: "logging",
						},
						Spec: v1beta1.SyslogNGClusterOutputSpec{
							SyslogNGOutputSpec: v1beta1.SyslogNGOutputSpec{
								Syslog: &output.SyslogOutput{
									Host: "127.0.0.1",
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: `@version: current

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

destination "clusteroutput_logging_clusterout" {
    syslog("127.0.0.1" persist_name("clusteroutput_logging_clusterout"));
};

filter "flow_default_test-flow_ns_filter" {
    match("default" value("json.kubernetes.namespace_name") type("string"));
};
log {
    source("main_input");
    filter("flow_default_test-flow_ns_filter");
    log {
        destination("clusteroutput_logging_clusterout");
    };
};
`,
		},
		"flow referencing non-existent cluster output": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Name:         "test",
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							GlobalOutputRefs: []string{
								"clusterout",
							},
						},
					},
				},
				ClusterOutputs:      nil,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantErr: true,
		},
		"clusterFlow referencing non-existent cluster output": {
			input: Input{
				Name:         "test",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				ClusterFlows: []v1beta1.SyslogNGClusterFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGClusterFlowSpec{
							GlobalOutputRefs: []string{
								"clusterout",
							},
						},
					},
				},
				ClusterOutputs:      nil,
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantErr: true,
		},
		"parser": {
			input: Input{
				Name:         "test",
				Namespace:    "logging",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Filters: []v1beta1.SyslogNGFilter{
								{
									Parser: &filter.ParserConfig{
										Regexp: &filter.RegexpParser{
											Patterns: []string{
												".*test_field -> (?<test_field>.*)$",
											},
											Prefix: ".regexp.",
										},
									},
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: `@version: current

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

filter "flow_default_test-flow_ns_filter" {
    match("default" value("json.kubernetes.namespace_name") type("string"));
};
parser "flow_default_test-flow_filters_0" {
    regexp-parser(patterns(".*test_field -> (?<test_field>.*)$") prefix(".regexp."));
};
log {
    source("main_input");
    filter("flow_default_test-flow_ns_filter");
    parser("flow_default_test-flow_filters_0");
};
`,
		},
		"filter with name": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Namespace:    "logging",
				Name:         "test",
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Filters: []v1beta1.SyslogNGFilter{
								{
									ID: "remove message",
									Rewrite: []filter.RewriteConfig{
										{
											Unset: &filter.UnsetConfig{
												FieldName: "MESSAGE",
											},
										},
									},
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

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

filter "flow_default_test-flow_ns_filter" {
	match("default" value("json.kubernetes.namespace_name") type("string"));
};
rewrite "flow_default_test-flow_filters_remove message" {
    unset(value("MESSAGE"));
};
log {
    source("main_input");
	filter("flow_default_test-flow_ns_filter");
    rewrite("flow_default_test-flow_filters_remove message");
};
`),
		},
		"groupunset": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{},
				Namespace:    "logging",
				Name:         "test",
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Filters: []v1beta1.SyslogNGFilter{
								{
									ID: "remove message",
									Rewrite: []filter.RewriteConfig{
										{
											GroupUnset: &filter.GroupUnsetConfig{
												Pattern: ".SDATA.*",
											},
										},
									},
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

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

filter "flow_default_test-flow_ns_filter" {
	match("default" value("json.kubernetes.namespace_name") type("string"));
};
rewrite "flow_default_test-flow_filters_remove message" {
    groupunset(values(".SDATA.*"));
};
log {
    source("main_input");
	filter("flow_default_test-flow_ns_filter");
    rewrite("flow_default_test-flow_filters_remove message");
};
`),
		},
		"custom json key delimiter": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					JSONKeyDelimiter: ";",
				},
				Namespace: "logging",
				Name:      "test",
				Flows: []v1beta1.SyslogNGFlow{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "default",
							Name:      "test-flow",
						},
						Spec: v1beta1.SyslogNGFlowSpec{
							Filters: []v1beta1.SyslogNGFilter{
								{
									Match: &filter.MatchConfig{
										Regexp: &filter.RegexpMatchExpr{
											Pattern: "asdf",
											Value:   "ghjk",
										},
									},
								},
							},
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

@include "scl.conf"

source "main_input" {
    channel {
        source {
            network(flags("no-parse") port(601) transport("tcp"));
        };
        parser {
            json-parser(prefix("json;") key-delimiter(";"));
        };
    };
};

filter "flow_default_test-flow_ns_filter" {
    match("default" value("json;kubernetes;namespace_name") type("string"));
};
filter "flow_default_test-flow_filters_0" {
    match("asdf" value("ghjk"));
};
log {
    source("main_input");
    filter("flow_default_test-flow_ns_filter");
    filter("flow_default_test-flow_filters_0");
};
`),
		},
		"custom json key prefix": {
			input: Input{
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					JSONKeyPrefix: "asdf.",
				},
				Namespace:           "logging",
				Name:                "test",
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

@include "scl.conf"

source "main_input" {
    channel {
        source {
            network(flags("no-parse") port(601) transport("tcp"));
        };
        parser {
            json-parser(prefix("asdf."));
        };
    };
};
`),
		},
		"source metrics": {
			input: Input{
				Namespace: "logging",
				Name:      "test",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					SourceMetrics: []filter.MetricsProbe{
						{
							Key: "example",
							Labels: filter.ArrowMap{
								"a": "b",
							},
							Level: 2,
						},
						{
							Key: "example2",
							Labels: filter.ArrowMap{
								"c": "d",
							},
							Level: 3,
						},
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

@include "scl.conf"

source "main_input" {
    channel {
        source {
            network(flags("no-parse") port(601) transport("tcp"));
        };
        parser {
            json-parser(prefix("json."));
            metrics-probe(key("example") labels(
				"a" => "b"
				"logging" => "test"
			) level(2));
            metrics-probe(key("example2") labels(
				"c" => "d"
				"logging" => "test"
			) level(3));
        };
    };
};
`),
		},
		"date-parser default": {
			input: Input{
				Namespace: "logging",
				Name:      "test",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					SourceDateParser: &v1beta1.SourceDateParser{},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

@include "scl.conf"

source "main_input" {
    channel {
        source {
            network(flags("no-parse") port(601) transport("tcp"));
        };
        parser {
            json-parser(prefix("json."));
            date-parser(format("%FT%T.%f%z") template("${json.time}"));
        };
    };
};
`),
		},
		"date-parser custom": {
			input: Input{
				Namespace: "logging",
				Name:      "test",
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					SourceDateParser: &v1beta1.SourceDateParser{
						Format:   utils.StringPointer("asd"),
						Template: utils.StringPointer("bsd"),
					},
				},
				SecretLoaderFactory: &TestSecretLoaderFactory{},
				SourcePort:          601,
			},
			wantOut: Untab(`@version: current

@include "scl.conf"

source "main_input" {
    channel {
        source {
            network(flags("no-parse") port(601) transport("tcp"));
        };
        parser {
            json-parser(prefix("json."));
            date-parser(format("asd") template("bsd"));
        };
    };
};
`),
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			var buf strings.Builder
			err := RenderConfigInto(testCase.input, &buf)
			CheckError(t, testCase.wantErr, err)
			require.Equal(t, testCase.wantOut, buf.String())
		})
	}
}
