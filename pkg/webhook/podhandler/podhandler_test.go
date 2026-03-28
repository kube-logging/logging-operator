// Copyright © 2024 Kube logging authors
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

package podhandler

import (
	"testing"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newTestPodHandler() *PodHandler {
	return &PodHandler{Log: logr.Discard()}
}

func TestPodHandlerHelper_SingleSidecar(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "main"},
			},
		},
	}

	sideCars := []corev1.Container{
		{
			Name: "sidecar-1",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol", MountPath: "/var/log/app"},
			},
		},
	}
	volumes := []corev1.Volume{
		{Name: "vol", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
	}
	volumeMounts := []corev1.VolumeMount{
		{Name: "vol", MountPath: "/var/log/app"},
	}

	resp := p.podHandlerHelper(pod, sideCars, volumes, volumeMounts)
	if resp != nil {
		t.Fatalf("unexpected response: %v", resp)
	}

	if len(pod.Spec.Containers) != 2 {
		t.Fatalf("expected 2 containers, got %d", len(pod.Spec.Containers))
	}

	// main container should have the volumeMount
	mainMounts := pod.Spec.Containers[0].VolumeMounts
	if len(mainMounts) != 1 {
		t.Errorf("main container: expected 1 volumeMount, got %d", len(mainMounts))
	}

	// sidecar should have exactly 1 volumeMount (from creation, not duplicated)
	sidecarMounts := pod.Spec.Containers[1].VolumeMounts
	if len(sidecarMounts) != 1 {
		t.Errorf("sidecar: expected 1 volumeMount, got %d", len(sidecarMounts))
	}
}

func TestPodHandlerHelper_TwoSidecarsSameDir(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "main"},
			},
		},
	}

	// Two sidecars tailing two files in the same directory.
	// Each sidecar already has its own volumeMount from Container().
	sideCars := []corev1.Container{
		{
			Name: "sidecar-file1",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol", MountPath: "/var/log/app"},
			},
		},
		{
			Name: "sidecar-file2",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol", MountPath: "/var/log/app"},
			},
		},
	}
	volumes := []corev1.Volume{
		{Name: "vol", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
	}
	volumeMounts := []corev1.VolumeMount{
		{Name: "vol", MountPath: "/var/log/app"},
	}

	resp := p.podHandlerHelper(pod, sideCars, volumes, volumeMounts)
	if resp != nil {
		t.Fatalf("unexpected response: %v", resp)
	}

	if len(pod.Spec.Containers) != 3 {
		t.Fatalf("expected 3 containers (main + 2 sidecars), got %d", len(pod.Spec.Containers))
	}

	// main container should have the shared volumeMount
	mainMounts := pod.Spec.Containers[0].VolumeMounts
	if len(mainMounts) != 1 {
		t.Errorf("main container: expected 1 volumeMount, got %d", len(mainMounts))
	}

	// Each sidecar should have exactly 1 volumeMount (no duplicates)
	for i := 1; i <= 2; i++ {
		mounts := pod.Spec.Containers[i].VolumeMounts
		if len(mounts) != 1 {
			t.Errorf("sidecar %d (%s): expected 1 volumeMount, got %d",
				i, pod.Spec.Containers[i].Name, len(mounts))
		}
	}
}

