// Copyright © 2021 Cisco Systems, Inc. and/or its affiliates
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
	"strings"

	"github.com/cisco-open/operator-tools/pkg/utils"
	corev1 "k8s.io/api/core/v1"
)

func (r *Reconciler) placeholderPodFor(pvc corev1.PersistentVolumeClaim) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: r.FluentdObjectMeta(StatefulSetName+pvc.Name[strings.LastIndex(pvc.Name, "-"):], ComponentPlaceholder),
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "pause",
					Image:           r.Logging.Spec.FluentdSpec.Scaling.Drain.PauseImage.RepositoryWithTag(),
					ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Scaling.Drain.PauseImage.PullPolicy),
				},
			},
			NodeSelector:                  r.Logging.Spec.FluentdSpec.NodeSelector,
			Tolerations:                   r.Logging.Spec.FluentdSpec.Tolerations,
			Affinity:                      r.Logging.Spec.FluentdSpec.Affinity,
			PriorityClassName:             r.Logging.Spec.FluentdSpec.PodPriorityClassName,
			RestartPolicy:                 corev1.RestartPolicyNever,
			TerminationGracePeriodSeconds: utils.IntPointer64(0), // terminate immediately
		},
	}
}
