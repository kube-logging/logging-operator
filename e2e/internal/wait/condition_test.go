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
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// The ptr.Deref default is chosen per direction: false where Active must be
// true, true where it must be false. Both leave an unset Active unsatisfied so
// the caller polls again. One backwards would either pass on an unreconciled
// config or never pass at all, so cover all three states.
func TestAttachedAndExcessAcrossEveryActiveState(t *testing.T) {
	yes, no := true, false

	for _, c := range []struct {
		name            string
		active          *bool
		attach, exclude bool
	}{
		{"active true", &yes, true, false},
		{"active false", &no, false, true},
		{"active unset", nil, false, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.attach, attached(nil, "lg", c.active, "lg"))
			assert.Equal(t, c.exclude, excess([]string{"boom"}, "", c.active))
		})
	}
}

// Each predicate needs all three of its fields, or a config halfway through
// reconciling reads as settled.
func TestAttachedAndExcessNeedEveryField(t *testing.T) {
	yes, no := true, false

	assert.False(t, attached([]string{"boom"}, "lg", &yes, "lg"), "problems reported")
	assert.False(t, attached(nil, "other", &yes, "lg"), "naming a different Logging")
	assert.True(t, attached(nil, "lg", &yes, "lg"))

	assert.False(t, excess(nil, "", &no), "no problems reported")
	assert.False(t, excess([]string{"boom"}, "lg", &no), "still naming a Logging")
	assert.True(t, excess([]string{"boom"}, "", &no))
}

// A config check reports by exiting, so both terminal phases count and neither
// live one does.
func TestConfigCheckFinished(t *testing.T) {
	lbl := map[string]string{"my-unique-label": "configcheck"}

	for _, c := range []struct {
		name  string
		phase corev1.PodPhase
		want  bool
	}{
		{"succeeded counts as finished", corev1.PodSucceeded, true},
		{"failed counts as finished", corev1.PodFailed, true},
		{"running does not", corev1.PodRunning, false},
		{"pending does not", corev1.PodPending, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).
				WithObjects(pod("a", c.phase, lbl)).Build()

			met, err := ConfigCheck(lbl).Met(t.Context(), cl)
			require.NoError(t, err)
			assert.Equal(t, c.want, met)
		})
	}
}

// Cleared is not the negation of healthy: an unrelated problem leaves the
// config check's own one cleared while the Logging is still unhealthy.
func TestLoggingProblemConditions(t *testing.T) {
	failure := regexp.MustCompile(`^Configuration with checksum (.+) has failed. .*`)

	for _, c := range []struct {
		name                       string
		problems                   []string
		healthy, reported, cleared bool
	}{
		{"no problems", nil, true, false, true},
		{"the config check failed", []string{"Configuration with checksum abc has failed. boom"}, false, true, false},
		{"some other problem", []string{"unrelated"}, false, false, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(&v1beta1.Logging{
				ObjectMeta: metav1.ObjectMeta{Name: "lg"},
				Status:     v1beta1.LoggingStatus{Problems: c.problems, ProblemsCount: len(c.problems)},
			}).Build()

			healthy, err := LoggingHealthy("lg").Met(t.Context(), cl)
			require.NoError(t, err)
			assert.Equal(t, c.healthy, healthy)

			reported, err := LoggingProblem("lg", failure.MatchString).Met(t.Context(), cl)
			require.NoError(t, err)
			assert.Equal(t, c.reported, reported)

			cleared, err := LoggingProblemCleared("lg", failure.MatchString).Met(t.Context(), cl)
			require.NoError(t, err)
			assert.Equal(t, c.cleared, cleared)
		})
	}
}

// Active alone would be a race: a drainer that finishes between two polls is
// never seen active, and the wait then burns its budget on a Job that did run.
func TestJobStartedCountsTerminalStatesToo(t *testing.T) {
	for _, c := range []struct {
		name   string
		status batchv1.JobStatus
		want   bool
	}{
		{"created but not started", batchv1.JobStatus{}, false},
		{"running", batchv1.JobStatus{Active: 1}, true},
		{"already succeeded", batchv1.JobStatus{Succeeded: 1}, true},
		{"already failed", batchv1.JobStatus{Failed: 1}, true},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(&batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{Name: "d", Namespace: "ns"},
				Status:     c.status,
			}).Build()

			met, err := JobStarted("ns", "d").Met(t.Context(), cl)
			require.NoError(t, err)
			assert.Equal(t, c.want, met)
		})
	}
}

// A Deployment wanting no replicas would otherwise read as ready: nothing is
// unready, and the Available condition can still be True from before it was
// scaled down.
func TestDeploymentWantingNoReplicasIsNeverReady(t *testing.T) {
	zero := int32(0)
	for _, c := range []struct {
		name     string
		replicas *int32
	}{
		{"replicas unset", nil},
		{"explicitly zero", &zero},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "es", Namespace: "ns"},
				Spec:       appsv1.DeploymentSpec{Replicas: c.replicas},
				Status: appsv1.DeploymentStatus{Conditions: []appsv1.DeploymentCondition{
					{Type: appsv1.DeploymentAvailable, Status: corev1.ConditionTrue},
				}},
			}).Build()

			met, err := Deployment("ns", "es").Met(t.Context(), cl)
			require.NoError(t, err)
			assert.False(t, met)
		})
	}
}

// Ready alone can hold while the Available condition still reports the old
// rollout, so all three have to agree.
func TestDeploymentNeedsEveryReplicaAndTheCondition(t *testing.T) {
	two := int32(2)
	availableCond := func(s corev1.ConditionStatus) []appsv1.DeploymentCondition {
		return []appsv1.DeploymentCondition{{Type: appsv1.DeploymentAvailable, Status: s}}
	}

	for _, c := range []struct {
		name   string
		status appsv1.DeploymentStatus
		want   bool
	}{
		{"all ready and available", appsv1.DeploymentStatus{ReadyReplicas: 2, AvailableReplicas: 2, Conditions: availableCond(corev1.ConditionTrue)}, true},
		{"not all ready", appsv1.DeploymentStatus{ReadyReplicas: 1, AvailableReplicas: 2, Conditions: availableCond(corev1.ConditionTrue)}, false},
		{"not all available", appsv1.DeploymentStatus{ReadyReplicas: 2, AvailableReplicas: 1, Conditions: availableCond(corev1.ConditionTrue)}, false},
		{"condition still false", appsv1.DeploymentStatus{ReadyReplicas: 2, AvailableReplicas: 2, Conditions: availableCond(corev1.ConditionFalse)}, false},
		{"no condition yet", appsv1.DeploymentStatus{ReadyReplicas: 2, AvailableReplicas: 2}, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(&appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{Name: "es", Namespace: "ns"},
				Spec:       appsv1.DeploymentSpec{Replicas: &two},
				Status:     c.status,
			}).Build()

			met, err := Deployment("ns", "es").Met(t.Context(), cl)
			require.NoError(t, err)
			assert.Equal(t, c.want, met)
		})
	}
}
