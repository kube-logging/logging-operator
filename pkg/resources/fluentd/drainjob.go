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

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func (r *Reconciler) drainerJobFor(pvc corev1.PersistentVolumeClaim, fluentdSpec v1beta1.FluentdSpec) (*batchv1.Job, error) {
	bufVolName := r.Logging.QualifiedName(fluentdSpec.BufferStorageVolume.PersistentVolumeClaim.PersistentVolumeSource.ClaimName)

	fluentdContainer := fluentContainer(withoutFluentOutLogrotate(&fluentdSpec))
	fluentdContainer.VolumeMounts = append(fluentdContainer.VolumeMounts, corev1.VolumeMount{
		Name:      bufVolName,
		MountPath: bufferPath,
	})
	containers := []corev1.Container{
		fluentdContainer,
		drainWatchContainer(&fluentdSpec.Scaling.Drain, bufVolName),
	}
	if c := r.bufferMetricsSidecarContainer(); c != nil {
		containers = append(containers, *c)
	}

	var initContainers []corev1.Container
	if i := generateInitContainer(fluentdSpec); i != nil {
		initContainers = append(initContainers, *i)
	}
	if c := r.tmpDirHackContainer(); c != nil {
		initContainers = append(initContainers, *c)
	}

	spec := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      r.getDrainerLabels(fluentdSpec),
				Annotations: fluentdSpec.Scaling.Drain.Annotations,
			},
			Spec: corev1.PodSpec{
				Volumes:                   r.generateVolume(),
				ServiceAccountName:        r.getServiceAccount(),
				ImagePullSecrets:          fluentdSpec.Image.ImagePullSecrets,
				InitContainers:            initContainers,
				Containers:                containers,
				NodeSelector:              fluentdSpec.NodeSelector,
				Tolerations:               fluentdSpec.Tolerations,
				Affinity:                  fluentdSpec.Affinity,
				TopologySpreadConstraints: fluentdSpec.TopologySpreadConstraints,
				PriorityClassName:         fluentdSpec.PodPriorityClassName,
				SecurityContext:           fluentdSpec.Security.PodSecurityContext,
				RestartPolicy:             corev1.RestartPolicyNever,
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
	for _, n := range fluentdSpec.ExtraVolumes {
		if err := n.ApplyVolumeForPodSpec(&spec.Template.Spec); err != nil {
			return nil, err
		}
	}
	return &batchv1.Job{
		ObjectMeta: r.FluentdObjectMeta(StatefulSetName+pvc.Name[strings.LastIndex(pvc.Name, "-"):]+"-drainer", ComponentDrainer),
		Spec:       spec,
	}, nil
}

func drainWatchContainer(cfg *v1beta1.FluentdDrainConfig, bufferVolumeName string) corev1.Container {
	return corev1.Container{
		Env: []corev1.EnvVar{
			{
				Name:  "BUFFER_PATH",
				Value: bufferPath,
			},
			{
				Name:  "CHECK_INTERVAL",
				Value: drainerCheckInterval,
			},
		},
		Image:           cfg.Image.RepositoryWithTag(),
		ImagePullPolicy: corev1.PullPolicy(cfg.Image.PullPolicy),
		Name:            "drain-watch",
		VolumeMounts: []corev1.VolumeMount{
			{
				MountPath: bufferPath,
				Name:      bufferVolumeName,
				ReadOnly:  true,
			},
		},
		Resources:       *cfg.Resources,
		SecurityContext: cfg.SecurityContext,
	}
}

func withoutFluentOutLogrotate(spec *v1beta1.FluentdSpec) *v1beta1.FluentdSpec {
	res := spec.DeepCopy()
	res.FluentOutLogrotate = nil
	return res
}

func (r *Reconciler) getDrainerLabels(fluentdSpec v1beta1.FluentdSpec) map[string]string {
	labels := r.Logging.GetFluentdLabels(ComponentDrainer, fluentdSpec)

	for key, value := range fluentdSpec.Scaling.Drain.Labels {
		labels[key] = value
	}

	return labels
}
