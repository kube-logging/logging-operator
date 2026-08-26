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

package wait

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

// DeploymentAvailable returns a condition function that checks if a deployment
// is available with all replicas ready.
func DeploymentAvailable(t *testing.T, cl client.Client, ctx context.Context, namespace, name string) func() bool {
	return func() bool {
		deployment := &appsv1.Deployment{}
		if err := cl.Get(ctx, client.ObjectKey{
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
