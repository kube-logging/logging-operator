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
	"github.com/banzaicloud/logging-operator/pkg/resources/templates"
	"github.com/banzaicloud/logging-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

type ConfigCheckResult struct {
	Valid bool
	Ready bool
}

const ConfigKey = "generated.conf"

func (r *Reconciler) appconfigMap() runtime.Object {
	data := make(map[string][]byte)
	data[AppConfigKey] = []byte(*r.config)
	return &corev1.Secret{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(AppSecretConfigName), util.MergeLabels(r.Logging.Labels, labelSelector), r.Logging),
		Data: data,
	}
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

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: pod.Namespace, Name: pod.Name}, pod)
	if err == nil {
		// check pod status and write into the configmap
		switch pod.Status.Phase {
		case v1.PodSucceeded:
			return &ConfigCheckResult{
				Valid: true,
				Ready: true,
			}, nil
		case v1.PodPending:
			fallthrough
		case v1.PodRunning:
			return &ConfigCheckResult{}, nil
		case v1.PodFailed:
			return &ConfigCheckResult{
				Ready: true,
				Valid: false,
			}, nil
		default:
			return nil, errors.Errorf("invalid pod status %s, unable to a validate config", pod.Status.Phase)
		}
	}

	if err != nil && !apierrors.IsNotFound(err) {
		return nil, errors.WrapIff(err, "failed to get configcheck pod %s:%s", pod.Namespace, pod.Name)
	}

	checkSecret := r.newCheckSecret(hashKey)

	err = r.Client.Create(context.TODO(), checkSecret)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return nil, errors.WrapIf(err, "failed to create secret for fluentd configcheck")
	}

	err = r.Client.Create(context.TODO(), pod)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to create pod for fluentd configcheck")
	}

	return &ConfigCheckResult{}, nil
}

func (r *Reconciler) configCheckCleanup(currentHash string) ([]string, error) {
	var multierr error
	var removedHashes = make([]string, 0)
	for configHash, _ := range r.Logging.Status.ConfigCheckResults {
		if configHash == currentHash {
			continue
		}
		if err := r.Client.Delete(context.TODO(), r.newCheckSecret(configHash)); err != nil {
			if !apierrors.IsNotFound(err) {
				multierr = errors.Combine(multierr,
					errors.Wrapf(err, "failed to remove config check secret %s", configHash))
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
	return removedHashes, multierr
}

func (r *Reconciler) newCheckSecret(hashKey string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-%s", hashKey)),
			util.MergeLabels(r.Logging.Labels, labelSelector),
			r.Logging,
		),
		Data: map[string][]byte{
			ConfigKey: []byte(*r.config),
		},
	}
}

func (r *Reconciler) newCheckPod(hashKey string) *v1.Pod {
	pod := &v1.Pod{
		ObjectMeta: templates.FluentdObjectMeta(
			r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-%s", hashKey)),
			util.MergeLabels(r.Logging.Labels, labelSelector),
			r.Logging,
		),
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Volumes: []v1.Volume{
				{
					Name: "config",
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: r.Logging.QualifiedName(fmt.Sprintf("fluentd-configcheck-%s", hashKey)),
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name: "fluentd",
					Image: fmt.Sprintf("%s:%s",
						r.Logging.Spec.FluentdSpec.Image.Repository, r.Logging.Spec.FluentdSpec.Image.Tag),
					ImagePullPolicy: v1.PullPolicy(r.Logging.Spec.FluentdSpec.Image.PullPolicy),
					Args: []string{
						"fluentd", "-c",
						fmt.Sprintf("/fluentd/etc/%s", ConfigKey),
						"--dry-run",
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "config",
							MountPath: "/fluentd/etc/",
						},
					},
				},
			},
		},
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
	return pod
}
