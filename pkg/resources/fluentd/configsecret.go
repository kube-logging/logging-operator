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

func (r *Reconciler) generateConfigSecret() (map[string][]byte, error) {
	input := fluentdConfig{
		IgnoreSameLogInterval:     r.Logging.Spec.FluentdSpec.IgnoreSameLogInterval,
		IgnoreRepeatedLogInterval: r.Logging.Spec.FluentdSpec.IgnoreRepeatedLogInterval,
		EnableMsgpackTimeSupport:  r.Logging.Spec.FluentdSpec.EnableMsgpackTimeSupport,
		Workers:                   r.Logging.Spec.FluentdSpec.Workers,
		LogLevel:                  r.Logging.Spec.FluentdSpec.LogLevel,
	}

	input.RootDir = r.Logging.Spec.FluentdSpec.RootDir
	if input.RootDir == "" {
		input.RootDir = bufferPath
	}

	if r.Logging.Spec.FluentdSpec.Metrics != nil {
		input.Monitor.Enabled = true
		input.Monitor.Port = r.Logging.Spec.FluentdSpec.Metrics.Port
		input.Monitor.Path = r.Logging.Spec.FluentdSpec.Metrics.Path
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
	configMap, err := r.generateConfigSecret()
	if err != nil {
		return nil, nil, err
	}
	configMap["fluentlog.conf"] = []byte(fmt.Sprintf(fluentLog, r.Logging.Spec.FluentdSpec.FluentLogDestination))
	configs := &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(SecretConfigName, ComponentFluentd),
		Data:       configMap,
	}

	return configs, reconciler.StatePresent, nil
}
