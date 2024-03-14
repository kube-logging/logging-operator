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

package model

import "testing"

func Test_hasIntersection(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "no intersection empty",
			args: args{
				a: []string{},
				b: []string{},
			},
			want: false,
		},
		{
			name: "no intersection nonempty",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"d", "e"},
			},
			want: false,
		},
		{
			name: "has intersection",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"b"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIntersection(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("hasIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
