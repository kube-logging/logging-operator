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

package fluentd

import (
	"context"
	"fmt"
	"hash/fnv"

	"emperror.dev/errors"
	"github.com/cisco-open/operator-tools/pkg/reconciler"
	"github.com/kube-logging/logging-operator/pkg/compression"
	"github.com/kube-logging/logging-operator/pkg/resources/configcheck"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigCheckResult struct {
	Valid   bool
	Ready   bool
	Message string
}

func (r *Reconciler) appConfigSecret() (runtime.Object, reconciler.DesiredState, error) {
	data := make(map[string][]byte)

	if r.Logging.Spec.FluentdSpec.CompressConfigFile {
		AppConfigKeyCompress := AppConfigKey + ".gz"
		data[AppConfigKeyCompress] = compression.CompressString(*r.config, r.Log)
	} else {
		data[AppConfigKey] = []byte(*r.config)
	}

	return &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(AppSecretConfigName, ComponentFluentd),
		Data:       data,
	}, reconciler.StatePresent, nil
}

func (r *Reconciler) configHash() (string, error) {
	hasher := fnv.New32()
	_, err := hasher.Write([]byte(*r.config))
	if err != nil {
		return "", errors.WrapIf(err, "failed to calculate hash for the configmap data")
	}
	return fmt.Sprintf("%x", hasher.Sum32()), nil
}

func (r *Reconciler) configCheck(ctx context.Context) (*ConfigCheckResult, error) {
	hashKey, err := r.configHash()
	if err != nil {
		return nil, err
	}

	checkSecret, err := r.newCheckSecret(hashKey)
	if err != nil {
		return nil, err
	}
	configcheck.WithHashLabel(checkSecret, hashKey)
	checkSecretAppConfig, err := r.newCheckSecretAppConfig(hashKey)
	if err != nil {
		return nil, err
	}
	configcheck.WithHashLabel(checkSecretAppConfig, hashKey)

	err = r.Client.Create(ctx, checkSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret for fluentd configcheck")
	}
	err = r.Client.Create(ctx, checkSecretAppConfig)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret with Generated for fluentd configcheck")
	}

	checkOutputSecret, err := r.newCheckOutputSecret(hashKey)
	if err != nil {
		return nil, err
	}
	configcheck.WithHashLabel(checkOutputSecret, hashKey)
	err = r.Client.Create(ctx, checkOutputSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create output secret for fluentd configcheck")
	}

	pod := r.newCheckPod(hashKey)
	configcheck.WithHashLabel(pod, hashKey)

	existingPods := &corev1.PodList{}
	err = r.Client.List(ctx, existingPods, client.MatchingLabels(pod.Labels))
	if err != nil {
		return nil, errors.WrapIf(err, "failed to list existing configcheck pods")
	}

	podsByPhase := make(map[corev1.PodPhase]int)
	for _, p := range existingPods.Items {
		podsByPhase[p.Status.Phase] += 1
	}

	if podsByPhase[corev1.PodPending] > 0 {
		return &ConfigCheckResult{
			Ready:   false,
			Message: "there are pending configcheck pods, need to back off",
		}, nil
	}
	if podsByPhase[corev1.PodRunning] > 0 {
		return &ConfigCheckResult{
			Ready:   false,
			Message: "there are running configcheck pods, need to back off",
		}, nil
	}

	err = r.Client.Get(ctx, types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, pod)
	if err == nil {
		// check pod status and write into the configmap
		switch pod.Status.Phase {
		case corev1.PodSucceeded:
			return &ConfigCheckResult{
				Valid: true,
				Ready: true,
			}, nil
		case corev1.PodPending:
			fallthrough
		case corev1.PodRunning:
			return &ConfigCheckResult{}, nil
		case corev1.PodFailed:
			return &ConfigCheckResult{
				Ready: true,
				Valid: false,
			}, nil
		case corev1.PodUnknown:
			fallthrough
		default:
			return nil, errors.Errorf("invalid pod status %s, unable to a validate config", pod.Status.Phase)
		}
	}

	if err != nil && !apierrors.IsNotFound(err) {
		return nil, errors.WrapIff(err, "failed to get configcheck pod %s:%s", pod.Namespace, pod.Name)
	}

	err = r.Client.Create(ctx, pod)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create pod for fluentd configcheck")
	}

	return &ConfigCheckResult{}, nil
}

