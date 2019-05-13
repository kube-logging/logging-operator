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

package templates

import (
	loggingv1alpha1 "github.com/banzaicloud/logging-operator/pkg/apis/logging/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PluginsObjectMeta creates an objectMeta for resource plugin
func PluginsObjectMeta(name string, labels map[string]string, namespace string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
	return o
}

// FluentdObjectMeta creates an objectMeta for resource fluentd
func FluentdObjectMeta(name string, labels map[string]string, fluentd *loggingv1alpha1.Fluentd) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      name,
		Namespace: fluentd.Namespace,
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: fluentd.APIVersion,
				Kind:       fluentd.Kind,
				Name:       fluentd.Name,
				UID:        fluentd.UID,
			},
		},
	}
	return o
}

// FluentdObjectMetaClusterScope creates an objectMeta for resource fluentd
func FluentdObjectMetaClusterScope(name string, labels map[string]string, fluentd *loggingv1alpha1.Fluentd) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   name,
		Labels: labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: fluentd.APIVersion,
				Kind:       fluentd.Kind,
				Name:       fluentd.Name,
				UID:        fluentd.UID,
			},
		},
	}
	return o
}

// FluentbitObjectMeta creates an objectMeta for resource fluentbit
func FluentbitObjectMeta(name string, labels map[string]string, fluentbit *loggingv1alpha1.Fluentbit) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      name,
		Namespace: fluentbit.Namespace,
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: fluentbit.APIVersion,
				Kind:       fluentbit.Kind,
				Name:       fluentbit.Name,
				UID:        fluentbit.UID,
			},
		},
	}
	return o
}

// FluentbitObjectMetaClusterScope creates an cluster scoped objectMeta for resource fluentbit
func FluentbitObjectMetaClusterScope(name string, labels map[string]string, fluentbit *loggingv1alpha1.Fluentbit) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   name,
		Labels: labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: fluentbit.APIVersion,
				Kind:       fluentbit.Kind,
				Name:       fluentbit.Name,
				UID:        fluentbit.UID,
			},
		},
	}
	return o
}
