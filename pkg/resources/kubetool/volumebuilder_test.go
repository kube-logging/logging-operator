// Copyright © 2019 Banzai Cloud
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

package kubetool

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestNewVolumeBuilder(t *testing.T) {
	tests := []struct {
		name string
		want *VolumeBuilder
	}{
		{
			name: "constructor",
			want: &VolumeBuilder{Volume: corev1.Volume{}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVolumeBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewVolumeBuilder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVolumeBuilder_WithName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		v    *VolumeBuilder
		args args
		want *VolumeBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{name: "NewName"},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithNewName",
			v:    NewVolumeBuilder(),
			args: args{name: "NewName"},
			want: &VolumeBuilder{Volume: corev1.Volume{Name: "NewName"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeBuilder.WithName() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeBuilder_WithVolumeSource(t *testing.T) {
	type args struct {
		volumeSource corev1.VolumeSource
	}
	tests := []struct {
		name string
		v    *VolumeBuilder
		args args
		want *VolumeBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{volumeSource: corev1.VolumeSource{}},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithNewVolumeSource",
			v:    NewVolumeBuilder(),
			args: args{volumeSource: corev1.VolumeSource{}},
			want: &VolumeBuilder{Volume: corev1.Volume{VolumeSource: corev1.VolumeSource{}}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithVolumeSource(tt.args.volumeSource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeBuilder.WithVolumeSource() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeBuilder_WithEmptyDir(t *testing.T) {
	type args struct {
		emptyDir corev1.EmptyDirVolumeSource
	}
	tests := []struct {
		name string
		v    *VolumeBuilder
		args args
		want *VolumeBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{emptyDir: corev1.EmptyDirVolumeSource{}},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithEmptyDir",
			v:    NewVolumeBuilder(),
			args: args{emptyDir: corev1.EmptyDirVolumeSource{Medium: corev1.StorageMediumMemory}}, // non default Medium
			want: &VolumeBuilder{Volume: corev1.Volume{VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{Medium: corev1.StorageMediumMemory}}}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithEmptyDir(tt.args.emptyDir); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeBuilder.WithEmptyDir() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeBuilder_WithHostPath(t *testing.T) {
	type args struct {
		hostPath corev1.HostPathVolumeSource
	}
	tests := []struct {
		name string
		v    *VolumeBuilder
		args args
		want *VolumeBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{hostPath: corev1.HostPathVolumeSource{}},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithHostPath",
			v:    NewVolumeBuilder(),
			args: args{hostPath: corev1.HostPathVolumeSource{Path: "/foo/bar/baz"}}, // non default Medium
			want: &VolumeBuilder{Volume: corev1.Volume{VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/foo/bar/baz"}}}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithHostPath(tt.args.hostPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeBuilder.WithHostPath() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeBuilder_WithHostPathFromPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		v    *VolumeBuilder
		args args
		want *VolumeBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{path: "/foo/bar/baz"},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithNewHostPath",
			v:    NewVolumeBuilder(),
			args: args{path: "/foo/bar/baz"},
			want: &VolumeBuilder{Volume: corev1.Volume{VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/foo/bar/baz"}}}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithHostPathFromPath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeBuilder.WithHostPathFromPath() = %v\nwant %v", got, tt.want)
			}
		})
	}
}
