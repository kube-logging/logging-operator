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

package harness

import (
	"fmt"
	"strings"

	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	receiverPort = 8080

	// Generous on purpose: more lines makes a tag easier to find.
	receiverLogTail = 100

	wholeLog = -1
)

// Receiver is the test receiver the chart installs, where a suite looks to see
// which logs arrived.
type Receiver struct {
	env *Env
}

func ReceiverName(release string) string { return release + "-test-receiver" }

// ReceiverURL and ReceiverURLIn take a release rather than an Env because
// fixture builds its Outputs before there is one.
func ReceiverURL(release, tag string) string {
	return fmt.Sprintf("http://%s:%d/%s", ReceiverName(release), receiverPort, tag)
}

func ReceiverURLIn(release, namespace, tag string) string {
	return fmt.Sprintf("http://%s.%s:%d/%s", ReceiverName(release), namespace, receiverPort, tag)
}

func (r Receiver) URL(tag string) string { return ReceiverURL(r.env.Release, tag) }

// MustReceive does not echo the tail each poll: the archived cluster dump
// already carries the receiver's log.
func (r Receiver) MustReceive(tags ...string) {
	r.env.T.Helper()

	var outstanding pending
	require.Eventuallyf(r.env.T, func() bool {
		logs, err := r.Logs()
		if err != nil {
			r.env.T.Logf("reading the test receiver: %v", err)
			return false
		}
		for _, tag := range tags {
			if !strings.Contains(logs, tag) {
				outstanding.set(tag)
				return false
			}
		}
		return true
	}, r.env.waitBudget(), waitInterval, "the test receiver never logged %s", &outstanding)
}

// MustNotReceive is a point-in-time check, since an absence cannot be waited
// for. It belongs after whatever wait establishes that the pipeline is running,
// and it reads the whole log rather than the tail MustReceive polls: a tag
// that arrived early would otherwise have scrolled out of view.
func (r Receiver) MustNotReceive(tags ...string) {
	r.env.T.Helper()

	logs, err := r.logs(wholeLog)
	require.NoError(r.env.T, err)
	for _, tag := range tags {
		assert.NotContains(r.env.T, logs, tag)
	}
}

// Scale takes the receiver away and brings it back, which is how a drain test
// makes the aggregator buffer instead of deliver.
func (r Receiver) Scale(replicas int32) {
	r.env.T.Helper()
	name := ReceiverName(r.env.Release)
	_, err := r.env.cluster.clientset.AppsV1().Deployments(r.env.ControlNamespace).UpdateScale(r.env.Ctx, name,
		&autoscalingv1.Scale{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: r.env.ControlNamespace},
			Spec:       autoscalingv1.ScaleSpec{Replicas: replicas},
		}, metav1.UpdateOptions{})
	require.NoError(r.env.T, err)
}

func (r Receiver) Logs() (string, error) {
	return r.logs(receiverLogTail)
}

func (r Receiver) logs(tail int64) (string, error) {
	pods, err := r.env.cluster.pods(r.env.Ctx, r.env.ControlNamespace, map[string]string{types.NameLabel: ReceiverName(r.env.Release)})
	if err != nil {
		return "", err
	}
	var out strings.Builder
	for _, pod := range pods {
		logs, err := r.env.cluster.podLogs(r.env.Ctx, pod.Namespace, pod.Name, "", tail)
		if err != nil {
			return "", err
		}
		out.WriteString(logs)
	}
	return out.String(), nil
}