func (r *Reconciler) newCheckSecret(hashKey string) (*corev1.Secret, error) {
	data, err := r.generateConfigSecret()
	if err != nil {
		return nil, err
	}
	if r.Logging.Spec.FluentdSpec.CompressConfigFile {
		ConfigCheckKeyCompress := ConfigCheckKey + ".gz"
		data[ConfigCheckKeyCompress] = compression.CompressString(*r.config, r.Log)
	} else {
		data[ConfigCheckKey] = []byte(*r.config)
	}
	data["fluent.conf"] = []byte(fluentdConfigCheckTemplate)
	return &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(fmt.Sprintf("fluentd-configcheck-%s", hashKey), ComponentConfigCheck),
		Data:       data,
	}, nil
}

func (r *Reconciler) newCheckSecretAppConfig(hashKey string) (*corev1.Secret, error) {
	data := make(map[string][]byte)

	if r.Logging.Spec.FluentdSpec.CompressConfigFile {
		ConfigCheckKeyCompress := ConfigCheckKey + ".gz"
		data[ConfigCheckKeyCompress] = compression.CompressString(*r.config, r.Log)
	} else {
		data[ConfigCheckKey] = []byte(*r.config)
	}
	return &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(fmt.Sprintf("fluentd-configcheck-app-%s", hashKey), ComponentConfigCheck),
		Data:       data,
	}, nil
}

func (r *Reconciler) newCheckOutputSecret(hashKey string) (*corev1.Secret, error) {
	obj, _, err := r.outputSecret(r.secrets, OutputSecretPath)
	if err != nil {
		return nil, err
	}
	if secret, ok := obj.(*corev1.Secret); ok {
		secret.ObjectMeta = r.FluentdObjectMeta(fmt.Sprintf("fluentd-configcheck-output-%s", hashKey), ComponentConfigCheck)
		return secret, nil
	}
	return nil, errors.New("output secret is invalid, unable to create output secret for config check")
}

func (r *Reconciler) newCheckPod(hashKey string) *corev1.Pod {

	volumes := r.volumesCheckPod(hashKey)
	container := r.containerCheckPod(hashKey)
	initContainer := r.initContainerCheckPod()

	pod := &corev1.Pod{
		ObjectMeta: r.FluentdObjectMeta(fmt.Sprintf("fluentd-configcheck-%s", hashKey), ComponentConfigCheck),
		Spec: corev1.PodSpec{
			RestartPolicy:      corev1.RestartPolicyNever,
			ServiceAccountName: r.getServiceAccount(),
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
			Volumes:          volumes,
			ImagePullSecrets: r.Logging.Spec.FluentdSpec.Image.ImagePullSecrets,
			InitContainers:   initContainer,
			Containers:       container,
		},
	}
	if r.Logging.Spec.FluentdSpec.ConfigCheckAnnotations != nil {
		pod.Annotations = r.Logging.Spec.FluentdSpec.ConfigCheckAnnotations
	}
	if r.Logging.Spec.FluentdSpec.TLS.Enabled {
		tlsVolume := corev1.Volume{
			Name: "fluentd-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.FluentdSpec.TLS.SecretName,
				},
			},
		}
		pod.Spec.Volumes = append(pod.Spec.Volumes, tlsVolume)
		volumeMount := corev1.VolumeMount{
			Name:      "fluentd-tls",
			MountPath: "/fluentd/tls/",
		}
		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, volumeMount)
	}
	for _, n := range r.Logging.Spec.FluentdSpec.ExtraVolumes {
		if err := n.ApplyVolumeForPodSpec(&pod.Spec); err != nil {
			r.Log.Error(err, "Fluentd Config check pod extraVolume attachment failed.")
		}
	}

	return pod
}

