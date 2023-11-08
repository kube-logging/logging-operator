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
