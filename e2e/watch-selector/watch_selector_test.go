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
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
	"sigs.k8s.io/e2e-framework/third_party/helm"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/common/cond"
	"github.com/kube-logging/logging-operator/e2e/common/setup"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/stretchr/testify/require"
)

var TestTempDir string

func init() {
	var ok bool
	TestTempDir, ok = os.LookupEnv("PROJECT_DIR")
	if !ok {
		TestTempDir = "../.."
	}
	TestTempDir = filepath.Join(TestTempDir, "build/_test")
	err := os.MkdirAll(TestTempDir, os.FileMode(0755))
	if err != nil {
		panic(err)
	}
}

func TestWatchSelectors(t *testing.T) {
	common.Initialize(t)
	ns := "test"
	releaseNameOverride := "e2e"
	common.WithCluster("watch-selector", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
			options.Args = []string{"-enable-leader-election=true", "-watch-labeled-children=true", "-watch-labeled-secrets=true"}
		}))

		ctx := context.Background()

		// Managed logging resource which creates a fluentd pod with a secret named: watch-selector-test-fluentd
		logging := v1beta1.Logging{
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
		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))

		// Unmanaged resources
		common.RequireNoError(t, installFluentdSts(c))

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
		common.RequireNoError(t, c.GetClient().Create(ctx, unmanagedSecret))

		require.Eventually(t, func() bool {
			if isManagedFluentdPodRunning := cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: ns, Name: logging.Name + "-fluentd-0"}); !isManagedFluentdPodRunning() {
				t.Logf("managed fluentd pod is not running")
				return false
			}

			if isUnmanagedFluentdPodRunning := cond.PodShouldBeRunning(t, c.GetClient(), client.ObjectKey{Namespace: "fluentd", Name: "fluentd-0"}); !isUnmanagedFluentdPodRunning() {
				t.Logf("unmanaged fluentd pod is not running")
				return false
			}

			return true
		}, 5*time.Minute, 3*time.Second)

		deployedLogging := &v1beta1.Logging{}
		common.RequireNoError(t, c.GetClient().Get(ctx, client.ObjectKeyFromObject(&logging), deployedLogging))

		// Check if the managed resources are actually controlled by the logging resource
		managedSts := &appsv1.StatefulSet{}
		common.RequireNoError(t, c.GetClient().Get(ctx, client.ObjectKey{Namespace: ns, Name: deployedLogging.Name + "-fluentd"}, managedSts))
		stsOwnerRefMeta := metav1.GetControllerOf(managedSts)
		require.NotNil(t, stsOwnerRefMeta)

		require.Equal(t, deployedLogging.APIVersion, stsOwnerRefMeta.APIVersion)
		require.Equal(t, deployedLogging.Kind, stsOwnerRefMeta.Kind)
		require.Equal(t, deployedLogging.Name, stsOwnerRefMeta.Name)
		require.True(t, *stsOwnerRefMeta.Controller)

		managedSecret := &corev1.Secret{}
		common.RequireNoError(t, c.GetClient().Get(ctx, client.ObjectKey{Namespace: ns, Name: deployedLogging.Name + "-fluentd"}, managedSecret))
		secretOwnerRefMeta := metav1.GetControllerOf(managedSecret)
		require.NotNil(t, secretOwnerRefMeta)

		require.Equal(t, deployedLogging.APIVersion, secretOwnerRefMeta.APIVersion)
		require.Equal(t, deployedLogging.Kind, secretOwnerRefMeta.Kind)
		require.Equal(t, deployedLogging.Name, secretOwnerRefMeta.Name)
		require.True(t, *secretOwnerRefMeta.Controller)

		// Check if the unmanaged resources are actually not controlled by the operator
		unmanagedSts := &appsv1.StatefulSet{}
		common.RequireNoError(t, c.GetClient().Get(ctx, client.ObjectKey{Namespace: "fluentd", Name: "fluentd"}, unmanagedSts))
		secretOwnerRefMeta = metav1.GetControllerOf(unmanagedSts)
		require.Nil(t, secretOwnerRefMeta)

		secret := &corev1.Secret{}
		common.RequireNoError(t, c.GetClient().Get(ctx, client.ObjectKeyFromObject(unmanagedSecret), secret))
		secretOwnerRefMeta = metav1.GetControllerOf(secret)
		require.Nil(t, secretOwnerRefMeta)

	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		err := c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{ns, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
		if err != nil {
			return err
		}

		loggingOperatorName := "logging-operator-" + releaseNameOverride
		t.Logf("Collecting coverage files from logging-operator: %s/%s", ns, loggingOperatorName)
		err = c.CollectTestCoverageFiles(ns, loggingOperatorName)
		if err != nil {
			t.Logf("Failed collecting coverage files: %s", err)
		}
		return err

	}, func(o *cluster.Options) {
		if o.Scheme == nil {
			o.Scheme = runtime.NewScheme()
		}
		common.RequireNoError(t, v1beta1.AddToScheme(o.Scheme))
		common.RequireNoError(t, apiextensionsv1.AddToScheme(o.Scheme))
		common.RequireNoError(t, appsv1.AddToScheme(o.Scheme))
		common.RequireNoError(t, batchv1.AddToScheme(o.Scheme))
		common.RequireNoError(t, corev1.AddToScheme(o.Scheme))
		common.RequireNoError(t, rbacv1.AddToScheme(o.Scheme))
	})
}

func installFluentdSts(c common.Cluster) error {
	manager := helm.New(c.KubeConfigFilePath())

	if err := manager.RunRepo(helm.WithArgs("add", "fluent", "https://fluent.github.io/helm-charts")); err != nil {
		return fmt.Errorf("failed to add fluent repo: %v", err)
	}

	if err := manager.RunInstall(
		helm.WithName("fluentd"),
		helm.WithChart("fluent/fluentd"),
		helm.WithArgs("--create-namespace"),
		helm.WithNamespace("fluentd"),
		helm.WithArgs("--set", "kind=StatefulSet"),
		helm.WithWait(),
	); err != nil {
		return fmt.Errorf("failed to install fluentd: %v", err)
	}

	return nil
}
