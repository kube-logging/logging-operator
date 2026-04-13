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

func TestPodHandlerHelper(t *testing.T) {
	emptyDirVol := func(name string) corev1.Volume {
		return corev1.Volume{Name: name, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}
	}
	mount := func(name, path string) corev1.VolumeMount {
		return corev1.VolumeMount{Name: name, MountPath: path}
	}
	sidecar := func(name string, mounts ...corev1.VolumeMount) corev1.Container {
		return corev1.Container{Name: name, VolumeMounts: mounts}
	}

	tests := []struct {
		name           string
		containers     []corev1.Container
		targetIdx      int
		sideCars       []corev1.Container
		volumes        []corev1.Volume
		volumeMounts   []corev1.VolumeMount
		wantDenied     bool
		wantContainers int
		wantMounts     map[string]int // container name → expected mount count
	}{
		{
			name:           "single sidecar gets mount added to target",
			containers:     []corev1.Container{{Name: "main"}},
			targetIdx:      0,
			sideCars:       []corev1.Container{sidecar("sidecar-1", mount("vol", "/var/log/app"))},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 2,
			wantMounts:     map[string]int{"main": 1, "sidecar-1": 1},
		},
		{
			name:       "two sidecars same dir — no duplicate mounts",
			containers: []corev1.Container{{Name: "main"}},
			targetIdx:  0,
			sideCars: []corev1.Container{
				sidecar("sidecar-file1", mount("vol", "/var/log/app")),
				sidecar("sidecar-file2", mount("vol", "/var/log/app")),
			},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 3,
			wantMounts:     map[string]int{"main": 1, "sidecar-file1": 1, "sidecar-file2": 1},
		},
		{
			name:       "two sidecars different dirs — target gets both mounts",
			containers: []corev1.Container{{Name: "main"}},
			targetIdx:  0,
			sideCars: []corev1.Container{
				sidecar("sidecar-app", mount("vol-app", "/var/log/app")),
				sidecar("sidecar-sys", mount("vol-sys", "/var/log/sys")),
			},
			volumes:        []corev1.Volume{emptyDirVol("vol-app"), emptyDirVol("vol-sys")},
			volumeMounts:   []corev1.VolumeMount{mount("vol-app", "/var/log/app"), mount("vol-sys", "/var/log/sys")},
			wantContainers: 3,
			wantMounts:     map[string]int{"main": 2, "sidecar-app": 1, "sidecar-sys": 1},
		},
		{
			name:       "multi-container pod — only target gets mount",
			containers: []corev1.Container{{Name: "app"}, {Name: "nginx"}},
			targetIdx:  0, // target is "app"
			sideCars: []corev1.Container{
				sidecar("sidecar-file1", mount("vol", "/var/log/app")),
			},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 3,
			wantMounts:     map[string]int{"app": 1, "nginx": 0, "sidecar-file1": 1},
		},
		{
			name:       "duplicate sidecar name is denied",
			containers: []corev1.Container{{Name: "main"}, {Name: "existing-sidecar"}},
			targetIdx:  0,
			sideCars:   []corev1.Container{{Name: "existing-sidecar"}},
			wantDenied: true,
		},
		{
			name: "existing compatible mount on target is not duplicated",
			containers: []corev1.Container{
				{Name: "main", VolumeMounts: []corev1.VolumeMount{mount("vol", "/var/log/app")}},
			},
			targetIdx:      0,
			sideCars:       []corev1.Container{sidecar("sidecar-1", mount("vol", "/var/log/app"))},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 2,
			wantMounts:     map[string]int{"main": 1, "sidecar-1": 1},
		},
		{
			name: "conflicting mount (same path, different volume) is denied",
			containers: []corev1.Container{
				{Name: "main", VolumeMounts: []corev1.VolumeMount{mount("my-logs", "/var/log/app")}},
			},
			targetIdx:    0,
			sideCars:     []corev1.Container{sidecar("sidecar-1", mount("webhook-vol", "/var/log/app"))},
			volumes:      []corev1.Volume{emptyDirVol("webhook-vol")},
			volumeMounts: []corev1.VolumeMount{mount("webhook-vol", "/var/log/app")},
			wantDenied:   true,
		},
		{
			name:       "out of range targetContainerIdx is denied",
			containers: []corev1.Container{{Name: "main"}},
			targetIdx:  5,
			sideCars:   []corev1.Container{sidecar("sidecar-1", mount("vol", "/var/log/app"))},
			wantDenied: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestPodHandler()
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				Spec:       corev1.PodSpec{Containers: tt.containers},
			}

			resp := p.podHandlerHelper(pod, tt.targetIdx, tt.sideCars, tt.volumes, tt.volumeMounts)

			if tt.wantDenied {
				if resp == nil {
					t.Fatal("expected Denied response, got nil")
				}
				return
			}
			if resp != nil {
				t.Fatalf("unexpected Denied response: %v", resp)
			}
			if got := len(pod.Spec.Containers); got != tt.wantContainers {
				t.Fatalf("containers: want %d, got %d", tt.wantContainers, got)
			}
			for _, c := range pod.Spec.Containers {
				if want, ok := tt.wantMounts[c.Name]; ok {
					if got := len(c.VolumeMounts); got != want {
						t.Errorf("container %q: want %d volumeMounts, got %d", c.Name, want, got)
					}
				}
			}
		})
	}
}

