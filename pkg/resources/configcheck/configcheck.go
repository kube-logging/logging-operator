// Copyright © 2022 Kube logging authors
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

package configcheck

import (
	"context"

	"emperror.dev/errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const HashLabel = "logging.banzaicloud.io/config-hash"

func WithHashLabel(accessor v1.Object, hash string) {
	l := accessor.GetLabels()
	if l == nil {
		l = map[string]string{}
	}
	l[HashLabel] = hash
	accessor.SetLabels(l)
}

func hasHashLabel(accessor v1.Object, hash string) (has bool, match bool) {
	l := accessor.GetLabels()
	var val string
	val, has = l[HashLabel]
	return has, val == hash
}

type ConfigCheckCleaner struct {
	client client.Client
	labels client.MatchingLabels
}

func NewConfigCheckCleaner(c client.Client, component string) *ConfigCheckCleaner {
	return &ConfigCheckCleaner{
		client: c,
		labels: client.MatchingLabels{
			"app.kubernetes.io/component": component,
		},
	}
}

// SecretCleanup cleans up configcheck secrets that have the logging.banzaicloud.io/config-hash label, but
// doesn't match the current config hash
func (c *ConfigCheckCleaner) SecretCleanup(ctx context.Context, hash string) (multierr error) {
	allCheckSecrets := &corev1.SecretList{}
	if err := c.client.List(ctx, allCheckSecrets, c.labels); err != nil {
		return errors.Wrap(err, "failed to list configcheck secrets")
	}

	for _, secret := range allCheckSecrets.Items {
		if _, match := hasHashLabel(&secret, hash); match {
			continue
		}
		if err := client.IgnoreNotFound(c.client.Delete(ctx, &secret)); err != nil {
			multierr = errors.Combine(multierr,
				errors.Wrapf(err, "failed to remove config check secret %s", secret.Name))
			continue
		}
	}

	return
}

// PodCleanup cleans up configcheck pods that have the logging.banzaicloud.io/config-hash label, but
// doesn't match the current config hash
func (c *ConfigCheckCleaner) PodCleanup(ctx context.Context, hash string) (multierr error) {
	allCheckPods := &corev1.PodList{}
	if err := c.client.List(ctx, allCheckPods, c.labels); err != nil {
		return errors.Wrap(err, "failed to list configcheck pods")
	}

	for _, pod := range allCheckPods.Items {
		if _, match := hasHashLabel(&pod, hash); match {
			continue
		}
		if err := client.IgnoreNotFound(c.client.Delete(ctx, &pod)); err != nil {
			multierr = errors.Combine(multierr,
				errors.Wrapf(err, "failed to remove config check pod %s", pod.Name))
			continue
		}
	}

	return
}
