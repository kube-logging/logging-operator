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
	"slices"

	"github.com/cisco-open/operator-tools/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

func podInPhase(name string, phases []corev1.PodPhase, opts ...client.ListOption) Condition {
	return Condition{
		Name: name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var pods corev1.PodList
			if err := cl.List(ctx, &pods, opts...); err != nil {
				return false, err
			}
			for _, pod := range pods.Items {
				if slices.Contains(phases, pod.Status.Phase) {
					return true, nil
				}
			}
			return false, nil
		},
	}
}

// PodRunning is the escape hatch; the constructors below carry their own name
// and selector so a call site passes only what varies.
func PodRunning(name string, opts ...client.ListOption) Condition {
	return podInPhase(name, []corev1.PodPhase{corev1.PodRunning}, opts...)
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

// PodFinished is the counterpart of PodRunning: a config check reports by
// running to completion rather than by staying up.
func PodFinished(name string, opts ...client.ListOption) Condition {
	return podInPhase(name, []corev1.PodPhase{corev1.PodSucceeded, corev1.PodFailed}, opts...)
}

func ConfigCheck(labels map[string]string) Condition {
	return PodFinished("config check", client.MatchingLabels(labels))
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

// LoggingHealthy is a Logging the operator has found nothing wrong with.
func LoggingHealthy(name string) Condition {
	return Condition{
		Name: "Logging " + name + " without problems",
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			problems, err := loggingProblems(ctx, cl, name)
			return err == nil && len(problems) == 0, err
		},
	}
}

// LoggingProblem and LoggingProblemCleared take a matcher rather than a string
// because a failing config check names the checksum in the message.
func LoggingProblem(name string, match func(string) bool) Condition {
	return Condition{
		Name: "Logging " + name + " reporting the problem",
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			problems, err := loggingProblems(ctx, cl, name)
			return err == nil && slices.ContainsFunc(problems, match), err
		},
	}
}

func LoggingProblemCleared(name string, match func(string) bool) Condition {
	return Condition{
		Name: "Logging " + name + " clearing the problem",
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			problems, err := loggingProblems(ctx, cl, name)
			return err == nil && !slices.ContainsFunc(problems, match), err
		},
	}
}

func loggingProblems(ctx context.Context, cl client.Reader, name string) ([]string, error) {
	var logging v1beta1.Logging
	if err := cl.Get(ctx, client.ObjectKey{Name: name}, &logging); err != nil {
		return nil, err
	}
	return logging.Status.Problems, nil
}

// JobStarted is the drainer the operator creates when an aggregator replica
// goes away. Presence is not enough, since the Job exists before it runs, but
// Active alone would be a race: a Job that finishes between two polls is never
// seen active. Counting the terminal states too makes any observation after it
// started sufficient.
func JobStarted(namespace, name string) Condition {
	return Condition{
		Name: "drainer Job " + name + " started",
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var job batchv1.Job
			if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &job); err != nil {
				return false, err
			}
			return job.Status.Active+job.Status.Succeeded+job.Status.Failed > 0, nil
		},
	}
}

// gone reports the object having been collected. NotFound is the answer rather
// than an error, and any other failure leaves the caller polling.
func gone(kind, namespace, name string, obj client.Object) Condition {
	return Condition{
		Name: kind + " " + name + " gone",
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, obj)
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		},
	}
}

func JobGone(namespace, name string) Condition {
	return gone("Job", namespace, name, &batchv1.Job{})
}

func PodGone(namespace, name string) Condition {
	return gone("Pod", namespace, name, &corev1.Pod{})
}

func PVCGone(namespace, name string) Condition {
	return gone("PersistentVolumeClaim", namespace, name, &corev1.PersistentVolumeClaim{})
}

// Deployment is ready when every replica it wants is ready, available, and the
// Available condition agrees. Ready alone can hold while the condition still
// reports the old rollout.
//
// Wanting no replicas is deliberately never ready: the old helper bailed on a
// nil Spec.Replicas, and taking it as zero would otherwise let a Deployment
// scaled to nothing satisfy a readiness wait with a stale Available condition.
func Deployment(namespace, name string) Condition {
	return Condition{
		Name: "Deployment " + name,
		Met: func(ctx context.Context, cl client.Reader) (bool, error) {
			var deployment appsv1.Deployment
			if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, &deployment); err != nil {
				return false, err
			}
			want := ptr.Deref(deployment.Spec.Replicas, 0)
			if want == 0 || deployment.Status.ReadyReplicas != want || deployment.Status.AvailableReplicas != want {
				return false, nil
			}
			for _, c := range deployment.Status.Conditions {
				if c.Type == appsv1.DeploymentAvailable {
					return c.Status == corev1.ConditionTrue, nil
				}
			}
			return false, nil
		},
	}
}
