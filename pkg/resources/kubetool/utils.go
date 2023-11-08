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

package kubetool

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func FindContainerByName(cnrs []corev1.Container, name string) *corev1.Container {
	for i := range cnrs {
		cnr := &cnrs[i]
		if cnr.Name == name {
			return cnr
		}
	}
	return nil
}

func FindVolumeByName(vols []corev1.Volume, name string) *corev1.Volume {
	for i := range vols {
		vol := &vols[i]
		if vol.Name == name {
			return vol
		}
	}
	return nil
}

func FindVolumeMountByName(mnts []corev1.VolumeMount, name string) *corev1.VolumeMount {
	for i := range mnts {
		mnt := &mnts[i]
		if mnt.Name == name {
			return mnt
		}
	}
	return nil
}

func JobSuccessfullyCompleted(job *batchv1.Job) bool {
	return job.Status.CompletionTime != nil && job.Status.Succeeded > 0
}
