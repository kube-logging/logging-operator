// Copyright © 2022 Banzai Cloud
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

package syslogng

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/banzaicloud/operator-tools/pkg/merge"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigCheckResult struct {
	Valid   bool
	Ready   bool
	Message string
}

func (r *Reconciler) configHash() (string, error) {
	hasher := fnv.New32()
	if _, err := io.WriteString(hasher, r.config); err != nil {
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
	err = r.Client.Create(ctx, checkSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret for syslog-ng configcheck")
	}

	checkOutputSecret, err := r.newCheckOutputSecret(hashKey)
	if err != nil {
		return nil, err
	}
	err = r.Client.Create(ctx, checkOutputSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create output secret for syslog-ng configcheck")
	}

	pod, err := r.newCheckPod(hashKey)
	if err != nil {
		return nil, errors.WrapIff(err, "failed to create resource description for config check pod %s", hashKey)
	}

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
		return nil, errors.WrapIf(err, "failed to create pod for syslog-ng configcheck")
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
		checkPod, err := r.newCheckPod(configHash)
		if err != nil {
			multierr = errors.Append(multierr, errors.WrapIff(err, "failed to create resource description for config check pod %s", configHash))
			continue
		}
		if err := r.Client.Delete(context.TODO(), checkPod); err != nil {
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
	return &corev1.Secret{
		ObjectMeta: r.SyslogNGObjectMeta(fmt.Sprintf("syslog-ng-configcheck-%s", hashKey), ComponentConfigCheck),
		Data: map[string][]byte{
			configKey: []byte(r.config),
		},
	}, nil
}

func (r *Reconciler) newCheckOutputSecret(hashKey string) (*corev1.Secret, error) {
	obj, _, err := r.outputSecret(r.secrets, outputSecretPath)
	if err != nil {
		return nil, err
	}
	if secret, ok := obj.(*corev1.Secret); ok {
		secret.ObjectMeta = r.SyslogNGObjectMeta(fmt.Sprintf("syslog-ng-configcheck-output-%s", hashKey), ComponentConfigCheck)
		return secret, nil
	}
	return nil, errors.New("output secret is invalid, unable to create output secret for config check")
}

func (r *Reconciler) newCheckPod(hashKey string) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: r.SyslogNGObjectMeta(fmt.Sprintf("syslog-ng-configcheck-%s", hashKey), ComponentConfigCheck),
		Spec: corev1.PodSpec{
			RestartPolicy:      corev1.RestartPolicyNever,
			ServiceAccountName: r.getServiceAccountName(),
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: r.Logging.QualifiedName(fmt.Sprintf("syslog-ng-configcheck-%s", hashKey)),
						},
					},
				},
				{
					Name: "output-secret",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: r.Logging.QualifiedName(fmt.Sprintf("syslog-ng-configcheck-output-%s", hashKey)),
						},
					},
				},
			},
			Containers: []corev1.Container{
				{
					Name:            "syslog-ng",
					Image:           v1beta1.RepositoryWithTag(syslogngImageRepository, syslogngImageTag),
					ImagePullPolicy: corev1.PullIfNotPresent,
					Args: []string{
						"-s",
					},
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("400M"),
							corev1.ResourceCPU:    resource.MustParse("1000m"),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("100M"),
							corev1.ResourceCPU:    resource.MustParse("500m"),
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config",
							MountPath: configDir,
						},
						{
							Name:      "output-secret",
							MountPath: outputSecretPath,
						},
					},
				},
			},
		},
	}
	if r.Logging.Spec.SyslogNGSpec.TLS.Enabled {
		tlsVolume := corev1.Volume{
			Name: "syslog-ng-tls",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: r.Logging.Spec.SyslogNGSpec.TLS.SecretName,
				},
			},
		}
		pod.Spec.Volumes = append(pod.Spec.Volumes, tlsVolume)
		volumeMount := corev1.VolumeMount{
			Name:      "syslog-ng-tls",
			MountPath: "/syslog-ng/tls/",
		}
		pod.Spec.Containers[0].VolumeMounts = append(pod.Spec.Containers[0].VolumeMounts, volumeMount)
	}

	err := merge.Merge(pod, r.Logging.Spec.SyslogNGSpec.ConfigCheckPodOverrides)

	return pod, err
}
