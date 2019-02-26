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

package plugins

import (
	"bytes"
	"github.com/Masterminds/sprig"
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"text/template"
)

func generateFluentdConfig(plugin *loggingv1alpha1.LoggingPlugin, client client.Client) (string, string) {
	var finalConfig string
	// Generate filters
	for _, filter := range plugin.Spec.Filter {
		logrus.Info("Applying filter")
		values, err := GetDefaultValues(filter.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		values["pattern"] = plugin.Spec.Input.Label["app"]
		config, err := renderPlugin(filter, values, plugin.Namespace, client)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		finalConfig += config
	}

	// Generate output
	for _, output := range plugin.Spec.Output {
		values, err := GetDefaultValues(output.Type)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		values["pattern"] = plugin.Spec.Input.Label["app"]
		config, err := renderPlugin(output, values, plugin.Namespace, client)
		if err != nil {
			logrus.Infof("Error in rendering template: %s", err)
			return "", ""
		}
		finalConfig += config
	}
	return plugin.Name, finalConfig

}

// RenderPlugin general Plugin renderer
func renderPlugin(plugin loggingv1alpha1.Plugin, baseMap map[string]string, namespace string, client client.Client) (string, error) {
	rawTemplate, err := GetTemplate(plugin.Type)
	if err != nil {
		return "", err
	}
	for _, param := range plugin.Parameters {
		k, v := param.GetValue(namespace, client)
		baseMap[k] = v
	}

	t := template.New("PluginTemplate").Funcs(sprig.TxtFuncMap())
	t, err = t.Parse(rawTemplate)
	if err != nil {
		return "", err
	}
	tpl := new(bytes.Buffer)
	err = t.Execute(tpl, baseMap)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func (r *Reconciler) appConfigMap() runtime.Object {
	name, data := generateFluentdConfig(r.Plugin, r.Client)
	if name != "" {
		name = name + ".conf"
	}

	return &corev1.ConfigMap{
		ObjectMeta: templates.PluginsObjectMeta(appConfigMapName, util.MergeLabels(r.Plugin.Labels, labelSelector), r.Plugin),
		Data: map[string]string{
			name: data,
		},
	}
}
