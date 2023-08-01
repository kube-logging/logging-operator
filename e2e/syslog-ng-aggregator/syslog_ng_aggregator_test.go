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

package syslong_ng_aggregator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cisco-open/operator-tools/pkg/typeoverride"
	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"

	"github.com/kube-logging/logging-operator/pkg/resources/syslogng"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/filter"
	syslogngoutput "github.com/kube-logging/logging-operator/pkg/sdk/logging/model/syslogng/output"

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

func TestSyslogNGIsRunningAndForwardingLogs(t *testing.T) {
	common.Initialize(t)
	ns := "test"
	releaseNameOverride := "e2e"
	common.WithCluster("syslog-ng-1", t, func(t *testing.T, c common.Cluster) {
		setup.LoggingOperator(t, c, setup.LoggingOperatorOptionFunc(func(options *setup.LoggingOperatorOptions) {
			options.Namespace = ns
			options.NameOverride = releaseNameOverride
		}))

		ctx := context.Background()

		logging := v1beta1.Logging{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "syslog-ng-aggregator-test",
				Namespace: ns,
			},
			Spec: v1beta1.LoggingSpec{
				EnableRecreateWorkloadOnImmutableFieldChange: true,
				ControlNamespace: ns,
				FluentbitSpec: &v1beta1.FluentbitSpec{
					Network: &v1beta1.FluentbitNetwork{
						Keepalive: utils.BoolPointer(false),
					},
				},
				SyslogNGSpec: &v1beta1.SyslogNGSpec{
					StatefulSetOverrides: &typeoverride.StatefulSet{
						Spec: typeoverride.StatefulSetSpec{
							Template: typeoverride.PodTemplateSpec{
								Spec: typeoverride.PodSpec{
									Containers: []corev1.Container{
										{
											Name: syslogng.ContainerName,
											Resources: corev1.ResourceRequirements{
												Limits: corev1.ResourceList{
													corev1.ResourceCPU:    resource.MustParse("100m"),
													corev1.ResourceMemory: resource.MustParse("100M"),
												},
												Requests: corev1.ResourceList{
													corev1.ResourceCPU:    resource.MustParse("25m"),
													corev1.ResourceMemory: resource.MustParse("10M"),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
		testTag := "test.tag"

		common.RequireNoError(t, c.GetClient().Create(ctx, &logging))
		output := v1beta1.SyslogNGOutput{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-output",
				Namespace: ns,
			},
			Spec: v1beta1.SyslogNGOutputSpec{
				HTTP: &syslogngoutput.HTTPOutput{
					URL: fmt.Sprintf("http://%s-test-receiver:8080/%s", releaseNameOverride, testTag),
					Headers: []string{
						"Content-type: application/json",
					},
					Method: "POST",
					DiskBuffer: &syslogngoutput.DiskBuffer{
						DiskBufSize: 100 * 1024 * 1024,
						Reliable:    true,
						Dir:         syslogng.BufferPath,
					},
				},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &output))
		flow := v1beta1.SyslogNGFlow{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-flow",
				Namespace: ns,
			},
			Spec: v1beta1.SyslogNGFlowSpec{
				Match: &v1beta1.SyslogNGMatch{
					Regexp: &filter.RegexpMatchExpr{
						Pattern: "log-producer",
						Value:   "json.kubernetes.labels.my-unique-label",
						Type:    "string",
					},
				},
				LocalOutputRefs: []string{output.Name},
			},
		}
		common.RequireNoError(t, c.GetClient().Create(ctx, &flow))

		aggergatorLabels := map[string]string{
			"app.kubernetes.io/name":      "syslog-ng",
			"app.kubernetes.io/component": "syslog-ng",
		}
		operatorLabels := map[string]string{
			"app.kubernetes.io/name": releaseNameOverride,
		}
		producerLabels := map[string]string{
			"my-unique-label": "log-producer",
		}
		go setup.LogProducer(t, c.GetClient(), setup.LogProducerOptionFunc(func(options *setup.LogProducerOptions) {
			options.Namespace = ns
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
			if aggregatorRunning := cond.AnyPodShouldBeRunning(t, c.GetClient(), client.MatchingLabels(aggergatorLabels)); !aggregatorRunning() {
				t.Log("waiting for the aggregator")
				return false
			}

			cmd := common.CmdEnv(exec.Command("kubectl",
				"logs",
				"-n", ns,
				"-l", fmt.Sprintf("app.kubernetes.io/name=%s-test-receiver", releaseNameOverride,
				)), c)
			rawOut, err := cmd.Output()
			if err != nil {
				t.Logf("failed to get log consumer logs: %+v %s", err, rawOut)
				return false
			}
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), testTag)
		}, 5*time.Minute, 2*time.Second)

	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		return c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{ns, "default"},
			FilePath:   path,
			Limit:      100 * 1000,
		})
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
