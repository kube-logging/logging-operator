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

package configcheck

import (
	"context"
	"testing"

	"emperror.dev/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func failedPod(reason string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "fluentd-configcheck-deadbeef", Namespace: "logging"},
		Status: corev1.PodStatus{
			Phase:  corev1.PodFailed,
			Reason: reason,
		},
	}
}

func newClient(t *testing.T, pod *corev1.Pod, opts ...interceptor.Funcs) client.Client {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(scheme))

	builder := fake.NewClientBuilder().WithScheme(scheme).WithObjects(pod)
	if len(opts) > 0 {
		builder = builder.WithInterceptorFuncs(opts[0])
	}
	return builder.Build()
}

// TestRetryAbnormalFailure pins which failed pods are a verdict on the config. Only a container
// exiting non-zero is, and that leaves the pod-level Reason empty; a pod-level Reason means the pod
// was killed around the check, so it has to be cleared out and tried again.
func TestRetryAbnormalFailure(t *testing.T) {
	tests := []struct {
		name         string
		reason       string
		expectRetry  bool
		expectDelete bool
	}{
		{name: "InvalidConfigIsAVerdict", reason: ""},
		{name: "DeadlineExceeded", reason: "DeadlineExceeded", expectRetry: true, expectDelete: true},
		{name: "Evicted", reason: "Evicted", expectRetry: true, expectDelete: true},
		{name: "NodeAffinity", reason: "NodeAffinity", expectRetry: true, expectDelete: true},
		{name: "UnexpectedAdmissionError", reason: "UnexpectedAdmissionError", expectRetry: true, expectDelete: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pod := failedPod(tt.reason)
			c := newClient(t, pod)
			ctx := context.Background()

			message, err := RetryAbnormalFailure(ctx, c, log.Log, pod)
			require.NoError(t, err)

			if !tt.expectRetry {
				assert.Empty(t, message, "a config verdict must be left to the caller to record")
			} else {
				assert.Contains(t, message, tt.reason)
				assert.Contains(t, message, pod.Name)
			}

			err = c.Get(ctx, client.ObjectKeyFromObject(pod), &corev1.Pod{})
			assert.Equal(t, tt.expectDelete, apierrors.IsNotFound(err),
				"deleting the pod is what lets the next reconcile create a fresh one")
		})
	}
}

// TestRetryAbnormalFailureAlreadyDeleted covers the pod being reaped between the read and the
// delete: there is nothing left to clean up, so the retry still stands.
func TestRetryAbnormalFailureAlreadyDeleted(t *testing.T) {
	pod := failedPod("DeadlineExceeded")
	c := newClient(t, pod)
	ctx := context.Background()
	require.NoError(t, c.Delete(ctx, pod))

	message, err := RetryAbnormalFailure(ctx, c, log.Log, pod)
	require.NoError(t, err)
	assert.Contains(t, message, "DeadlineExceeded")
}

// TestRetryAbnormalFailureDeleteFails pins that a failed delete is reported instead of being
// reported as a retry that never happened.
func TestRetryAbnormalFailureDeleteFails(t *testing.T) {
	pod := failedPod("DeadlineExceeded")
	c := newClient(t, pod, interceptor.Funcs{
		Delete: func(context.Context, client.WithWatch, client.Object, ...client.DeleteOption) error {
			return errors.New("boom")
		},
	})

	message, err := RetryAbnormalFailure(context.Background(), c, log.Log, pod)
	require.Error(t, err)
	assert.Empty(t, message)
}

func TestStillRunningContainer(t *testing.T) {
	running := corev1.ContainerState{Running: &corev1.ContainerStateRunning{}}
	terminated := corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{}}

	tests := []struct {
		name     string
		pod      *corev1.Pod
		expected string
	}{
		{
			name:     "NoStatuses",
			pod:      &corev1.Pod{},
			expected: "",
		},
		{
			name: "EverythingTerminated",
			pod: &corev1.Pod{Status: corev1.PodStatus{
				InitContainerStatuses: []corev1.ContainerStatus{{Name: "tmp-dir-hack", State: terminated}},
				ContainerStatuses:     []corev1.ContainerStatus{{Name: "fluentd", State: terminated}},
			}},
			expected: "",
		},
		{
			// The case the log line exists for: a native sidecar is an init container, so a
			// helper that never exits shows up there while fluentd has already finished.
			name: "NativeSidecarStillRunning",
			pod: &corev1.Pod{Status: corev1.PodStatus{
				InitContainerStatuses: []corev1.ContainerStatus{
					{Name: "tmp-dir-hack", State: terminated},
					{Name: "geoip-refresh", State: running},
				},
				ContainerStatuses: []corev1.ContainerStatus{{Name: "fluentd", State: terminated}},
			}},
			expected: "geoip-refresh",
		},
		{
			name: "DryRunItselfStillRunning",
			pod: &corev1.Pod{Status: corev1.PodStatus{
				InitContainerStatuses: []corev1.ContainerStatus{{Name: "tmp-dir-hack", State: terminated}},
				ContainerStatuses:     []corev1.ContainerStatus{{Name: "fluentd", State: running}},
			}},
			expected: "fluentd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, stillRunningContainer(tt.pod))
		})
	}
}
