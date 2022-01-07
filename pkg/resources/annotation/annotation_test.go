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

package annotation

import (
	"reflect"
	"testing"
)

func TestNewHandler(t *testing.T) {
	type args struct {
		containerNames []string
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "empty",
			args: args{containerNames: []string{}},
			want: &Handler{
				containerPaths:       make(ContainerPaths),
				defaultContainerName: "",
				Config:               defaults,
			},
		},
		{
			name: "withContainerNames",
			args: args{containerNames: []string{"foo", "bar", "baz"}},
			want: &Handler{
				containerPaths:       ContainerPaths{"foo": FilePaths{}, "bar": FilePaths{}, "baz": FilePaths{}},
				defaultContainerName: "foo",
				Config:               defaults,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.containerNames); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTailerDescriptor(t *testing.T) {
	type args struct {
		tailerDescriptor TailerDescriptor
	}
	tests := []struct {
		name string
		h    *Handler
		args args
		want *Handler
	}{
		{
			name: "allnil",
			h:    nil,
			args: args{tailerDescriptor: ""},
			want: nil,
		},
		{
			name: "emptyHandler",
			h:    NewHandler([]string{}),
			args: args{tailerDescriptor: "foo:/var/foo"},
			want: NewHandler([]string{}),
		},
		{
			name: "emptyDescriptor",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: ""},
			want: NewHandler([]string{"foo"}),
		},
		{
			name: "descriptorWithValidKey",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: "foo:/var/foo"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo"}
				return result
			}(),
		},
		{
			name: "descriptorWithInvalidKeySkipsDescriptor",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: "bar:/var/foo"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{}
				return result
			}(),
		},
		{
			name: "descriptorWithoutKeyUsesDefaultContainerName",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: "/var/foo"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo"}
				return result
			}(),
		},
		{
			name: "descriptorWithoutPath",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: "bar:"},
			want: NewHandler([]string{"foo"}),
		},
		{
			name: "descriptorWithInvalidPath",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: "bar:invalid/path"},
			want: NewHandler([]string{"foo"}),
		},
		{
			name: "descriptorWithSpecialChars",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerDescriptor: "foo:/var/pos/tail-foo_bar.db"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/pos/tail-foo_bar.db"}
				return result
			}(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.h.addTailerDescriptor(tt.args.tailerDescriptor)
			if !reflect.DeepEqual(tt.h, tt.want) {
				t.Errorf("\naddTailerDescriptor() = %v\nwant %v", tt.h, tt.want)
			}
		})
	}
}

func TestAddTailerAnnotation(t *testing.T) {
	type args struct {
		tailerAnnotation TailerAnnotation
	}
	tests := []struct {
		name string
		h    *Handler
		args args
		want *Handler
	}{
		{
			name: "allnil",
			h:    nil,
			args: args{tailerAnnotation: ""},
			want: nil,
		},
		{
			name: "emptyHandler",
			h:    NewHandler([]string{}),
			args: args{tailerAnnotation: "foo:/var/foo"},
			want: NewHandler([]string{}),
		},
		{
			name: "emptyAnnotation",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: ""},
			want: NewHandler([]string{"foo"}),
		},
		{
			name: "validSingleAnnotation",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo"}
				return result
			}(),
		},
		{
			name: "validSingleAnnotationWithUnnecessarySeparatorAtEnd",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo,"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo"}
				return result
			}(),
		},
		{
			name: "validSingleAnnotationWithUnnecessarySeparatorAtBeginning",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: ",foo:/var/foo"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo"}
				return result
			}(),
		},
		{
			name: "validMultiAnnotation",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo,foo:/var/log/zzz"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo", "/var/log/zzz"}
				return result
			}(),
		},
		{
			name: "partiallyValidMultiAnnotation",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo,foo:var/log/zzz"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo"}
				return result
			}(),
		},
		{
			name: "validMultiAnnotationSpaceAfterDelimiterSkipped",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo,     foo:/var/log/zzz"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo", "/var/log/zzz"}
				return result
			}(),
		},
		{
			name: "validMultiAnnotationSpaceBeforeDelimiterSkipped",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo     ,foo:/var/log/zzz"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo", "/var/log/zzz"}
				return result
			}(),
		},
		{
			name: "validMultiAnnotationSpaceAroundDelimiterSkipped",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/foo   ,  foo:/var/log/zzz"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/foo", "/var/log/zzz"}
				return result
			}(),
		},
		{
			name: "emptyDescriptorInAnnotation",
			h:    NewHandler([]string{"foo"}),
			args: args{tailerAnnotation: "foo:/var/log,,,foo:/var/foo"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/log", "/var/foo"}
				return result
			}(),
		},
		{
			name: "multiKeyMultiPathAnnotation",
			h:    NewHandler([]string{"foo", "bar"}),
			args: args{tailerAnnotation: "foo:/var/log,foo:/var/foo,bar:/dev/zero,bar:/dev/urandom"},
			want: func() *Handler {
				result := NewHandler([]string{"foo"})
				result.containerPaths["foo"] = []string{"/var/log", "/var/foo"}
				result.containerPaths["bar"] = []string{"/dev/zero", "/dev/urandom"}
				return result
			}(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.h.AddTailerAnnotation(tt.args.tailerAnnotation)
			if !reflect.DeepEqual(tt.h, tt.want) {
				t.Errorf("\nAddTailerAnnotation() = %v\nwant %v", tt.h, tt.want)
			}
		})
	}
}

