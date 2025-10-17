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

package cond

import (
	"context"
	"testing"

	"github.com/cisco-open/operator-tools/pkg/utils"
	"github.com/kube-logging/logging-operator/e2e/common"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PodShouldBeRunning(t *testing.T, cl client.Reader, key client.ObjectKey) func() bool {
	return func() bool {
		var pod corev1.Pod
		err := cl.Get(context.Background(), key, &pod)

		if pod.Status.Phase == corev1.PodRunning {
			return true
		}

		if err == nil {
			t.Logf("pod %s is in phase %s", key, pod.Status.Phase)
		} else if !apierrors.IsNotFound(err) {
			t.Logf("an error occurred while getting pod %s: %v", key, err)
		}

		return false
	}
}

func AnyPodShouldBeRunning(t *testing.T, cl client.Reader, opts ...client.ListOption) func() bool {
	return func() bool {
		var podList corev1.PodList
		if err := cl.List(context.Background(), &podList, opts...); err != nil {
			t.Logf("an error occurred while listing pods: %v", err)
		}
		for _, pod := range podList.Items {
			if pod.Status.Phase == corev1.PodRunning {
				return true
			}
		}
		return false
	}
}

func AnyPodShouldBeFinished(t *testing.T, cl client.Reader, opts ...client.ListOption) func() bool {
	return func() bool {
		var podList corev1.PodList
		if err := cl.List(context.Background(), &podList, opts...); err != nil {
			t.Logf("an error occurred while listing pods: %v", err)
		}
		for _, pod := range podList.Items {
			if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
				return true
			}
		}
		return false
	}
}

func ResourceShouldBeAbsent(t *testing.T, cl client.Reader, obj client.Object) func() bool {
	return func() bool {
		err := cl.Get(context.Background(), client.ObjectKeyFromObject(obj), obj)
		if apierrors.IsNotFound(err) {
			return true
		}
		if err != nil {
			t.Logf("an error occurred while getting %q resource: %v", obj.GetObjectKind().GroupVersionKind(), err)
		}
		return false
	}
}

func ResourceShouldBePresent(t *testing.T, cl client.Reader, obj client.Object) func() bool {
	return func() bool {
		err := cl.Get(context.Background(), client.ObjectKeyFromObject(obj), obj)
		if err == nil {
			return true
		}
		if !apierrors.IsNotFound(err) {
			t.Logf("an error occurred while getting %q resource: %v", obj.GetObjectKind().GroupVersionKind(), err)
		}
		return false
	}
}
func CheckFluentdStatus(t *testing.T, c *common.Cluster, ctx *context.Context, fluentd *v1beta1.FluentdConfig, loggingName string) bool {
	fluentdInstanceName := fluentd.Name
	cluster := *c

	if len(fluentd.Status.Problems) != 0 {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(fluentd), fluentd))
		t.Logf("%s should have 0 problems, problems=%v", fluentdInstanceName, fluentd.Status.Problems)
		return false
	}
	if fluentd.Status.Logging != loggingName {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(fluentd), fluentd))
		t.Logf("%s should have it's logging field filled, found: %s, expect:%s", fluentdInstanceName, fluentd.Status.Logging, loggingName)
		return false
	}
	if !*fluentd.Status.Active {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(fluentd), fluentd))
		t.Logf("%s should have it's active field set as true, found: %v", fluentdInstanceName, *fluentd.Status.Active)
		return false
	}

	return true
}

func CheckExcessFluentdStatus(t *testing.T, c *common.Cluster, ctx *context.Context, fluentd *v1beta1.FluentdConfig) bool {
	fluentdInstanceName := fluentd.Name
	cluster := *c

	if len(fluentd.Status.Problems) == 0 {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(fluentd), fluentd))
		t.Logf("%s should have it's problems field filled", fluentdInstanceName)
		return false
	}
	if fluentd.Status.Logging != "" {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(fluentd), fluentd))
		t.Logf("%s should have it's logging field empty, found: %s", fluentdInstanceName, fluentd.Status.Logging)
		return false
	}
	if *fluentd.Status.Active {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(fluentd), fluentd))
		t.Logf("%s should have it's active field set as false, found: %v", fluentdInstanceName, *fluentd.Status.Active)
		return false
	}

	return true
}

func CheckSyslogNGStatus(t *testing.T, c *common.Cluster, ctx *context.Context, syslogNG *v1beta1.SyslogNGConfig, loggingName string) bool {
	instanceName := syslogNG.Name
	cluster := *c

	if len(syslogNG.Status.Problems) != 0 {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(syslogNG), syslogNG))
		t.Logf("%s should have 0 problems, problems=%v", instanceName, syslogNG.Status.Problems)
		return false
	}
	if syslogNG.Status.Logging != loggingName {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(syslogNG), syslogNG))
		t.Logf("%s should have it's logging field filled, found: %s, expect:%s", instanceName, syslogNG.Status.Logging, loggingName)
		return false
	}
	if !*syslogNG.Status.Active {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(syslogNG), syslogNG))
		t.Logf("%s should have it's active field set as true, found: %v", instanceName, *syslogNG.Status.Active)
		return false
	}

	return true
}

func CheckExcessSyslogNGStatus(t *testing.T, c *common.Cluster, ctx *context.Context, syslogNG *v1beta1.SyslogNGConfig) bool {
	instanceName := syslogNG.Name
	cluster := *c

	if len(syslogNG.Status.Problems) == 0 {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(syslogNG), syslogNG))
		t.Logf("%s should have it's problems field filled", instanceName)
		return false
	}
	if syslogNG.Status.Logging != "" {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(syslogNG), syslogNG))
		t.Logf("%s should have it's logging field empty, found: %s", instanceName, syslogNG.Status.Logging)
		return false
	}
	if *syslogNG.Status.Active {
		common.RequireNoError(t, cluster.GetClient().Get(*ctx, utils.ObjectKeyFromObjectMeta(syslogNG), syslogNG))
		t.Logf("%s should have it's active field set as false, found: %v", instanceName, *syslogNG.Status.Active)
		return false
	}

	return true
}

// DeploymentAvailable returns a condition function that checks if a deployment
// is available with all replicas ready.
func DeploymentAvailable(t *testing.T, c client.Client, ctx *context.Context, namespace, name string) func() bool {
	return func() bool {
		deployment := &appsv1.Deployment{}
		if err := c.Get(*ctx, client.ObjectKey{
			Name:      name,
			Namespace: namespace,
		}, deployment); err != nil {
			t.Logf("Failed to get deployment %s/%s: %v", namespace, name, err)
			return false
		}

		if deployment.Spec.Replicas == nil {
			return false
		}
		desiredReplicas := *deployment.Spec.Replicas

		if deployment.Status.ReadyReplicas != desiredReplicas {
			t.Logf("Deployment %s/%s: %d/%d replicas ready",
				namespace, name, deployment.Status.ReadyReplicas, desiredReplicas)
			return false
		}

		if deployment.Status.AvailableReplicas != desiredReplicas {
			t.Logf("Deployment %s/%s: %d/%d replicas available",
				namespace, name, deployment.Status.AvailableReplicas, desiredReplicas)
			return false
		}

		for _, condition := range deployment.Status.Conditions {
			if condition.Type == appsv1.DeploymentAvailable {
				if condition.Status == corev1.ConditionTrue {
					t.Logf("Deployment %s/%s is available", namespace, name)
					return true
				}
				t.Logf("Deployment %s/%s Available condition is %s: %s",
					namespace, name, condition.Status, condition.Message)
				return false
			}
		}

		t.Logf("Deployment %s/%s has no Available condition", namespace, name)
		return false
	}
}
