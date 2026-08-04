// Copyright © 2026 Kube logging authors
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

package wait

import (
	"context"

	"github.com/cisco-open/operator-tools/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Condition carries no testing.T and no client, so a suite can name one without
// assembling it and whoever evaluates it supplies both.
type Condition struct {
	Name string
	Met  func(context.Context, client.Reader) (bool, error)
}

// PodRunning is the escape hatch; the constructors below carry their own name
// and selector so a call site passes only what varies.
func PodRunning(name string, opts ...client.ListOption) Condition {
	return Condition{
		Name: name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var pods corev1.PodList
			if err := cl.List(ctx, &pods, opts...); err != nil {
				return false, err
			}
			for _, pod := range pods.Items {
				if pod.Status.Phase == corev1.PodRunning {
					return true, nil
				}
			}
			return false, nil
		},
	}
}

func Operator(release string) Condition {
	return PodRunning("operator", client.MatchingLabels{types.NameLabel: release})
}

func Producer(labels map[string]string) Condition {
	return PodRunning("producer", client.MatchingLabels(labels))
}

func FluentdAggregator(namespace string) Condition {
	return aggregator("fluentd", namespace)
}

func SyslogNGAggregator(namespace string) Condition {
	return aggregator("syslog-ng", namespace)
}

func aggregator(kind, namespace string) Condition {
	return PodRunning(kind+" aggregator in "+namespace,
		client.MatchingLabels{types.NameLabel: kind, types.ComponentLabel: kind},
		client.InNamespace(namespace))
}
