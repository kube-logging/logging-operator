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
	"context"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"text/template"

	"github.com/Masterminds/sprig"
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var log = logf.Log.WithName("plugins.configmap")

func generateFluentdConfig(plugin *loggingv1alpha1.Plugin, client client.Client) (string, string) {
	var finalConfig string
	// Generate filters
	for _, filter := range plugin.Spec.Filter {
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

// RenderPlugin general FPlugin renderer
func renderPlugin(plugin loggingv1alpha1.FPlugin, baseMap map[string]string, namespace string, client client.Client) (string, error) {
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
	appConfigData := map[string]string{}
	labels := map[string]string{}
	for _, plugin := range r.PluginList.Items {
		labels = util.MergeLabels(labels, plugin.Labels)
		name, data := generateFluentdConfig(&plugin, r.Client)
		if name != "" {
			name = name + ".conf"
		}
		appConfigData[name] = data
	}
	pluginConfigMapNamespace := r.Namespace
	cmLog := log.WithValues("pluginConfigMapNamespace", pluginConfigMapNamespace)
	fluentdList := loggingv1alpha1.FluentdList{}
	err := r.Client.List(context.TODO(), client.MatchingLabels(map[string]string{}), &fluentdList)
	if err != nil {
		cmLog.Error(err, "Reconciler query failed.")
	}

	if len(fluentdList.Items) > 0 {
		cmLog = log.WithValues("pluginConfigMapNamespace", pluginConfigMapNamespace, "FluentdNamespace", fluentdList.Items[0].Namespace)
		cmLog.Info("Check Fluentd Namespace")
		if pluginConfigMapNamespace != fluentdList.Items[0].Namespace {
			pluginConfigMapNamespace = fluentdList.Items[0].Namespace
			cmLog = log.WithValues("pluginConfigMapNamespace", pluginConfigMapNamespace)
			cmLog.Info("Plugin ConfigMap Namespace Updated")

		}
	} else {
		log.Info("The is no Fluentd resource available")
	}
	return &corev1.ConfigMap{
		ObjectMeta: templates.PluginsObjectMeta(appConfigMapName, util.MergeLabels(map[string]string{}, labelSelector), pluginConfigMapNamespace),
		Data:       appConfigData,
	}
}
