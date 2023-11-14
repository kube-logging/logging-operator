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

package fluentd

import (
	"bytes"
	"fmt"
	"html/template"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type fluentdConfig struct {
	LogLevel string
	Monitor  struct {
		Enabled bool
		Port    int32
		Path    string
	}
	IgnoreSameLogInterval     string
	IgnoreRepeatedLogInterval string
	Workers                   int32
	RootDir                   string
	EnableMsgpackTimeSupport  bool
}

func generateConfig(input fluentdConfig) (string, error) {
	output := new(bytes.Buffer)
	tmpl, err := template.New("test").Parse(fluentdInputTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse template")
	}
	err = tmpl.Execute(output, input)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute template")
	}
	return output.String(), nil
}

func (r *Reconciler) generateConfigSecret(fluentdSpec v1beta1.FluentdSpec) (map[string][]byte, error) {
	input := fluentdConfig{
		IgnoreSameLogInterval:     fluentdSpec.IgnoreSameLogInterval,
		IgnoreRepeatedLogInterval: fluentdSpec.IgnoreRepeatedLogInterval,
		EnableMsgpackTimeSupport:  fluentdSpec.EnableMsgpackTimeSupport,
		Workers:                   fluentdSpec.Workers,
		LogLevel:                  fluentdSpec.LogLevel,
	}

	input.RootDir = fluentdSpec.RootDir
	if input.RootDir == "" {
		input.RootDir = bufferPath
	}

	if fluentdSpec.Metrics != nil {
		input.Monitor.Enabled = true
		input.Monitor.Port = fluentdSpec.Metrics.Port
		input.Monitor.Path = fluentdSpec.Metrics.Path
	}

	inputConfig, err := generateConfig(input)
	if err != nil {
		return nil, err
	}

	configs := map[string][]byte{
		"fluent.conf":  []byte(fluentdDefaultTemplate),
		"input.conf":   []byte(inputConfig),
		"devnull.conf": []byte(fluentdOutputTemplate),
	}
	return configs, nil
}

func (r *Reconciler) secretConfig() (runtime.Object, reconciler.DesiredState, error) {
	configMap, err := r.generateConfigSecret(*r.fluentdSpec)
	if err != nil {
		return nil, nil, err
	}
	configMap["fluentlog.conf"] = []byte(fmt.Sprintf(fluentLog, r.fluentdSpec.FluentLogDestination))
	configs := &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(SecretConfigName, ComponentFluentd),
		Data:       configMap,
	}

	return configs, reconciler.StatePresent, nil
}