func (r *Reconciler) volumesCheckPod(hashKey string) (v []corev1.Volume) {
	v = []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-%s", hashKey)),
				},
			},
		},
		{
			Name: "output-secret",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-output-%s", hashKey)),
				},
			},
		},
	}

	if r.Logging.Spec.FluentdSpec.CompressConfigFile {
		v = append(v, corev1.Volume{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		})
		v = append(v, corev1.Volume{
			Name: "app-config-compress",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-app-%s", hashKey)),
				},
			},
		})
	} else {
		v = append(v, corev1.Volume{
			Name: "app-config",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-app-%s", hashKey)),
				},
			},
		})
	}

	return v
}

func (r *Reconciler) containerCheckPod(hashKey string) []corev1.Container {
	containerArgs := []string{
		"timeout", "10",
		"fluentd", "-c",
		fmt.Sprintf("/fluentd/etc/%s", ConfigKey),
	}
	containerArgs = append(containerArgs, r.Logging.Spec.FluentdSpec.ExtraArgs...)

	container := []corev1.Container{
		{
			Name:            "fluentd",
			Image:           r.Logging.Spec.FluentdSpec.Image.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Image.PullPolicy),
			Args:            containerArgs,
			Env:             r.Logging.Spec.FluentdSpec.EnvVars,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "config",
					MountPath: "/fluentd/etc/",
				},
				{
					Name:      "app-config",
					MountPath: "/fluentd/app-config/",
				},
				{
					Name:      "output-secret",
					MountPath: OutputSecretPath,
				},
			},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:                r.Logging.Spec.FluentdSpec.Security.SecurityContext.RunAsUser,
				RunAsGroup:               r.Logging.Spec.FluentdSpec.Security.SecurityContext.RunAsGroup,
				ReadOnlyRootFilesystem:   r.Logging.Spec.FluentdSpec.Security.SecurityContext.ReadOnlyRootFilesystem,
				AllowPrivilegeEscalation: r.Logging.Spec.FluentdSpec.Security.SecurityContext.AllowPrivilegeEscalation,
				Privileged:               r.Logging.Spec.FluentdSpec.Security.SecurityContext.Privileged,
				RunAsNonRoot:             r.Logging.Spec.FluentdSpec.Security.SecurityContext.RunAsNonRoot,
				SELinuxOptions:           r.Logging.Spec.FluentdSpec.Security.SecurityContext.SELinuxOptions,
			},
			Resources: r.Logging.Spec.FluentdSpec.ConfigCheckResources,
		},
	}

	return container
}

func (r *Reconciler) initContainerCheckPod() []corev1.Container {
	var initContainer []corev1.Container
	if r.Logging.Spec.FluentdSpec.CompressConfigFile {
		initContainer = []corev1.Container{
			{
				Name:            "config-reloader",
				Image:           r.Logging.Spec.FluentdSpec.ConfigReloaderImage.RepositoryWithTag(),
				ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Image.PullPolicy),
				Resources:       r.Logging.Spec.FluentdSpec.ConfigReloaderResources,
				Args: []string{
					"--init-mode=true",
					"--volume-dir-archive=/tmp/archive",
					"--dir-for-unarchive=/fluentd/app-config",
					"-webhook-url=http://127.0.0.1:24444/api/config.reload",
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      "app-config",
						MountPath: "/fluentd/app-config",
					},
					{
						Name:      "app-config-compress",
						MountPath: "/tmp/archive",
					},
				},
			},
		}
	} else {
		initContainer = []corev1.Container{}
	}

	return initContainer
}