// TestPodHandlerHelper_MultiCallSimulatesHandle simulates the Handle() pattern:
// multiple sequential calls to podHandlerHelper for a multi-container pod,
// where each call targets a different original container's files.
func TestPodHandlerHelper_MultiCallSimulatesHandle(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "multi-call-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app"},
				{Name: "nginx"},
			},
		},
	}

	mount := func(name, path string) corev1.VolumeMount {
		return corev1.VolumeMount{Name: name, MountPath: path}
	}
	emptyDirVol := func(name string) corev1.Volume {
		return corev1.Volume{Name: name, VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}
	}

	// Call 1: sidecars for "app" container (2 files in /var/log/app), targetIdx=0
	sideCars1 := []corev1.Container{
		{Name: "app-file1-sidecar", VolumeMounts: []corev1.VolumeMount{mount("vol-app", "/var/log/app")}},
		{Name: "app-file2-sidecar", VolumeMounts: []corev1.VolumeMount{mount("vol-app", "/var/log/app")}},
	}
	volumes1 := []corev1.Volume{emptyDirVol("vol-app")}
	mounts1 := []corev1.VolumeMount{mount("vol-app", "/var/log/app")}

	resp := p.podHandlerHelper(pod, 0, sideCars1, volumes1, mounts1)
	if resp != nil {
		t.Fatalf("call 1: unexpected Denied: %v", resp)
	}

	// After call 1: [app, nginx, app-file1-sidecar, app-file2-sidecar]
	if got := len(pod.Spec.Containers); got != 4 {
		t.Fatalf("after call 1: want 4 containers, got %d", got)
	}

	// Call 2: sidecars for "nginx" container (1 file in /var/log/nginx), targetIdx=1
	sideCars2 := []corev1.Container{
		{Name: "nginx-access-sidecar", VolumeMounts: []corev1.VolumeMount{mount("vol-nginx", "/var/log/nginx")}},
	}
	volumes2 := []corev1.Volume{emptyDirVol("vol-nginx")}
	mounts2 := []corev1.VolumeMount{mount("vol-nginx", "/var/log/nginx")}

	resp = p.podHandlerHelper(pod, 1, sideCars2, volumes2, mounts2)
	if resp != nil {
		t.Fatalf("call 2: unexpected Denied: %v", resp)
	}

	// After call 2: [app, nginx, app-file1-sidecar, app-file2-sidecar, nginx-access-sidecar]
	if got := len(pod.Spec.Containers); got != 5 {
		t.Fatalf("after call 2: want 5 containers, got %d", got)
	}

	// Verify scoped mounts: each original container only gets its own mount,
	// not the other container's mount.
	wantMounts := map[string][]string{
		"app":                  {"/var/log/app"},
		"nginx":                {"/var/log/nginx"},
		"app-file1-sidecar":    {"/var/log/app"},
		"app-file2-sidecar":    {"/var/log/app"},
		"nginx-access-sidecar": {"/var/log/nginx"},
	}

	for _, c := range pod.Spec.Containers {
		expected, ok := wantMounts[c.Name]
		if !ok {
			t.Errorf("unexpected container %q", c.Name)
			continue
		}
		if got := len(c.VolumeMounts); got != len(expected) {
			t.Errorf("container %q: want %d mounts %v, got %d", c.Name, len(expected), expected, got)
		}
		for i, vm := range c.VolumeMounts {
			if i < len(expected) && vm.MountPath != expected[i] {
				t.Errorf("container %q mount[%d]: want path %q, got %q", c.Name, i, expected[i], vm.MountPath)
			}
		}
	}
}

