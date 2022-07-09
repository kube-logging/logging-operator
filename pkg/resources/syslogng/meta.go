// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	util "github.com/banzaicloud/operator-tools/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SyslogNGObjectMeta creates an objectMeta for resource syslog-ng
func (r *Reconciler) SyslogNGObjectMeta(name, component string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      r.Logging.QualifiedName(name),
		Namespace: r.Logging.Spec.ControlNamespace,
		Labels:    r.Logging.GetSyslogNGLabels(component),
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: r.Logging.APIVersion,
				Kind:       r.Logging.Kind,
				Name:       r.Logging.Name,
				UID:        r.Logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return *o.DeepCopy()
}

// SyslogNGObjectMetaClusterScope creates an objectMeta for resource syslog-ng
func (r *Reconciler) SyslogNGObjectMetaClusterScope(name, component string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   r.Logging.QualifiedName(name),
		Labels: r.Logging.GetSyslogNGLabels(component),
		OwnerReferences: []metav1.OwnerReference{
			{
				APIVersion: r.Logging.APIVersion,
				Kind:       r.Logging.Kind,
				Name:       r.Logging.Name,
				UID:        r.Logging.UID,
				Controller: util.BoolPointer(true),
			},
		},
	}
	return o
}
