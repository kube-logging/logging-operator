// Copyright © 2025 Kube logging authors
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
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

func scheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(s))
	require.NoError(t, appsv1.AddToScheme(s))
	require.NoError(t, v1beta1.AddToScheme(s))
	return s
}

func pod(name string, phase corev1.PodPhase, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "ns", Labels: labels},
		Status:     corev1.PodStatus{Phase: phase},
	}
}

func deployment(desired, ready, available int32, cond *appsv1.DeploymentCondition) *appsv1.Deployment {
	d := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: "d", Namespace: "ns"},
		Spec:       appsv1.DeploymentSpec{Replicas: &desired},
		Status:     appsv1.DeploymentStatus{ReadyReplicas: ready, AvailableReplicas: available},
	}
	if cond != nil {
		d.Status.Conditions = []appsv1.DeploymentCondition{*cond}
	}
	return d
}

func available(status corev1.ConditionStatus) *appsv1.DeploymentCondition {
	return &appsv1.DeploymentCondition{Type: appsv1.DeploymentAvailable, Status: status}
}

func TestPodShouldBeRunning(t *testing.T) {
	for _, c := range []struct {
		name    string
		objects []client.Object
		want    bool
	}{
		{"running", []client.Object{pod("p", corev1.PodRunning, nil)}, true},
		{"pending", []client.Object{pod("p", corev1.PodPending, nil)}, false},
		{"absent", nil, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(c.objects...).Build()
			got := PodShouldBeRunning(t, cl, client.ObjectKey{Namespace: "ns", Name: "p"})()
			require.Equal(t, c.want, got)
		})
	}
}

func TestAnyPodShouldBeRunning(t *testing.T) {
	lbl := map[string]string{"app": "x"}
	for _, c := range []struct {
		name    string
		objects []client.Object
		want    bool
	}{
		{"one running among several", []client.Object{
			pod("a", corev1.PodPending, lbl), pod("b", corev1.PodRunning, lbl),
		}, true},
		{"none running", []client.Object{
			pod("a", corev1.PodPending, lbl), pod("b", corev1.PodFailed, lbl),
		}, false},
		{"running but not matching the selector", []client.Object{
			pod("a", corev1.PodRunning, map[string]string{"app": "other"}),
		}, false},
		{"no pods at all", nil, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(c.objects...).Build()
			got := AnyPodShouldBeRunning(t, cl, client.MatchingLabels(lbl))()
			require.Equal(t, c.want, got)
		})
	}
}

func TestAnyPodShouldBeFinished(t *testing.T) {
	lbl := map[string]string{"app": "x"}
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
			got := AnyPodShouldBeFinished(t, cl, client.MatchingLabels(lbl))()
			require.Equal(t, c.want, got)
		})
	}
}

func TestResourceShouldBeAbsentAndPresent(t *testing.T) {
	existing := pod("p", corev1.PodRunning, nil)

	withPod := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(existing).Build()
	empty := fake.NewClientBuilder().WithScheme(scheme(t)).Build()

	require.False(t, ResourceShouldBeAbsent(t, withPod, pod("p", corev1.PodRunning, nil))())
	require.True(t, ResourceShouldBeAbsent(t, empty, pod("p", corev1.PodRunning, nil))())

	require.True(t, ResourceShouldBePresent(t, withPod, pod("p", corev1.PodRunning, nil))())
	require.False(t, ResourceShouldBePresent(t, empty, pod("p", corev1.PodRunning, nil))())
}

func TestDeploymentAvailable(t *testing.T) {
	for _, c := range []struct {
		name string
		dep  *appsv1.Deployment
		want bool
	}{
		{"ready, available and reporting Available", deployment(2, 2, 2, available(corev1.ConditionTrue)), true},
		{"Available condition is False", deployment(1, 1, 1, available(corev1.ConditionFalse)), false},
		{"no Available condition at all", deployment(1, 1, 1, nil), false},
		{"not all replicas ready", deployment(2, 1, 2, available(corev1.ConditionTrue)), false},
		{"not all replicas available", deployment(2, 2, 1, available(corev1.ConditionTrue)), false},
	} {
		t.Run(c.name, func(t *testing.T) {
			cl := fake.NewClientBuilder().WithScheme(scheme(t)).WithObjects(c.dep).Build()
			require.Equal(t, c.want, DeploymentAvailable(t, cl, t.Context(), "ns", "d")())
		})
	}

	t.Run("deployment absent", func(t *testing.T) {
		cl := fake.NewClientBuilder().WithScheme(scheme(t)).Build()
		require.False(t, DeploymentAvailable(t, cl, t.Context(), "ns", "d")())
	})
}