// TestPodHandlerHelper_DuplicateVolumeSkipped verifies that when the same
// volume name is passed in two sequential calls (e.g., two containers sharing
// a log directory), the volume is added once and the second call skips it
// gracefully instead of returning Denied.
func TestPodHandlerHelper_DuplicateVolumeSkipped(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "app"},
				{Name: "worker"},
			},
		},
	}

	mount := func(name, path string) corev1.VolumeMount {
		return corev1.VolumeMount{Name: name, MountPath: path}
	}
	sharedVol := corev1.Volume{Name: "shared-vol", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}

	// Call 1: targets "app"
	resp := p.podHandlerHelper(pod, 0,
		[]corev1.Container{{Name: "sc-app", VolumeMounts: []corev1.VolumeMount{mount("shared-vol", "/var/log/shared")}}},
		[]corev1.Volume{sharedVol},
		[]corev1.VolumeMount{mount("shared-vol", "/var/log/shared")},
	)
	if resp != nil {
		t.Fatalf("call 1: unexpected Denied: %v", resp)
	}

	// Call 2: targets "worker" with the same volume name
	resp = p.podHandlerHelper(pod, 1,
		[]corev1.Container{{Name: "sc-worker", VolumeMounts: []corev1.VolumeMount{mount("shared-vol", "/var/log/shared")}}},
		[]corev1.Volume{sharedVol},
		[]corev1.VolumeMount{mount("shared-vol", "/var/log/shared")},
	)
	if resp != nil {
		t.Fatalf("call 2: unexpected Denied (volume should be skipped, not denied): %v", resp)
	}

	// Volume should appear exactly once
	volCount := 0
	for _, v := range pod.Spec.Volumes {
		if v.Name == "shared-vol" {
			volCount++
		}
	}
	if volCount != 1 {
		t.Errorf("want 1 volume named shared-vol, got %d", volCount)
	}
}

// TestPodHandlerHelper_IncompatibleVolumeSourceDenied verifies that if the pod
// already has a volume with the same name but a non-EmptyDir source, the
// mutation is denied rather than silently proceeding.
func TestPodHandlerHelper_IncompatibleVolumeSourceDenied(t *testing.T) {
	p := newTestPodHandler()

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{Name: "main"}},
			Volumes: []corev1.Volume{
				{Name: "vol", VolumeSource: corev1.VolumeSource{
					HostPath: &corev1.HostPathVolumeSource{Path: "/host/logs"},
				}},
			},
		},
	}

	mount := func(name, path string) corev1.VolumeMount {
		return corev1.VolumeMount{Name: name, MountPath: path}
	}

	resp := p.podHandlerHelper(pod, 0,
		[]corev1.Container{{Name: "sidecar-1", VolumeMounts: []corev1.VolumeMount{mount("vol", "/var/log/app")}}},
		[]corev1.Volume{{Name: "vol", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}}},
		[]corev1.VolumeMount{mount("vol", "/var/log/app")},
	)
	if resp == nil {
		t.Fatal("expected Denied for incompatible volume source, got nil")
	}
}

func TestFindVolumeMount(t *testing.T) {
	tests := []struct {
		name      string
		mounts    []corev1.VolumeMount
		mountPath string
		wantFound bool
		wantName  string
	}{
		{"found", []corev1.VolumeMount{{Name: "a", MountPath: "/var/log/app"}}, "/var/log/app", true, "a"},
		{"not found", []corev1.VolumeMount{{Name: "a", MountPath: "/var/log/app"}}, "/var/log/other", false, ""},
		{"nil slice", nil, "/var/log/app", false, ""},
		{"empty slice", []corev1.VolumeMount{}, "/var/log/app", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := findVolumeMount(tt.mounts, tt.mountPath)
			if found != tt.wantFound {
				t.Errorf("findVolumeMount(%v, %q) found=%v, want %v", tt.mounts, tt.mountPath, found, tt.wantFound)
			}
			if found && got.Name != tt.wantName {
				t.Errorf("findVolumeMount(%v, %q) name=%q, want %q", tt.mounts, tt.mountPath, got.Name, tt.wantName)
			}
		})
	}
}
