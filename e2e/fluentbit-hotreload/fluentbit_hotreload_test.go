// Copyright © 2024 Kube logging authors
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

package fluentbit_hotreload

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"

	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/e2e/common/cond"
	"github.com/kube-logging/logging-operator/e2e/common/setup"
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

var tags = "time"
var realTimeBuffer = &output.Buffer{
	Tags:        &tags,
	Timekey:     "1s",
	TimekeyWait: "0s",
}
var producerLabels = map[string]string{
	"my-unique-label": "log-producer",
}

func TestFluentbitHotReload(t *testing.T) {
	common.Initialize(t)
	nsInfra := "infra"
	nsTenant := "tenant"
	tagInfra := "tag_infra"
	tagTenant := "tag_tenant"

	release := "fluentbit-hotreload"
	common.WithCluster(release, t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = nsInfra
			options.NameOverride = release
		}))

		ctx := context.Background()

		common.RequireNoError(t, c.GetClient().Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: nsTenant,
			},
		}))

		common.LoggingInfra(ctx, t, c.GetClient(), nsInfra, release, tagInfra, realTimeBuffer, producerLabels, &v1beta1.HotReload{})
		common.LoggingTenant(ctx, t, c.GetClient(), nsTenant, nsInfra, release, tagTenant, realTimeBuffer, producerLabels)

		aggregatorLabels := map[string]string{
			"app.kubernetes.io/name":      "fluentd",
			"app.kubernetes.io/component": "fluentd",
		}
		operatorLabels := map[string]string{
			"app.kubernetes.io/name": release,
		}

		// start log producer in the tenant namespace
		go setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = nsTenant
			options.Labels = producerLabels
		}))

		require.Eventually(t, func() bool {
			if operatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(operatorLabels))(); !operatorRunning {
				t.Log("waiting for the operator")
				return false
			}
			if producerRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(producerLabels))(); !producerRunning {
				t.Log("waiting for the producer")
				return false
			}
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggregatorLabels), client.InNamespace(nsInfra)); !aggregatorRunning() {
				t.Log("waiting for the infra aggregator")
				return false
			}
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggregatorLabels), client.InNamespace(nsTenant)); !aggregatorRunning() {
				t.Log("waiting for the tenant aggregator")
				return false
			}
			return true
		}, 5*time.Minute, 3*time.Second)

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl",
				"logs",
				"-n", nsInfra,
				"--tail", "100",
				"-l", fmt.Sprintf("app.kubernetes.io/name=%s-test-receiver", release)), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs should contain no tenant, only infra logs: %s", rawOut)
			return !strings.Contains(string(rawOut), tagTenant) && strings.Contains(string(rawOut), tagInfra)
		}, 5*time.Minute, 3*time.Second)

		common.LoggingRoute(ctx, t, c.GetClient())

		require.Eventually(t, func() bool {
			cmd := common.CmdEnv(exec.Command("kubectl",
				"logs",
				"-n", nsInfra,
				"--tail", "30",
				"-l", fmt.Sprintf("app.kubernetes.io/name=%s-test-receiver", release)), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %v", err)
				return false
			}
			t.Logf("log consumer logs should contain tenant logs: %s", rawOut)
			return strings.Contains(string(rawOut), tagTenant)
		}, 5*time.Minute, 3*time.Second)

		ds := &appsv1.DaemonSet{}
		err := c.GetClient().Get(ctx, types.NamespacedName{
			Namespace: nsInfra,
			Name:      "infra-fluentbit",
		}, ds)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), ds.Generation, "generation should not be incremented for a reloadable agent")

	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		err := c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{nsInfra, nsTenant, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
		if err != nil {
			return err
		}

		loggingOperatorName := "logging-operator-" + release
		t.Logf("Collecting coverage files from logging-operator: %s/%s", nsInfra, loggingOperatorName)
		err = c.CollectTestCoverageFiles(nsInfra, loggingOperatorName)
		if err != nil {
			t.Logf("Failed collecting coverage files: %s", err)
		}

		return nil
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
