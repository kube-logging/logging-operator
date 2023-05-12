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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FluentbitObjectMeta creates an objectMeta for resource fluentbit
func (r *Reconciler) FluentbitObjectMeta(name string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:      r.nameProvider.ComponentName(name),
		Namespace: r.Logging.Spec.ControlNamespace,
		Labels:    r.getFluentBitLabels(),
		OwnerReferences: []metav1.OwnerReference{
			r.nameProvider.OwnerRef(),
		},
	}
	return *o.DeepCopy()
}

// FluentbitObjectMetaClusterScope creates an cluster scoped objectMeta for resource fluentbit
func (r *Reconciler) FluentbitObjectMetaClusterScope(name string) metav1.ObjectMeta {
	o := metav1.ObjectMeta{
		Name:   r.nameProvider.ComponentName(name),
		Labels: r.getFluentBitLabels(),
		OwnerReferences: []metav1.OwnerReference{
			r.nameProvider.OwnerRef(),
		},
	}
	return o
}
