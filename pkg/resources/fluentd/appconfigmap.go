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
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
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
	data[AppConfigKey] = []byte(*r.config)
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

func (r *Reconciler) configCheck() (*ConfigCheckResult, error) {
	hashKey, err := r.configHash()
	if err != nil {
		return nil, err
	}

	pod := r.newCheckPod(hashKey)

	existingPods := &corev1.PodList{}
	err = r.Client.List(context.TODO(), existingPods, client.MatchingLabels(pod.Labels))
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

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, pod)
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

	checkSecret, err := r.newCheckSecret(hashKey)
	if err != nil {
		return nil, err
	}
	checkOutputSecret, err := r.newCheckOutputSecret(hashKey)
	if err != nil {
		return nil, err
	}

	err = r.Client.Create(context.TODO(), checkSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret for fluentd configcheck")
	}
	err = r.Client.Create(context.TODO(), checkOutputSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create output secret for fluentd configcheck")
	}

	err = r.Client.Create(context.TODO(), pod)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create pod for fluentd configcheck")
	}

	return &ConfigCheckResult{}, nil
}

func (r *Reconciler) configCheckCleanup(currentHash string) (removedHashes []string, multierr error) {
	for configHash := range r.Logging.Status.ConfigCheckResults {
		if configHash == currentHash {
			continue
		}
		newSecret, err := r.newCheckSecret(configHash)
		if err != nil {
			multierr = errors.Combine(multierr,
				errors.Wrapf(err, "failed to create config check secret %s", configHash))
			continue
		}
		if err := r.Client.Delete(context.TODO(), newSecret); err != nil {
			if !apierrors.IsNotFound(err) {
				multierr = errors.Combine(multierr,
					errors.Wrapf(err, "failed to remove config check secret %s", configHash))
				continue
			}
		}
		checkOutputSecret, err := r.newCheckOutputSecret(configHash)
		if err != nil {
			multierr = errors.Combine(multierr,
				errors.Wrapf(err, "failed to create config check output secret %s", configHash))
			continue
		}
		if err := r.Client.Delete(context.TODO(), checkOutputSecret); err != nil {
			if !apierrors.IsNotFound(err) {
				multierr = errors.Combine(multierr,
					errors.Wrapf(err, "failed to remove config check output secret %s", configHash))
				continue
			}
		}
		if err := r.Client.Delete(context.TODO(), r.newCheckPod(configHash)); err != nil {
			if !apierrors.IsNotFound(err) {
				multierr = errors.Combine(multierr,
					errors.Wrapf(err, "failed to remove config check pod %s", configHash))
				continue
			}
		}
		removedHashes = append(removedHashes, configHash)
	}
	return
}

func (r *Reconciler) newCheckSecret(hashKey string) (*corev1.Secret, error) {
	data, err := r.generateConfigSecret()
	if err != nil {
		return nil, err
	}
	data[ConfigCheckKey] = []byte(*r.config)
	data["fluent.conf"] = []byte(fluentdConfigCheckTemplate)
	return &corev1.Secret{
		ObjectMeta: r.FluentdObjectMeta(fmt.Sprintf("fluentd-configcheck-%s", hashKey), ComponentConfigCheck),
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
			Volumes: []corev1.Volume{
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
			},
			ImagePullSecrets: r.Logging.Spec.FluentdSpec.Image.ImagePullSecrets,
			Containers: []corev1.Container{
				{
					Name:            "fluentd",
					Image:           r.Logging.Spec.FluentdSpec.Image.RepositoryWithTag(),
					ImagePullPolicy: corev1.PullPolicy(r.Logging.Spec.FluentdSpec.Image.PullPolicy),
					Args: []string{
						"fluentd", "-c",
						fmt.Sprintf("/fluentd/etc/%s", ConfigKey),
						"--dry-run",
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config",
							MountPath: "/fluentd/etc/",
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
			},
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
