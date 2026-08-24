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
	"os/exec"
	"strings"

	"github.com/cisco-open/operator-tools/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/e2e/common"
)

const (
	receiverPort = 8080

	// Generous on purpose: more lines makes a tag easier to find and its
	// absence harder to claim.
	receiverLogTail = 100
)

// Receiver is the test receiver the chart installs, where a suite looks to see
// which logs arrived.
type Receiver struct {
	env *Env
}

func ReceiverName(release string) string { return release + "-test-receiver" }

// These take a release rather than an Env so the unmigrated suites can use
// them too. They go when the last one migrates, leaving the methods.
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
// for. It belongs after whatever wait establishes that the pipeline is running.
func (r Receiver) MustNotReceive(tags ...string) {
	r.env.T.Helper()

	logs, err := r.Logs()
	require.NoError(r.env.T, err)
	for _, tag := range tags {
		assert.NotContains(r.env.T, logs, tag)
	}
}

func (r Receiver) Logs() (string, error) {
	out, err := common.CmdEnv(exec.Command("kubectl",
		"logs",
		"-n", r.env.ControlNamespace,
		"--tail", fmt.Sprint(receiverLogTail),
		"-l", fmt.Sprintf("%s=%s", types.NameLabel, ReceiverName(r.env.Release))), r.env.Cluster).Output()
	return string(out), err
}
