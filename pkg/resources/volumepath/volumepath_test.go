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

package volumepath

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *List
	}{
		{
			name: "constructor",
			want: &List{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInit(t *testing.T) {
	type args struct {
		strings []string
	}
	tests := []struct {
		name string
		args args
		want *List
	}{
		{
			name: "nil",
			args: args{strings: nil},
			want: nil,
		},
		{
			name: "empty",
			args: args{strings: []string{}},
			want: Reference(List([]string{})),
		},
		{
			name: "singleItem",
			args: args{strings: []string{"foobar"}},
			want: Reference(List([]string{"foobar"})),
		},
		{
			name: "multiItem",
			args: args{strings: []string{"foo", "bar"}},
			want: Reference(List([]string{"foo", "bar"})),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := Init(tt.args.strings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniq(t *testing.T) {
	tests := []struct {
		name string
		s    *List
		want *List
	}{
		{
			name: "nil",
			s:    nil,
			want: nil,
		},
		{
			name: "empty",
			s:    New(),
			want: New(),
		},
		{
			name: "alreadyUniqeVolumes",
			s:    Init([]string{"/already", "/unique", "/path"}),
			want: Init([]string{"/already", "/unique", "/path"}),
		},
		{
			name: "nonUniqeVolumes",
			s:    Init([]string{"/not", "/not", "/unique", "/path", "/unique", "/volumes", "/volumes"}),
			want: Init([]string{"/not", "/unique", "/path", "/volumes"}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Uniq(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Uniq() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApply(t *testing.T) {
	type args struct {
		fn ApplyFn
	}
	tests := []struct {
		name string
		s    *List
		args args
		want *List
	}{
		{name: "AllNil",
			s:    nil,
			args: args{fn: nil},
			want: nil,
		},
		{name: "nilReceiver",
			s: nil,
			args: args{
				fn: ApplyFn(func(strs []string, idx int) *string {
					result := "foobar"
					return &result
				}),
			},
			want: nil,
		},
		{name: "nilFn",
			s:    Init([]string{"foo", "bar", "baz"}),
			args: args{fn: nil},
			want: Init([]string{"foo", "bar", "baz"}),
		},
		{name: "ModifyFn",
			s: Init([]string{"foo", "bar", "baz"}),
			args: args{
				fn: ApplyFn(func(strs []string, idx int) *string {
					str := strs[idx]
					str += "!"
					return &str
				}),
			},
			want: Init([]string{"foo!", "bar!", "baz!"}),
		},
		{name: "FilterByString",
			s: Init([]string{"foo", "bar", "baz"}),
			args: args{
				fn: ApplyFn(func(strs []string, idx int) *string {
					if strs[idx] == "foo" {
						return nil
					}
					return &strs[idx]
				}),
			},
			want: Init([]string{"bar", "baz"}),
		},
		{name: "FilterByIndex",
			s: Init([]string{"foo", "bar", "baz"}),
			args: args{
				fn: ApplyFn(func(strs []string, idx int) *string {
					if idx%2 == 0 {
						return &strs[idx]
					}
					return nil
				}),
			},
			want: Init([]string{"foo", "baz"}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Apply(tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirst(t *testing.T) {
	tests := []struct {
		name string
		s    *List
		want *string
	}{
		{
			name: "nil",
			s:    nil,
			want: nil,
		},
		{
			name: "empty",
			s:    New(),
			want: nil,
		},
		{
			name: "singleString",
			s:    Init([]string{"foo"}),
			want: StringReference("foo"),
		},
		{
			name: "multiStrings",
			s:    Init([]string{"foo", "bar", "baz"}),
			want: StringReference("foo"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		name string
		s    *List
		want *string
	}{
		{
			name: "nil",
			s:    nil,
			want: nil,
		},
		{
			name: "empty",
			s:    New(),
			want: nil,
		},
		{
			name: "singleString",
			s:    Init([]string{"foo"}),
			want: StringReference("foo"),
		},
		{
			name: "multiStrings",
			s:    Init([]string{"foo", "bar", "baz"}),
			want: StringReference("baz"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTopLevelPathList(t *testing.T) {
	tests := []struct {
		name string
		l    *List
		want *List
	}{
		{
			name: "nil",
			l:    nil,
			want: nil,
		},
		{
			name: "empty",
			l:    New(),
			want: New(),
		},
		{
			name: "singlePath",
			l:    Init([]string{"/var/log"}),
			want: Init([]string{"/var/log"}),
		},
		{
			name: "multiPath",
			l:    Init([]string{"/var/log", "/tmp/test/"}),
			want: Init([]string{"/var/log", "/tmp/test/"}),
		},
		{
			name: "singleNestedPath",
			l:    Init([]string{"/var/log", "/var/log/nginx"}),
			want: Init([]string{"/var/log", "/var/log"}),
		},
		{
			name: "multiNestedPath",
			l:    Init([]string{"/var/log", "/var/log/nginx", "/var"}),
			want: Init([]string{"/var", "/var", "/var"}),
		},
		{
			name: "mixedPath",
			l:    Init([]string{"/foo", "/var/log", "/bar/", "/var/log/nginx", "/var/bar", "/should/remain/intact", "/var/bar/baz"}),
			want: Init([]string{"/foo", "/var/log", "/bar/", "/var/log", "/var/bar", "/should/remain/intact", "/var/bar"}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.TopLevelPathList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.TopLevelPathList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertFilePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty",
			args: args{path: ""},
			want: "",
		},
		{
			name: "Root",
			args: args{path: "/"},
			want: "",
		},
		{
			name: "SingleDir",
			args: args{path: "/var"},
			want: "var",
		},
		{
			name: "MultiDir",
			args: args{path: "/var/log"},
			want: "var-log",
		},
		{
			name: "MultiDirClosingBackslash",
			args: args{path: "/var/log/"},
			want: "var-log",
		},
		{
			name: "ClosingBackslashMultipleTimes",
			args: args{path: "/var/log///"},
			want: "var-log",
		},
		{
			name: "OpeningBackslashMultipleTimes",
			args: args{path: "///var/log"},
			want: "var-log",
		},
		{
			name: "InvalidPath",
			args: args{path: "var"},
			want: "",
		},
		{
			name: "EscapeToDNS1123",
			args: args{path: "/var/log/nginx/error.log"},
			want: "var-log-nginx-error-log",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertFilePath(tt.args.path); got != tt.want {
				t.Errorf("ConvertFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveInvalidPath(t *testing.T) {
	type args struct {
		validatorFn ApplyFn
	}
	tests := []struct {
		name string
		l    *List
		args args
		want *List
	}{
		{
			name: "Allnil",
			l:    nil,
			args: args{validatorFn: nil},
			want: nil,
		},
		{
			name: "NilValidatorMustUseDefault",
			l:    Init([]string{"/foo", "bar", "/bar/baz"}),
			args: args{validatorFn: nil},
			want: Init([]string{"/foo", "/bar/baz"}),
		},
		{
			name: "CustomValidatorShouldApplied",
			l:    Init([]string{"/foo", "bar", "/bar/baz"}),
			args: args{validatorFn: ApplyFn(
				func(strs []string, idx int) *string {
					if len(strs[idx]) > 4 {
						return nil
					}
					return &strs[idx]
				},
			)},
			want: Init([]string{"/foo", "bar"}),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.RemoveInvalidPath(tt.args.validatorFn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List.RemoveInvalidPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
