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

package fluentd

import (
	"github.com/banzaicloud/logging-operator/pkg/k8sutil"
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	pspv1beta1 "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) clusterPodSecurityPolicy() (runtime.Object, k8sutil.DesiredState) {
	if r.Logging.Spec.FluentdSpec.Security.PodSecurityPolicyCreate {

		return &pspv1beta1.PodSecurityPolicy{
			ObjectMeta: templates.FluentdObjectMeta(r.Logging.QualifiedName(roleName), r.Logging.Labels, r.Logging),
			Spec: pspv1beta1.PodSecurityPolicySpec{
				Privileged:                      false,
				DefaultAddCapabilities:          nil,
				RequiredDropCapabilities:        nil,
				AllowedCapabilities:             nil,
				Volumes:                         nil,
				HostNetwork:                     false,
				HostPorts:                       nil,
				HostPID:                         false,
				HostIPC:                         false,
				SELinux:                         pspv1beta1.SELinuxStrategyOptions{},
				RunAsUser:                       pspv1beta1.RunAsUserStrategyOptions{},
				RunAsGroup:                      nil,
				SupplementalGroups:              pspv1beta1.SupplementalGroupsStrategyOptions{},
				FSGroup:                         pspv1beta1.FSGroupStrategyOptions{},
				ReadOnlyRootFilesystem:          false,
				DefaultAllowPrivilegeEscalation: nil,
				AllowPrivilegeEscalation:        nil,
				AllowedHostPaths:                nil,
				AllowedFlexVolumes:              nil,
				AllowedCSIDrivers:               nil,
				AllowedUnsafeSysctls:            nil,
				ForbiddenSysctls:                nil,
				AllowedProcMountTypes:           nil,
			},
		}, k8sutil.StatePresent

	}
	return nil, k8sutil.StatePresent
}
