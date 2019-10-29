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

package templates

import (
	"github.com/banzaicloud/logging-operator/api/v1beta1"
	"github.com/banzaicloud/logging-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FluentbitObjectMeta creates an objectMeta for resource fluentbit
func FluentbitObjectMeta(name string, labels map[string]string, logging *v1beta1.Logging) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      name,
		Namespace: logging.Spec.ControlNamespace,
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: logging.APIVersion,
				Kind:       logging.Kind,
				Name:       logging.Name,
				UID:        logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}

// FluentbitObjectMetaClusterScope creates an cluster scoped objectMeta for resource fluentbit
func FluentbitObjectMetaClusterScope(name string, labels map[string]string, logging *v1beta1.Logging) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   name,
		Labels: labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: logging.APIVersion,
				Kind:       logging.Kind,
				Name:       logging.Name,
				UID:        logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}

// FluentdObjectMeta creates an objectMeta for resource fluentd
func FluentdObjectMeta(name string, labels map[string]string, logging *v1beta1.Logging) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      name,
		Namespace: logging.Spec.ControlNamespace,
		Labels:    labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: logging.APIVersion,
				Kind:       logging.Kind,
				Name:       logging.Name,
				UID:        logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}

// FluentdObjectMetaClusterScope creates an objectMeta for resource fluentd
func FluentdObjectMetaClusterScope(name string, labels map[string]string, logging *v1beta1.Logging) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   name,
		Labels: labels,
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: logging.APIVersion,
				Kind:       logging.Kind,
				Name:       logging.Name,
				UID:        logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}
