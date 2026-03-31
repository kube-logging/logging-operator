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
		sideCars       []corev1.Container
		volumes        []corev1.Volume
		volumeMounts   []corev1.VolumeMount
		wantDenied     bool
		wantContainers int
		wantMounts     map[string]int // container name → expected mount count
	}{
		{
			name:           "single sidecar gets mount added to main",
			containers:     []corev1.Container{{Name: "main"}},
			sideCars:       []corev1.Container{sidecar("sidecar-1", mount("vol", "/var/log/app"))},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 2,
			wantMounts:     map[string]int{"main": 1, "sidecar-1": 1},
		},
		{
			name:       "two sidecars same dir — no duplicate mounts",
			containers: []corev1.Container{{Name: "main"}},
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
			name:       "two sidecars different dirs — main gets both mounts",
			containers: []corev1.Container{{Name: "main"}},
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
			name:       "multi-container pod — both originals get mounts",
			containers: []corev1.Container{{Name: "app"}, {Name: "nginx"}},
			sideCars: []corev1.Container{
				sidecar("sidecar-file1", mount("vol", "/var/log/app")),
				sidecar("sidecar-file2", mount("vol", "/var/log/app")),
			},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 4,
			wantMounts:     map[string]int{"app": 1, "nginx": 1, "sidecar-file1": 1, "sidecar-file2": 1},
		},
		{
			name:       "duplicate sidecar name is denied",
			containers: []corev1.Container{{Name: "main"}, {Name: "existing-sidecar"}},
			sideCars:   []corev1.Container{{Name: "existing-sidecar"}},
			wantDenied: true,
		},
		{
			name: "existing mount on main is not duplicated",
			containers: []corev1.Container{
				{Name: "main", VolumeMounts: []corev1.VolumeMount{mount("existing", "/var/log/app")}},
			},
			sideCars:       []corev1.Container{sidecar("sidecar-1", mount("vol", "/var/log/app"))},
			volumes:        []corev1.Volume{emptyDirVol("vol")},
			volumeMounts:   []corev1.VolumeMount{mount("vol", "/var/log/app")},
			wantContainers: 2,
			wantMounts:     map[string]int{"main": 1, "sidecar-1": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newTestPodHandler()
			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: "test-pod"},
				Spec:       corev1.PodSpec{Containers: tt.containers},
			}

			resp := p.podHandlerHelper(pod, tt.sideCars, tt.volumes, tt.volumeMounts)

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

func TestHasVolumeMount(t *testing.T) {
	tests := []struct {
		name      string
		mounts    []corev1.VolumeMount
		mountPath string
		want      bool
	}{
		{"found", []corev1.VolumeMount{{Name: "a", MountPath: "/var/log/app"}}, "/var/log/app", true},
		{"not found", []corev1.VolumeMount{{Name: "a", MountPath: "/var/log/app"}}, "/var/log/other", false},
		{"nil slice", nil, "/var/log/app", false},
		{"empty slice", []corev1.VolumeMount{}, "/var/log/app", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasVolumeMount(tt.mounts, tt.mountPath); got != tt.want {
				t.Errorf("hasVolumeMount(%v, %q) = %v, want %v", tt.mounts, tt.mountPath, got, tt.want)
			}
		})
	}
}
