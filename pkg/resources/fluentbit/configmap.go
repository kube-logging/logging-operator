/*
 * Copyright © 2019 Banzai Cloud
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fluentbit

import (
	"bytes"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"text/template"
)

type fluentbitTLSConfig struct {
	Enabled    bool
	SharedKey  string
	CACertFile string
	CertFile   string
	KeyFile    string
}

type fluentBitConfig struct {
	Namespace string
	TLS       fluentbitTLSConfig
	Monitor   map[string]string
	Output    map[string]string
}

func (r *Reconciler) configMap() runtime.Object {
	var monitorConfig map[string]string
	if _, ok := r.Fluentbit.Spec.Annotations["prometheus.io/port"]; ok {
		monitorConfig = map[string]string{
			"Port": r.Fluentbit.Spec.Annotations["prometheus.io/port"],
		}
	}
	tlsConfig := fluentbitTLSConfig{
		Enabled:   r.Fluentbit.Spec.TLS.Enabled,
		SharedKey: r.Fluentbit.Spec.TLS.SharedKey,
	}
	if r.Fluentbit.Spec.TLS.SecretType == "tls" {
		tlsConfig.CertFile = "/fluent-bit/tls/tls.crt"
		tlsConfig.KeyFile = "/fluent-bit/tls/tls.key"
		tlsConfig.CACertFile = "/fluent-bit/tls/ca.crt"
	} else {
		tlsConfig.CertFile = "/fluent-bit/tls/clientCert"
		tlsConfig.KeyFile = "/fluent-bit/tls/clientKey"
		tlsConfig.CACertFile = "/fluent-bit/tls/caCert"
	}
	input := fluentBitConfig{
		Namespace: r.Fluentbit.Namespace,
		TLS:       tlsConfig,
		Monitor:   monitorConfig,
	}
	return &corev1.ConfigMap{
		ObjectMeta: templates.FluentbitObjectMeta(fluentbitConfigMapName, r.Fluentbit.Labels, r.Fluentbit),
		Data: map[string]string{
			"fluent-bit.conf": generateConfig(input),
		},
	}
}

func generateConfig(input fluentBitConfig) string {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentBitConfigTemplate)
	if err != nil {
		return ""
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return ""
	}
	outputString := output.String()
	return outputString
}