func TestFilePathsForContainer(t *testing.T) {
	type args struct {
		containerName string
	}
	tests := []struct {
		name             string
		h                *Handler
		tailerAnnotation TailerAnnotation
		args             args
		want             FilePaths
	}{
		{
			name:             "nilHandler",
			h:                nil,
			tailerAnnotation: "",
			args:             args{containerName: ""},
			want:             nil,
		},
		{
			name:             "emptyHandler",
			h:                NewHandler([]string{}),
			tailerAnnotation: "",
			args:             args{containerName: ""},
			want:             FilePaths{},
		},
		{
			name:             "handlerWithEmptyPaths",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "",
			args:             args{containerName: "foo"},
			want:             FilePaths{},
		},
		{
			name:             "handlerWithInvalidArgs",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "foo:/var/log/nginx/access.log",
			args:             args{containerName: "Invalid"},
			want:             FilePaths{},
		},
		{
			name:             "singlePath",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "foo:/var/log/nginx/access.log",
			args:             args{containerName: "foo"},
			want:             FilePaths{"/var/log/nginx/access.log"},
		},
		{
			name:             "multiPath",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "foo:/var/log/nginx/access.log,foo:/var/log/nginx/error.log",
			args:             args{containerName: "foo"},
			want:             FilePaths{"/var/log/nginx/access.log", "/var/log/nginx/error.log"},
		},
		{
			name:             "multiKeyMultiPath",
			h:                NewHandler([]string{"foo", "bar"}),
			tailerAnnotation: "foo:/var/log/nginx/access.log,foo:/var/log/nginx/error.log, bar:/dev/zero, bar:/dev/urandom",
			args:             args{containerName: "bar"},
			want:             FilePaths{"/dev/zero", "/dev/urandom"},
		},
		{
			name:             "handlerWithEmptyArgsDefaultsToDefaultContainerName",
			h:                NewHandler([]string{"foo", "bar"}),
			tailerAnnotation: "foo:/var/log/nginx/access.log,foo:/var/log/nginx/error.log, bar:/dev/zero, bar:/dev/urandom",
			args:             args{containerName: ""},
			want:             FilePaths{"/var/log/nginx/access.log", "/var/log/nginx/error.log"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.h.AddTailerAnnotation(tt.tailerAnnotation)
			if got := tt.h.FilePathsForContainer(tt.args.containerName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\nFilePathsForContainer() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestAllFilePaths(t *testing.T) {
	tests := []struct {
		name             string
		h                *Handler
		tailerAnnotation TailerAnnotation
		want             FilePaths
	}{
		{
			name:             "nilHandler",
			h:                nil,
			tailerAnnotation: "",
			want:             nil,
		},
		{
			name:             "emptyHandler",
			h:                NewHandler([]string{}),
			tailerAnnotation: "",
			want:             FilePaths{},
		},
		{
			name:             "handlerWithEmptyPaths",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "",
			want:             FilePaths{},
		},
		{
			name:             "singleKeySinglePath",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "foo:/var/log/nginx/error.log",
			want:             FilePaths{"/var/log/nginx/error.log"},
		},
		{
			name:             "singleKeyMultiPath",
			h:                NewHandler([]string{"foo"}),
			tailerAnnotation: "foo:/var/log/nginx/error.log,foo:/var/log/nginx/access.log",
			want:             FilePaths{"/var/log/nginx/error.log", "/var/log/nginx/access.log"},
		},
		{
			name:             "multiKeyMultiPath",
			h:                NewHandler([]string{"foo", "bar"}),
			tailerAnnotation: "foo:/var/log/nginx/error.log,foo:/var/log/nginx/access.log, bar:/foo/bar, bar:/foo/bar/baz",
			want:             FilePaths{"/var/log/nginx/error.log", "/var/log/nginx/access.log", "/foo/bar", "/foo/bar/baz"},
		},
		{
			name:             "multiKeyMultiPathWithDuplications",
			h:                NewHandler([]string{"foo", "bar"}),
			tailerAnnotation: "foo:/var/log/nginx/error.log,foo:/var/log/nginx/access.log, bar:/foo/bar, bar:/var/log/nginx/error.log",
			want:             FilePaths{"/var/log/nginx/error.log", "/var/log/nginx/access.log", "/foo/bar", "/var/log/nginx/error.log"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.h.AddTailerAnnotation(tt.tailerAnnotation)
			if got := tt.h.AllFilePaths(); !matchingStringSlicesWithoutOrder(got, tt.want) {
				t.Errorf("\nAllFilePaths() = %v\nwant %v", got, tt.want)
			}
		})
	}
}

func matchingStringSlicesWithoutOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	diff := make(map[string]int, len(a))

	for _, s := range a {
		diff[s]++
	}

	for _, s := range b {
		if _, ok := diff[s]; !ok {
			return false
		}
		diff[s]--
		if diff[s] == 0 {
			delete(diff, s)
		}
	}

	return len(diff) == 0
}