func TestPodHandlerHelper_TwoSidecarsDifferentDirs(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "main"},
			},
		},
	}

	sideCars := []corev1.Container{
		{
			Name: "sidecar-app",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol-app", MountPath: "/var/log/app"},
			},
		},
		{
			Name: "sidecar-sys",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol-sys", MountPath: "/var/log/sys"},
			},
		},
	}
	volumes := []corev1.Volume{
		{Name: "vol-app", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
		{Name: "vol-sys", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
	}
	volumeMounts := []corev1.VolumeMount{
		{Name: "vol-app", MountPath: "/var/log/app"},
		{Name: "vol-sys", MountPath: "/var/log/sys"},
	}

	resp := p.podHandlerHelper(pod, sideCars, volumes, volumeMounts)
	if resp != nil {
		t.Fatalf("unexpected response: %v", resp)
	}

	if len(pod.Spec.Containers) != 3 {
		t.Fatalf("expected 3 containers, got %d", len(pod.Spec.Containers))
	}

	// main container should have both volumeMounts
	mainMounts := pod.Spec.Containers[0].VolumeMounts
	if len(mainMounts) != 2 {
		t.Errorf("main container: expected 2 volumeMounts, got %d", len(mainMounts))
	}

	// Each sidecar should have exactly 1 volumeMount
	for i := 1; i <= 2; i++ {
		mounts := pod.Spec.Containers[i].VolumeMounts
		if len(mounts) != 1 {
			t.Errorf("sidecar %d (%s): expected 1 volumeMount, got %d",
				i, pod.Spec.Containers[i].Name, len(mounts))
		}
	}
}

func TestPodHandlerHelper_MultiContainerPod(t *testing.T) {
	p := newTestPodHandler()

	// Pod with 2 original containers
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app"},
				{Name: "nginx"},
			},
		},
	}

	sideCars := []corev1.Container{
		{
			Name: "sidecar-file1",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol", MountPath: "/var/log/app"},
			},
		},
		{
			Name: "sidecar-file2",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol", MountPath: "/var/log/app"},
			},
		},
	}
	volumes := []corev1.Volume{
		{Name: "vol", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
	}
	volumeMounts := []corev1.VolumeMount{
		{Name: "vol", MountPath: "/var/log/app"},
	}

	resp := p.podHandlerHelper(pod, sideCars, volumes, volumeMounts)
	if resp != nil {
		t.Fatalf("unexpected response: %v", resp)
	}

	if len(pod.Spec.Containers) != 4 {
		t.Fatalf("expected 4 containers (2 original + 2 sidecars), got %d", len(pod.Spec.Containers))
	}

	// Both original containers should have the mount
	for i := 0; i < 2; i++ {
		mounts := pod.Spec.Containers[i].VolumeMounts
		if len(mounts) != 1 {
			t.Errorf("original container %d (%s): expected 1 volumeMount, got %d",
				i, pod.Spec.Containers[i].Name, len(mounts))
		}
	}

	// Sidecars should have exactly 1 mount each
	for i := 2; i < 4; i++ {
		mounts := pod.Spec.Containers[i].VolumeMounts
		if len(mounts) != 1 {
			t.Errorf("sidecar %d (%s): expected 1 volumeMount, got %d",
				i, pod.Spec.Containers[i].Name, len(mounts))
		}
	}
}

func TestPodHandlerHelper_DuplicateSidecarDenied(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "main"},
				{Name: "existing-sidecar"}, // already present
			},
		},
	}

	sideCars := []corev1.Container{
		{Name: "existing-sidecar"}, // name collision
	}

	resp := p.podHandlerHelper(pod, sideCars, nil, nil)
	if resp == nil {
		t.Fatal("expected Denied response for duplicate sidecar name, got nil")
	}
}

func TestPodHandlerHelper_ExistingMountNotDuplicated(t *testing.T) {
	p := newTestPodHandler()

	// main container already has a mount at /var/log/app
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name: "main",
					VolumeMounts: []corev1.VolumeMount{
						{Name: "existing", MountPath: "/var/log/app"},
					},
				},
			},
		},
	}

	sideCars := []corev1.Container{
		{
			Name: "sidecar-1",
			VolumeMounts: []corev1.VolumeMount{
				{Name: "vol", MountPath: "/var/log/app"},
			},
		},
	}
	volumes := []corev1.Volume{
		{Name: "vol", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
	}
	volumeMounts := []corev1.VolumeMount{
		{Name: "vol", MountPath: "/var/log/app"},
	}

	resp := p.podHandlerHelper(pod, sideCars, volumes, volumeMounts)
	if resp != nil {
		t.Fatalf("unexpected response: %v", resp)
	}

	// main should still have exactly 1 mount (not duplicated)
	mainMounts := pod.Spec.Containers[0].VolumeMounts
	if len(mainMounts) != 1 {
		t.Errorf("main container: expected 1 volumeMount (existing), got %d", len(mainMounts))
	}
}

func TestHasVolumeMount(t *testing.T) {
	mounts := []corev1.VolumeMount{
		{Name: "a", MountPath: "/var/log/app"},
		{Name: "b", MountPath: "/var/log/sys"},
	}

	if !hasVolumeMount(mounts, "/var/log/app") {
		t.Error("expected true for /var/log/app")
	}
	if !hasVolumeMount(mounts, "/var/log/sys") {
		t.Error("expected true for /var/log/sys")
	}
	if hasVolumeMount(mounts, "/var/log/other") {
		t.Error("expected false for /var/log/other")
	}
	if hasVolumeMount(nil, "/var/log/app") {
		t.Error("expected false for nil slice")
	}
}
