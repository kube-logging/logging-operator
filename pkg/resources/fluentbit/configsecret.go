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

package fluentbit

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources/fluentd"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/sdk/model/types"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type fluentBitConfig struct {
	Namespace string
	TLS       struct {
		Enabled   bool
		SharedKey string
	}
	Monitor struct {
		Enabled bool
		Port    int32
		Path    string
	}
	Output        map[string]string
	TargetHost    string
	TargetPort    int32
	Input         map[string]string
	Tailers       []map[string]string
	Filter        map[string]string
	BufferStorage map[string]string
}

func (r *Reconciler) configSecret() (runtime.Object, k8sutil.DesiredState) {
	if r.Logging.Spec.FluentbitSpec.CustomConfigSecret != "" {
		return &corev1.Secret{
			ObjectMeta: templates.FluentbitObjectMeta(
				r.Logging.QualifiedName(fluentBitSecretConfigName), r.Logging.Labels, r.Logging),
		}, k8sutil.StateAbsent
	}
	monitor := struct {
		Enabled bool
		Port    int32
		Path    string
	}{}
	if r.Logging.Spec.FluentbitSpec.Metrics != nil {
		monitor.Enabled = true
		monitor.Port = r.Logging.Spec.FluentbitSpec.Metrics.Port
		monitor.Path = r.Logging.Spec.FluentbitSpec.Metrics.Path
	}

	if r.Logging.Spec.FluentbitSpec.ContainerTail.Parser == "" {
		switch r.CRI {
		case "docker":
			r.Logging.Spec.FluentbitSpec.ContainerTail.Parser = "docker"
		case "containerd":
			r.Logging.Spec.FluentbitSpec.ContainerTail.Parser = "cri"
		default:
			r.Logging.Spec.FluentbitSpec.ContainerTail.Parser = "cri"
		}
	}

	mapper := types.NewStructToStringMapper(nil)
	fluentbitInput, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.ContainerTail)
	if err != nil {
		log.Error(err)
	}

	fluentbitTailers := make([]map[string]string, 0)
	if len(r.Logging.Spec.FluentbitSpec.Tailers) > 0 {
		for _, tailer := range r.Logging.Spec.FluentbitSpec.Tailers {
			mappedTailer, err := mapper.StringsMap(tailer)
			if err != nil {
				log.Error(err)
			} else {
				fluentbitTailers = append(fluentbitTailers, mappedTailer)
			}
		}
	}

	fluentbitFilter, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.FilterKubernetes)
	if err != nil {
		log.Error(err)
	}

	fluentbitBufferStorage, err := mapper.StringsMap(r.Logging.Spec.FluentbitSpec.BufferStorage)
	if err != nil {
		log.Error(err)
	}

	input := fluentBitConfig{
		Namespace: r.Logging.Spec.ControlNamespace,
		TLS: struct {
			Enabled   bool
			SharedKey string
		}{
			Enabled:   r.Logging.Spec.FluentbitSpec.TLS.Enabled,
			SharedKey: r.Logging.Spec.FluentbitSpec.TLS.SharedKey,
		},
		Monitor:       monitor,
		TargetHost:    fmt.Sprintf("%s.%s.svc", r.Logging.QualifiedName(fluentd.ServiceName), r.Logging.Spec.ControlNamespace),
		TargetPort:    r.Logging.Spec.FluentdSpec.Port,
		Input:         fluentbitInput,
		Tailers:       fluentbitTailers,
		Filter:        fluentbitFilter,
		BufferStorage: fluentbitBufferStorage,
	}
	if r.Logging.Spec.FluentbitSpec.TargetHost != "" {
		input.TargetHost = r.Logging.Spec.FluentbitSpec.TargetHost
	}
	if r.Logging.Spec.FluentbitSpec.TargetPort != 0 {
		input.TargetPort = r.Logging.Spec.FluentbitSpec.TargetPort
	}

	r.desiredConfig, err = generateConfig(input)
	if err != nil {
		log.Error(err)
	}

	return &corev1.Secret{
		ObjectMeta: templates.FluentbitObjectMeta(
			r.Logging.QualifiedName(fluentBitSecretConfigName), r.Logging.Labels, r.Logging),
		Data: map[string][]byte{
			"fluent-bit.conf": []byte(r.desiredConfig),
		},
	}, k8sutil.StatePresent
}

func generateConfig(input fluentBitConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentBitConfigTemplate)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", err
	}
	outputString := output.String()
	return outputString, nil
}
