// Copyright © 2023 Kube logging authors
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
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) pdb() (runtime.Object, reconciler.DesiredState, error) {
	if r.fluentdSpec.Pdb != nil {
		pdbSpec := r.fluentdSpec.Pdb

		return &policyv1.PodDisruptionBudget{
			ObjectMeta: r.FluentdObjectMeta(PodDisruptionBudgetName, ComponentFluentd),
			Spec: policyv1.PodDisruptionBudgetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: r.Logging.GetFluentdLabels(ComponentFluentd, *r.fluentdSpec),
				},
				MinAvailable:               pdbSpec.MinAvailable,
				MaxUnavailable:             pdbSpec.MaxUnavailable,
				UnhealthyPodEvictionPolicy: pdbSpec.UnhealthyPodEvictionPolicy,
			},
		}, reconciler.StatePresent, nil
	}
	return &policyv1.PodDisruptionBudget{
		ObjectMeta: r.FluentdObjectMeta(PodDisruptionBudgetName, ComponentFluentd),
		Spec:       policyv1.PodDisruptionBudgetSpec{},
	}, reconciler.StateAbsent, nil
}
