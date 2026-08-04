// Copyright © 2025 Kube logging authors
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

package watch_selector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/e2e-framework/third_party/helm"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/internal/harness"
	"github.com/kube-logging/logging-operator/e2e/internal/wait"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const (
	ns      = "test"
	release = "e2e"

	// unmanagedNS is created by the fluent chart, not by the harness.
	unmanagedNS = "fluentd"
)

func TestWatchSelectors(t *testing.T) {
	env := harness.New(t).
		WithCluster("watch-selector").
		WithRelease(release).
		WithControlNamespace(ns).
		WithOperatorArgs("-enable-leader-election=true", "-watch-labeled-children=true", "-watch-labeled-secrets=true").
		Start()

	logging := &v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "watch-selector-test",
			Namespace: ns,
		},
		Spec: v1beta1.LoggingSpec{
			ControlNamespace: ns,
			FluentbitSpec: &v1beta1.FluentbitSpec{
				ConfigHotReload: &v1beta1.HotReload{
					Image: v1beta1.ImageSpec{
						Repository: common.ConfigReloaderRepo,
						Tag:        common.ConfigReloaderTag,
					},
				},
				BufferVolumeImage: v1beta1.ImageSpec{
					Repository: common.NodeExporterRepo,
					Tag:        common.NodeExporterTag,
				},
			},
			FluentdSpec: &v1beta1.FluentdSpec{
				Image: v1beta1.ImageSpec{
					Repository: common.FluentdImageRepo,
					Tag:        common.FluentdImageTag,
				},
				ConfigReloaderImage: v1beta1.ImageSpec{
					Repository: common.ConfigReloaderRepo,
					Tag:        common.ConfigReloaderTag,
				},
				BufferVolumeImage: v1beta1.ImageSpec{
					Repository: common.NodeExporterRepo,
					Tag:        common.NodeExporterTag,
				},
			},
		},
	}
	env.Create(logging)

	require.NoError(t, installFluentdSts(env))

	unmanagedSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "unmanaged-fluentd-secret",
			Namespace: ns,
			Labels: map[string]string{
				"app": "fluentd",
			},
		},
		Data: map[string][]byte{
			"key": []byte("value"),
		},
	}
	env.Create(unmanagedSecret)

	env.WaitForRunning(
		wait.Pod(ns, logging.Name+"-fluentd-0"),
		wait.Pod(unmanagedNS, "fluentd-0"),
	)

	deployedLogging := &v1beta1.Logging{}
	require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKeyFromObject(logging), deployedLogging))

	// The managed resources have to be owned by the logging resource.
	managedSts := &appsv1.StatefulSet{}
	require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKey{Namespace: ns, Name: deployedLogging.Name + "-fluentd"}, managedSts))
	requireOwnedBy(t, deployedLogging, metav1.GetControllerOf(managedSts))

	managedSecret := &corev1.Secret{}
	require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKey{Namespace: ns, Name: deployedLogging.Name + "-fluentd"}, managedSecret))
	requireOwnedBy(t, deployedLogging, metav1.GetControllerOf(managedSecret))

	// The unmanaged ones have to be left alone.
	unmanagedSts := &appsv1.StatefulSet{}
	require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKey{Namespace: unmanagedNS, Name: "fluentd"}, unmanagedSts))
	require.Nil(t, metav1.GetControllerOf(unmanagedSts))

	secret := &corev1.Secret{}
	require.NoError(t, env.Client.Get(env.Ctx, client.ObjectKeyFromObject(unmanagedSecret), secret))
	require.Nil(t, metav1.GetControllerOf(secret))
}

func requireOwnedBy(t *testing.T, owner *v1beta1.Logging, ref *metav1.OwnerReference) {
	t.Helper()

	require.NotNil(t, ref)
	require.Equal(t, owner.APIVersion, ref.APIVersion)
	require.Equal(t, owner.Kind, ref.Kind)
	require.Equal(t, owner.Name, ref.Name)
	require.True(t, *ref.Controller)
}

func installFluentdSts(env *harness.Env) error {
	manager := helm.New(env.Cluster.KubeConfigFilePath())

	if err := manager.RunRepo(helm.WithArgs("add", "fluent", "https://fluent.github.io/helm-charts")); err != nil {
		return fmt.Errorf("failed to add fluent repo: %v", err)
	}

	if err := manager.RunInstall(
		helm.WithName("fluentd"),
		helm.WithChart("fluent/fluentd"),
		helm.WithArgs("--create-namespace"),
		helm.WithNamespace(unmanagedNS),
		helm.WithArgs("--set", "kind=StatefulSet"),
		helm.WithWait(),
	); err != nil {
		return fmt.Errorf("failed to install fluentd: %v", err)
	}

	return nil
}
