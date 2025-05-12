// Copyright © 2025 Kube logging authors
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

package axosyslog

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	"github.com/Masterminds/sprig/v3"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	axoSyslogConfigName = "axosyslog-config"
	axoSyslogConfigKey  = "axosyslog.conf"
)

func CreateAxoSyslogConfig(object any) (runtime.Object, reconciler.DesiredState, error) {
	axoSyslog, ok := object.(*v1beta1.AxoSyslog)
	if !ok {
		return nil, reconciler.StateAbsent, fmt.Errorf("expected *v1beta1.AxoSyslog, got %T", axoSyslog)
	}

	tmpl, err := template.New("AxoSyslogConfig").Funcs(sprig.FuncMap()).Parse(heredoc.Doc(`
		@version: current
		@include "scl.conf"

		options {
			stats(level(2) freq(0));
			keep-hostname(yes);
			keep-timestamp(yes);
			log-msg-size(64KiB);
			trim-large-messages(yes);
			time-reopen(1);
			dns-cache(yes);
			log-level("default");
			use-uniqid(yes);
			create-dirs(yes);
		};

		source "axosyslog-otlp" {
			channel {
				source { opentelemetry(port(4317) log-iw-size(1000000) log-fetch-limit(100000) workers(10) ` + "`__VARARGS__`" + `); };
				if {
					filterx {
						declare resource = otel_resource(${.otel_raw.resource});
						declare scope = otel_scope(${.otel_raw.scope});
						declare log = otel_logrecord(${.otel_raw.log});

						if ((log.observed_time_unix_nano == 0) ?? true) {
							log.observed_time_unix_nano = $R_UNIXTIME;
						};
					};
				};
			};
		};

		{{ range .Destinations }}
		destination "{{ .Name }}" { {{ .Config | nindent 2 }}
		};
		{{ end }}

		{{ range .LogPaths }}
		log "axosyslog-default" {
			source("axosyslog-otlp");
			{{ if .Filterx }}
			filterx { {{ .Filterx | nindent 4 }}
			};
			{{ end }}
			{{ if .Destination }}
			log {
				destination("{{ .Destination }}");
				flags(flow-control);
			};
			{{ end }}
		};
		{{ end }}`))
	if err != nil {
		return nil, reconciler.StateAbsent, err
	}

	var configBuffer bytes.Buffer
	if err := tmpl.Execute(&configBuffer, axoSyslog.Spec); err != nil {
		return nil, reconciler.StateAbsent, err
	}

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      axoSyslogConfigName,
			Namespace: axoSyslog.Namespace,
			Labels: map[string]string{
				LabelAppName:      commonAxoSyslogObjectValue,
				LabelAppComponent: commonAxoSyslogObjectValue,
			},
		},
		Data: map[string]string{
			axoSyslogConfigKey: configBuffer.String(),
		},
	}, reconciler.StatePresent, nil
}
