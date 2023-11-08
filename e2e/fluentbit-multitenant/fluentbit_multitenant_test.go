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

package fluentd_aggregator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestFluentbitSingleTenantPlusInfra(t *testing.T) {
	common.Initialize(t)
	nsInfra := "infra"
	nsTenant := "tenant"
	tagInfra := "tag_infra"
	tagTenant := "tag_tenant"

	release := "fluentbit-multitenant"
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

		loggingInfra(ctx, t, c.GetClient(), nsInfra, release, tagInfra)
		loggingTenant(ctx, t, c.GetClient(), nsTenant, nsInfra, release, tagTenant)

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
			t.Logf("log consumer logs: %s", rawOut)
			return strings.Contains(string(rawOut), tagTenant) && strings.Contains(string(rawOut), tagInfra)
		}, 5*time.Minute, 3*time.Second)

	}, func(t *testing.T, c common.Cluster) error {
		path := filepath.Join(TestTempDir, fmt.Sprintf("cluster-%s.log", t.Name()))
		t.Logf("Printing cluster logs to %s", path)
		return c.PrintLogs(common.PrintLogConfig{
			Namespaces: []string{nsInfra, nsTenant, "default"},
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

func loggingInfra(ctx context.Context, t *testing.T, c client.Client, nsInfra string, release string, tag string) {
	output := v1beta1.ClusterOutput{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http",
			Namespace: nsInfra,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				LoggingRef: "infra",
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    fmt.Sprintf("http://%s-test-receiver:8080/%s", release, tag),
					ContentType: "application/json",
					Buffer:      realTimeBuffer,
				},
			},
		},
	}

	common.RequireNoError(t, c.Create(ctx, &output))
	flow := v1beta1.ClusterFlow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flow",
			Namespace: nsInfra,
		},
		Spec: v1beta1.ClusterFlowSpec{
			LoggingRef: "infra",
			Match: []v1beta1.ClusterMatch{
				{
					ClusterSelect: &v1beta1.ClusterSelect{
						Labels: producerLabels,
					},
				},
			},
			GlobalOutputRefs: []string{output.Name},
		},
	}
	common.RequireNoError(t, c.Create(ctx, &flow))

	agent := v1beta1.FluentbitAgent{
		ObjectMeta: metav1.ObjectMeta{
			Name: "infra",
		},
		Spec: v1beta1.FluentbitSpec{
			LoggingRef: "infra",
		},
	}
	common.RequireNoError(t, c.Create(ctx, &agent))

	ap := v1beta1.LoggingRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tenants",
		},
		Spec: v1beta1.LoggingRouteSpec{
			Source: "infra",
			Targets: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "tenant",
						Operator: metav1.LabelSelectorOpExists,
					},
				},
			},
		},
	}
	common.RequireNoError(t, c.Create(ctx, &ap))

	logging := v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name: "infra",
			Labels: map[string]string{
				"tenant": "infra",
			},
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       "infra",
			ControlNamespace: nsInfra,
			FluentdSpec: &v1beta1.FluentdSpec{
				DisablePvc: true,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	}
	common.RequireNoError(t, c.Create(ctx, &logging))
}

func loggingTenant(ctx context.Context, t *testing.T, c client.Client, nsTenant, nsInfra, release, tag string) {
	output := v1beta1.Output{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "http",
			Namespace: nsTenant,
		},
		Spec: v1beta1.OutputSpec{
			LoggingRef: "tenant",
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    fmt.Sprintf("http://%s-test-receiver.%s:8080/%s", release, nsInfra, tag),
				ContentType: "application/json",
				Buffer:      realTimeBuffer,
			},
		},
	}

	common.RequireNoError(t, c.Create(ctx, &output))
	flow := v1beta1.Flow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "flow",
			Namespace: nsTenant,
		},
		Spec: v1beta1.FlowSpec{
			LoggingRef: "tenant",
			Match: []v1beta1.Match{
				{
					Select: &v1beta1.Select{
						Labels: producerLabels,
					},
				},
			},
			LocalOutputRefs: []string{output.Name},
		},
	}
	common.RequireNoError(t, c.Create(ctx, &flow))

	logging := v1beta1.Logging{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tenant",
			Labels: map[string]string{
				"tenant": "tenant",
			},
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       "tenant",
			ControlNamespace: nsTenant,
			WatchNamespaces:  []string{"tenant"},
			FluentdSpec: &v1beta1.FluentdSpec{
				DisablePvc: true,
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	}
	common.RequireNoError(t, c.Create(ctx, &logging))
}
