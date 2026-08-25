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
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
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

func Pod(namespace, name string) Condition {
	key := client.ObjectKey{Namespace: namespace, Name: name}
	return Condition{
		Name: "pod " + namespace + "/" + name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var pod corev1.Pod
			if err := cl.Get(ctx, key, &pod); err != nil {
				return false, client.IgnoreNotFound(err)
			}
			return pod.Status.Phase == corev1.PodRunning, nil
		},
	}
}

func Operator(release string) Condition {
	return PodRunning("operator", client.MatchingLabels{types.NameLabel: release})
}

func Producer(labels map[string]string) Condition {
	return PodRunning("producer", client.MatchingLabels(labels))
}

// Fluentbit is the agent daemonset, which a suite naming a
// fluentBitAgentNamespace expects somewhere other than the control namespace.
func Fluentbit(namespace string) Condition {
	return PodRunning("fluentbit in "+namespace,
		client.MatchingLabels{types.NameLabel: "fluentbit"},
		client.InNamespace(namespace))
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

// A detached config carries the verdict the Logging wrote back onto it:
// attached means it owns the aggregator, excess means another already does.
func attached(problems []string, logging string, active *bool, loggingName string) bool {
	return len(problems) == 0 && logging == loggingName && ptr.Deref(active, false)
}

func excess(problems []string, logging string, active *bool) bool {
	return len(problems) > 0 && logging == "" && !ptr.Deref(active, true)
}

func AttachedFluentd(namespace, name, loggingName string) Condition {
	return Condition{
		Name: "attached FluentdConfig " + name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var config v1beta1.FluentdConfig
			if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &config); err != nil {
				return false, err
			}
			return attached(config.Status.Problems, config.Status.Logging, config.Status.Active, loggingName), nil
		},
	}
}

func ExcessFluentd(namespace, name string) Condition {
	return Condition{
		Name: "excess FluentdConfig " + name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var config v1beta1.FluentdConfig
			if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &config); err != nil {
				return false, err
			}
			return excess(config.Status.Problems, config.Status.Logging, config.Status.Active), nil
		},
	}
}

func AttachedSyslogNG(namespace, name, loggingName string) Condition {
	return Condition{
		Name: "attached SyslogNGConfig " + name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var config v1beta1.SyslogNGConfig
			if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &config); err != nil {
				return false, err
			}
			return attached(config.Status.Problems, config.Status.Logging, config.Status.Active, loggingName), nil
		},
	}
}

func ExcessSyslogNG(namespace, name string) Condition {
	return Condition{
		Name: "excess SyslogNGConfig " + name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var config v1beta1.SyslogNGConfig
			if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &config); err != nil {
				return false, err
			}
			return excess(config.Status.Problems, config.Status.Logging, config.Status.Active), nil
		},
	}
}

// LoggingUsesFluentd is the other direction: the Logging naming the detached
// config, which it writes before the config's own status settles.
func LoggingUsesFluentd(name, configName string) Condition {
	return Condition{
		Name: "Logging " + name + " using FluentdConfig " + configName,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var logging v1beta1.Logging
			if err := cl.Get(ctx, client.ObjectKey{Name: name}, &logging); err != nil {
				return false, err
			}
			return logging.Status.FluentdConfigName == configName, nil
		},
	}
}

func LoggingUsesSyslogNG(name, configName string) Condition {
	return Condition{
		Name: "Logging " + name + " using SyslogNGConfig " + configName,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var logging v1beta1.Logging
			if err := cl.Get(ctx, client.ObjectKey{Name: name}, &logging); err != nil {
				return false, err
			}
			return logging.Status.SyslogNGConfigName == configName, nil
		},
	}
}
