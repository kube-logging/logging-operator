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

func TestNewVolumeMountBuilder(t *testing.T) {
	tests := []struct {
		name string
		want *VolumeMountBuilder
	}{
		{
			name: "contructor",
			want: &VolumeMountBuilder{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NewVolumeMountBuilder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nNewVolumeMountBuilder() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeMountBuilder_WithName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		v    *VolumeMountBuilder
		args args
		want *VolumeMountBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{name: ""},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithName",
			v:    NewVolumeMountBuilder(),
			args: args{name: "NewName"},
			want: &VolumeMountBuilder{VolumeMount: corev1.VolumeMount{Name: "NewName"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeMountBuilder.WithName() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeMountBuilder_WithMountPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		v    *VolumeMountBuilder
		args args
		want *VolumeMountBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{path: ""},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithMountPath",
			v:    NewVolumeMountBuilder(),
			args: args{path: "/var/log/foobar"},
			want: &VolumeMountBuilder{VolumeMount: corev1.VolumeMount{MountPath: "/var/log/foobar"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithMountPath(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeMountBuilder.WithMountPath() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeMountBuilder_WithSubPath(t *testing.T) {
	type args struct {
		subPath string
	}
	tests := []struct {
		name string
		v    *VolumeMountBuilder
		args args
		want *VolumeMountBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{subPath: ""},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithSubPath",
			v:    NewVolumeMountBuilder(),
			args: args{subPath: "/var/log/foobar"},
			want: &VolumeMountBuilder{VolumeMount: corev1.VolumeMount{SubPath: "/var/log/foobar"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithSubPath(tt.args.subPath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeMountBuilder.WithSubPath() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeMountBuilder_WithSubPathExpr(t *testing.T) {
	type args struct {
		subPathExpr string
	}
	tests := []struct {
		name string
		v    *VolumeMountBuilder
		args args
		want *VolumeMountBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{subPathExpr: ""},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithSubPathExpr",
			v:    NewVolumeMountBuilder(),
			args: args{subPathExpr: "/var/log/foobar"},
			want: &VolumeMountBuilder{VolumeMount: corev1.VolumeMount{SubPathExpr: "/var/log/foobar"}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithSubPathExpr(tt.args.subPathExpr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeMountBuilder.WithSubPathExpr() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeMountBuilder_WithReadOnly(t *testing.T) {
	type args struct {
		readOnly bool
	}
	tests := []struct {
		name string
		v    *VolumeMountBuilder
		args args
		want *VolumeMountBuilder
	}{
		{
			name: "nilReceiver",
			v:    nil,
			args: args{readOnly: true},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithReadOnlyTrue",
			v:    NewVolumeMountBuilder(),
			args: args{readOnly: true},
			want: &VolumeMountBuilder{VolumeMount: corev1.VolumeMount{ReadOnly: true}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithReadOnly(tt.args.readOnly); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeMountBuilder.WithReadOnly() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestVolumeMountBuilder_WithMountPropagation(t *testing.T) {
	type args struct {
		mountPropagation corev1.MountPropagationMode
	}
	tests := []struct {
		name string
		v    *VolumeMountBuilder
		args args
		want *VolumeMountBuilder
	}{
		{
			name: "nilReceiverReturnsNil",
			v:    nil,
			args: args{mountPropagation: corev1.MountPropagationHostToContainer},
			want: nil,
		},
		{
			name: "validReceiverReturnsReceiverWithMountPropagationMode",
			v:    NewVolumeMountBuilder(),
			args: args{mountPropagation: corev1.MountPropagationHostToContainer},
			want: &VolumeMountBuilder{VolumeMount: corev1.VolumeMount{MountPropagation: MountPropagationModeRef(corev1.MountPropagationHostToContainer)}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.WithMountPropagation(tt.args.mountPropagation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nVolumeMountBuilder.WithMountPropagation() = %v\nwant %v", got, tt.want)
			}
		})
	}
}
