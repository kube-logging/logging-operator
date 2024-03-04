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
	"github.com/spf13/cast"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/compression"
	"github.com/kube-logging/logging-operator/pkg/resources/configcheck"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

type ConfigCheckResult struct {
	Valid   bool
	Ready   bool
	Message string
}

func (r *Reconciler) appConfigSecret() (runtime.Object, reconciler.DesiredState, error) {
	data := make(map[string][]byte)

	if r.fluentdSpec.CompressConfigFile {
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

func (r *Reconciler) hasConfigCheckPod(ctx context.Context, hashKey string, fluentdSpec v1beta1.FluentdSpec) (bool, error) {
	var err error
	pod := r.newCheckPod(hashKey, fluentdSpec)

	p := &corev1.Pod{}
	err = r.Client.Get(ctx, client.ObjectKeyFromObject(pod), p)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Reconciler) CheckForObjectExistence(ctx context.Context, object client.Object) (*ConfigCheckResult, error) {
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(object), object); err != nil {
		if apierrors.IsNotFound(err) {
			objNotFoundMsg := fmt.Sprintf("object %s (kind: secret) in namespace %s not found", object.GetName(), object.GetNamespace())
			r.Log.Info(objNotFoundMsg)
			err = nil
		}
		errMsg := fmt.Sprintf("object %s (kind: secret) in namespace %s is not available", object.GetName(), object.GetNamespace())
		return &ConfigCheckResult{
			Ready: false, Valid: false, Message: errMsg,
		}, err
	}
	return nil, nil
}

func (r *Reconciler) configCheck(ctx context.Context) (*ConfigCheckResult, error) {
	hashKey, err := r.configHash()
	if err != nil {
		return nil, err
	}

	checkSecret, err := r.newCheckSecret(hashKey, *r.fluentdSpec)
	if err != nil {
		return nil, err
	}
	configcheck.WithHashLabel(checkSecret, hashKey)
	checkSecretAppConfig, err := r.newCheckSecretAppConfig(hashKey, *r.fluentdSpec)
	if err != nil {
		return nil, err
	}
	configcheck.WithHashLabel(checkSecretAppConfig, hashKey)

	err = r.Client.Create(ctx, checkSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret for fluentd configcheck")
	}
	if res, err := r.CheckForObjectExistence(ctx, checkSecret); res != nil {
		return res, errors.WrapIf(err, "failed to find secret for fluentd configcheck")
	}

	err = r.Client.Create(ctx, checkSecretAppConfig)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret with Generated for fluentd configcheck")
	}
	if res, err := r.CheckForObjectExistence(ctx, checkSecretAppConfig); res != nil {
		return res, errors.WrapIf(err, "failed to find secret with Generated for fluentd configcheck")
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
	if res, err := r.CheckForObjectExistence(ctx, checkOutputSecret); res != nil {
		return res, errors.WrapIf(err, "failed to find output secret for fluentd configcheck")
	}

	pod := r.newCheckPod(hashKey, *r.fluentdSpec)
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

func (r *Reconciler) newCheckSecret(hashKey string, fluentdSpec v1beta1.FluentdSpec) (*corev1.Secret, error) {
	data, err := r.generateConfigSecret(fluentdSpec)
	if err != nil {
		return nil, err
	}
	if fluentdSpec.CompressConfigFile {
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

func (r *Reconciler) newCheckSecretAppConfig(hashKey string, fluentdSpec v1beta1.FluentdSpec) (*corev1.Secret, error) {
	data := make(map[string][]byte)

	if fluentdSpec.CompressConfigFile {
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

func (r *Reconciler) newCheckPod(hashKey string, fluentdSpec v1beta1.FluentdSpec) *corev1.Pod {

	volumes := r.volumesCheckPod(hashKey, fluentdSpec)
	container := r.containerCheckPod(hashKey, fluentdSpec)
	initContainer := r.initContainerCheckPod(fluentdSpec)

	pod := &corev1.Pod{
		ObjectMeta: r.configCheckPodObjectMeta(fmt.Sprintf("fluentd-configcheck-%s", hashKey), ComponentConfigCheck),
		Spec: corev1.PodSpec{
			RestartPolicy:      corev1.RestartPolicyNever,
			ServiceAccountName: r.getServiceAccount(),
			NodeSelector:       fluentdSpec.NodeSelector,
			Tolerations:        fluentdSpec.Tolerations,
			Affinity:           fluentdSpec.Affinity,
			PriorityClassName:  fluentdSpec.PodPriorityClassName,
			SecurityContext:    fluentdSpec.Security.PodSecurityContext,
			Volumes:            volumes,
			ImagePullSecrets:   fluentdSpec.Image.ImagePullSecrets,
			InitContainers:     initContainer,
			Containers:         container,
		},
	}
	if fluentdSpec.ConfigCheckAnnotations != nil {
		pod.Annotations = fluentdSpec.ConfigCheckAnnotations
	}
	if fluentdSpec.TLS.Enabled {
		tlsVolume := corev1.Volume{
			Name: "fluentd-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: fluentdSpec.TLS.SecretName,
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
	for _, n := range fluentdSpec.ExtraVolumes {
		if err := n.ApplyVolumeForPodSpec(&pod.Spec); err != nil {
			r.Log.Error(err, "Fluentd Config check pod extraVolume attachment failed.")
		}
	}

	return pod
}

func (r *Reconciler) volumesCheckPod(hashKey string, fluentdSpec v1beta1.FluentdSpec) (v []corev1.Volume) {
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

	if fluentdSpec.CompressConfigFile {
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

func (r *Reconciler) containerCheckPod(hashKey string, fluentdSpec v1beta1.FluentdSpec) []corev1.Container {
	var containerArgs []string

	switch r.Logging.Spec.ConfigCheck.Strategy {
	case v1beta1.ConfigCheckStrategyTimeout:
		containerArgs = []string{
			"timeout", cast.ToString(r.Logging.Spec.ConfigCheck.TimeoutSeconds),
			"fluentd", "-c",
			fmt.Sprintf("/fluentd/etc/%s", ConfigKey),
		}
	case v1beta1.ConfigCheckStrategyDryRun:
		fallthrough
	default:
		containerArgs = []string{
			"fluentd", "-c",
			fmt.Sprintf("/fluentd/etc/%s", ConfigKey),
			"--dry-run",
		}
	}

	containerArgs = append(containerArgs, fluentdSpec.ExtraArgs...)

	container := []corev1.Container{
		{
			Name:            "fluentd",
			Image:           fluentdSpec.Image.RepositoryWithTag(),
			ImagePullPolicy: corev1.PullPolicy(fluentdSpec.Image.PullPolicy),
			Args:            containerArgs,
			Env:             fluentdSpec.EnvVars,
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
			SecurityContext: fluentdSpec.Security.SecurityContext,
			Resources:       fluentdSpec.ConfigCheckResources,
		},
	}

	return container
}

func (r *Reconciler) initContainerCheckPod(fluentdSpec v1beta1.FluentdSpec) []corev1.Container {
	var initContainer []corev1.Container
	if fluentdSpec.CompressConfigFile {
		initContainer = []corev1.Container{
			{
				Name:            "config-reloader",
				Image:           fluentdSpec.ConfigReloaderImage.RepositoryWithTag(),
				ImagePullPolicy: corev1.PullPolicy(fluentdSpec.Image.PullPolicy),
				Resources:       fluentdSpec.ConfigReloaderResources,
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

func (r *Reconciler) configCheckPodObjectMeta(name, component string) metav1.ObjectMeta {
	objectMeta := r.FluentdObjectMeta(name, component)

	for key, value := range r.Logging.Spec.ConfigCheck.Labels {
		objectMeta.Labels[key] = value
	}

	return objectMeta
}
