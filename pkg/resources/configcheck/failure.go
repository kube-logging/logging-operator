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
	"fmt"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// RetryAbnormalFailure deletes a failed configcheck pod that did not fail on the config itself and
// returns the message to report while it is retried, or an empty message when the failure is a
// genuine verdict on the config.
func RetryAbnormalFailure(ctx context.Context, c client.Client, log logr.Logger, pod *corev1.Pod) (string, error) {
	// A pod-level Reason (DeadlineExceeded, Evicted, Shutdown, NodeAffinity, ...) means the pod was
	// killed around the check; an invalid config only makes the container exit non-zero, which leaves
	// the pod-level Reason empty.
	if pod.Status.Reason == "" {
		return "", nil
	}

	log.Info("configcheck pod did not complete normally, deleting it to retry",
		"pod", pod.Name,
		"reason", pod.Status.Reason,
		"activeDeadlineSeconds", pod.Spec.ActiveDeadlineSeconds,
		"runningContainer", stillRunningContainer(pod))

	if err := client.IgnoreNotFound(c.Delete(ctx, pod)); err != nil {
		return "", errors.WrapIf(err, "failed to delete configcheck pod that did not complete normally")
	}

	return fmt.Sprintf("configcheck pod %s did not complete normally (reason: %s), deleted for retry",
		pod.Name, pod.Status.Reason), nil
}

// stillRunningContainer names the container that was still running when the pod stopped, so a
// helper that held up the check can be identified.
func stillRunningContainer(pod *corev1.Pod) string {
	for _, cs := range pod.Status.InitContainerStatuses {
		if cs.State.Running != nil {
			return cs.Name
		}
	}
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Running != nil {
			return cs.Name
		}
	}
	return ""
}
