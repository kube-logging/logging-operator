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

package extensionsconfig

import (
	"reflect"
	"testing"
)

func TestFluentBitConfigFilePath(t *testing.T) {
	type args struct {
		image    string
		filePath string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "allEmpty",
			args: args{image: "", filePath: ""},
			want: []string{},
		},
		{
			name: "emptyFilePath",
			args: args{image: HostTailer.FluentBitImage, filePath: ""},
			want: []string{},
		},
		{
			name: "emptyImage",
			args: args{image: "", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/stdout"},
		},
		{
			name: "imageVersionInvalid",
			args: args{image: "ghcr.io/fluent/fluent-bit@x.Arg!unset", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/stdout"},
		},
		{
			name: "imageVersionShortVersionLess",
			args: args{image: "1.4.5", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/stdout"},
		},
		{
			name: "imageVersionLongVersionLess",
			args: args{image: "ghcr.io/fluent/fluent-bit:1.4.5", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/stdout"},
		},
		{
			name: "imageVersionShortVersionEquals",
			args: args{image: "1.4.6", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/", "-p", "file=stdout"},
		},
		{
			name: "imageVersionLongversionEquals",
			args: args{image: "ghcr.io/fluent/fluent-bit:1.4.6", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/", "-p", "file=stdout"},
		},
		{
			name: "imageVersionShortVersionGreater",
			args: args{image: "1.4.7", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/", "-p", "file=stdout"},
		},
		{
			name: "imageVersionLongversionGreater",
			args: args{image: "ghcr.io/fluent/fluent-bit:1.4.7", filePath: "/dev/stdout"},
			want: []string{"-p", "path=/dev/", "-p", "file=stdout"},
		},
		{
			name: "longPath",
			args: args{image: "ghcr.io/fluent/fluent-bit:1.4.7", filePath: "/var/log/nginx/myaccess.log"},
			want: []string{"-p", "path=/var/log/nginx/", "-p", "file=myaccess.log"},
		},
		{
			name: "relativeFilePath",
			args: args{image: "ghcr.io/fluent/fluent-bit:1.4.7", filePath: "./myaccess.log"},
			want: []string{"-p", "path=./", "-p", "file=myaccess.log"},
		},
		{
			name: "noPathGiven",
			args: args{image: "ghcr.io/fluent/fluent-bit:1.4.7", filePath: "myaccess.log"},
			want: []string{"-p", "path=myaccess.log"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := fluentBitConfigFilePath(tt.args.image, tt.args.filePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FluentBitConfigFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEventTailerUsesEnvConfig(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  bool
	}{
		{name: "default", image: EventTailer.ImageWithTag, want: true},
		{name: "breakingVersion", image: "ghcr.io/kube-logging/eventrouter:1.0.0", want: true},
		{name: "newerVersion", image: "ghcr.io/kube-logging/eventrouter:1.2.0", want: true},
		{name: "olderVersion", image: "ghcr.io/kube-logging/eventrouter:0.5.0", want: false},
		{name: "olderPatch", image: "ghcr.io/kube-logging/eventrouter:0.4.0", want: false},
		{name: "registryWithPort", image: "localhost:5000/eventrouter:0.5.0", want: false},
		{name: "latestTag", image: "ghcr.io/kube-logging/eventrouter:latest", want: true},
		{name: "noTag", image: "ghcr.io/kube-logging/eventrouter", want: true},
		{name: "digestPin", image: "ghcr.io/kube-logging/eventrouter@sha256:abc123", want: true},
		{name: "legacyTagWithDigest", image: "ghcr.io/kube-logging/eventrouter:0.5.0@sha256:abc123", want: false},
		{name: "modernTagWithDigest", image: "ghcr.io/kube-logging/eventrouter:1.0.0@sha256:abc123", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EventTailer.UsesEnvConfig(tt.image); got != tt.want {
				t.Errorf("UsesEnvConfig(%q) = %v, want %v", tt.image, got, tt.want)
			}
		})
	}
}
