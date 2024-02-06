// Copyright © 2023 Kube logging authors
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

package common

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"emperror.dev/errors"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	v12 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/output"
)

var sequence uint32

func RequireNoError(t *testing.T, err error) {
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Received unexpected error:\n%#v %+v", err, errors.GetDetails(err)))
		t.FailNow()
	}
}

func Initialize(t *testing.T) {
	localSeq := atomic.AddUint32(&sequence, 1)
	shards := cast.ToUint32(os.Getenv("SHARDS"))
	shard := cast.ToUint32(os.Getenv("SHARD"))
	if shards > 0 {
		if localSeq%shards != shard {
			t.Skipf("skipping %s as sequence %d not in shard %d", t.Name(), localSeq, shard)
		}
	}
	t.Parallel()
}

func LoggingInfra(
	ctx context.Context,
	t *testing.T,
	c client.Client,
	nsInfra string,
	release string,
	tag string,
	buffer *output.Buffer,
	producerLabels map[string]string,
	hotReload *v1beta1.HotReload) {

	output := v1beta1.ClusterOutput{
		ObjectMeta: v1.ObjectMeta{
			Name:      "http",
			Namespace: nsInfra,
		},
		Spec: v1beta1.ClusterOutputSpec{
			OutputSpec: v1beta1.OutputSpec{
				LoggingRef: "infra",
				HTTPOutput: &output.HTTPOutputConfig{
					Endpoint:    fmt.Sprintf("http://%s-test-receiver:8080/%s", release, tag),
					ContentType: "application/json",
					Buffer:      buffer,
				},
			},
		},
	}

	RequireNoError(t, c.Create(ctx, &output))
	flow := v1beta1.ClusterFlow{
		ObjectMeta: v1.ObjectMeta{
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
	RequireNoError(t, c.Create(ctx, &flow))

	agent := v1beta1.FluentbitAgent{
		ObjectMeta: v1.ObjectMeta{
			Name: "infra",
		},
		Spec: v1beta1.FluentbitSpec{
			LoggingRef:      "infra",
			ConfigHotReload: hotReload,
		},
	}
	RequireNoError(t, c.Create(ctx, &agent))

	logging := v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
			Name: "infra",
			Labels: map[string]string{
				"tenant": "infra",
			},
		},
		Spec: v1beta1.LoggingSpec{
			LoggingRef:       "infra",
			ControlNamespace: nsInfra,
			FluentdSpec: &v1beta1.FluentdSpec{
				Image: v1beta1.ImageSpec{
					Tag: "v1.16-base",
				},
				DisablePvc: true,
				Resources: v12.ResourceRequirements{
					Requests: v12.ResourceList{
						v12.ResourceCPU:    resource.MustParse("50m"),
						v12.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &logging))
}

func LoggingTenant(
	ctx context.Context,
	t *testing.T,
	c client.Client,
	nsTenant,
	nsInfra,
	release,
	tag string,
	buffer *output.Buffer,
	producerLabels map[string]string) {
	output := v1beta1.Output{
		ObjectMeta: v1.ObjectMeta{
			Name:      "http",
			Namespace: nsTenant,
		},
		Spec: v1beta1.OutputSpec{
			LoggingRef: "tenant",
			HTTPOutput: &output.HTTPOutputConfig{
				Endpoint:    fmt.Sprintf("http://%s-test-receiver.%s:8080/%s", release, nsInfra, tag),
				ContentType: "application/json",
				Buffer:      buffer,
			},
		},
	}

	RequireNoError(t, c.Create(ctx, &output))
	flow := v1beta1.Flow{
		ObjectMeta: v1.ObjectMeta{
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
	RequireNoError(t, c.Create(ctx, &flow))

	logging := v1beta1.Logging{
		ObjectMeta: v1.ObjectMeta{
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
				Resources: v12.ResourceRequirements{
					Requests: v12.ResourceList{
						v12.ResourceCPU:    resource.MustParse("50m"),
						v12.ResourceMemory: resource.MustParse("50M"),
					},
				},
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &logging))
}

func LoggingRoute(ctx context.Context, t *testing.T, c client.Client) {
	ap := v1beta1.LoggingRoute{
		ObjectMeta: v1.ObjectMeta{
			Name: "tenants",
		},
		Spec: v1beta1.LoggingRouteSpec{
			Source: "infra",
			Targets: v1.LabelSelector{
				MatchExpressions: []v1.LabelSelectorRequirement{
					{
						Key:      "tenant",
						Operator: v1.LabelSelectorOpExists,
					},
				},
			},
		},
	}
	RequireNoError(t, c.Create(ctx, &ap))
}
