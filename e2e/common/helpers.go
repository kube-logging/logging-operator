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
	"testing"
	"time"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func RequireNoError(t *testing.T, err error) {
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Received unexpected error:\n%#v %+v", err, errors.GetDetails(err)))
		t.FailNow()
	}
}

func Initialize(t *testing.T) {
	t.Parallel()
}

// waitForPodReady is satisfied by the Ready condition, or by Running when the
// pod does not carry one yet. SetupCurlPod is the only caller left; suites wait
// through internal/wait.
func waitForPodReady(ctx context.Context, c client.Client, pod *corev1.Pod, pollInterval, pollTimeout time.Duration) error {
	return wait.PollUntilContextTimeout(ctx, pollInterval, pollTimeout, true, wait.ConditionWithContextFunc(func(ctx context.Context) (bool, error) {
		var updatedPod corev1.Pod
		err := c.Get(ctx, client.ObjectKeyFromObject(pod), &updatedPod)
		if client.IgnoreNotFound(err) != nil {
			return false, fmt.Errorf("failed to get pod status: %w", err)
		}

		isReady := updatedPod.Status.Phase == corev1.PodRunning
		for _, cond := range updatedPod.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return isReady, nil
	}))
}

// SetupCurlPod creates a curl pod for testing HTTP endpoints and waits for it to be ready
func SetupCurlPod(ctx context.Context, c client.Client, namespace, name string, pollInterval, pollTimeout time.Duration) (*corev1.Pod, error) {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "curl",
					Image:   "curlimages/curl:latest",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}

	if err := c.Create(ctx, pod); err != nil {
		return nil, fmt.Errorf("failed to create curl pod: %w", err)
	}

	if err := waitForPodReady(ctx, c, pod, pollInterval, pollTimeout); err != nil {
		return nil, fmt.Errorf("failed to wait for curl pod to be ready: %w", err)
	}

	return pod, nil
}
