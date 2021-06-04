// Copyright © 2021 Banzai Cloud
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

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

func (r *Reconciler) drainJobFor(pvc corev1.PersistentVolumeClaim) (*batchv1.Job, error) {
	bufVolName := r.Logging.QualifiedName(r.Logging.Spec.FluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName)

	var initContainers []corev1.Container
	if c := r.volumeMountHackContainer(); c != nil {
		initContainers = append(initContainers, *c)
	}

	fluentdContainer := r.fluentContainer() // TODO: don't redirect container logs
	fluentdContainer.VolumeMounts = append(fluentdContainer.VolumeMounts, corev1.VolumeMount{
		Name:      bufVolName,
		MountPath: bufferPath,
	})
	containers := []corev1.Container{
		*fluentdContainer,
		r.drainWatchContainer(bufVolName),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	spec := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: r.generatePodMeta(),
			Spec: corev1.PodSpec{
				Volumes:            r.generateVolume(),
				ServiceAccountName: r.getServiceAccount(),
				InitContainers:     initContainers,
				ImagePullSecrets:   r.Logging.Spec.FluentdSpec.Image.ImagePullSecrets,
				Containers:         containers,
				NodeSelector:       r.Logging.Spec.FluentdSpec.NodeSelector,
				Tolerations:        r.Logging.Spec.FluentdSpec.Tolerations,
				Affinity:           r.Logging.Spec.FluentdSpec.Affinity,
				PriorityClassName:  r.Logging.Spec.FluentdSpec.PodPriorityClassName,
				SecurityContext: &corev1.PodSecurityContext{
					RunAsNonRoot: r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsNonRoot,
					FSGroup:      r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.FSGroup,
					RunAsUser:    r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsUser,
					RunAsGroup:   r.Logging.Spec.FluentdSpec.Security.PodSecurityContext.RunAsGroup,
				},
				RestartPolicy: corev1.RestartPolicyNever,
			},
		},
	}

	spec.Template.Spec.Volumes = append(spec.Template.Spec.Volumes, corev1.Volume{
		Name: bufVolName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc.Name,
			},
		},
	})
	return &batchv1.Job{
		ObjectMeta: r.FluentdObjectMeta(StatefulSetName+pvc.Name[strings.LastIndex(pvc.Name, "-"):]+"-drain", ComponentFluentd),
		Spec:       spec,
	}, nil
}

func (r *Reconciler) drainWatchContainer(bufferVolumeName string) corev1.Container {
	return corev1.Container{
		Env: []corev1.EnvVar{
			{
				Name:  "BUFFER_PATH",
				Value: bufferPath,
			},
		},
		Image:           r.Logging.Spec.FluentdSpec.Scaling.DrainWatch.Image.RepositoryWithTag(),
		ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Scaling.DrainWatch.Image.PullPolicy),
		Name:            "drain-watch",
		VolumeMounts: []corev1.VolumeMount{
			{
				MountPath: bufferPath,
				Name:      bufferVolumeName,
				ReadOnly:  true,
			},
		},
	}
}
