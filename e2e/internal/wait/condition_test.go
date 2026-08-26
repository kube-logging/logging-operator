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

package wait

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The ptr.Deref default is chosen per direction: false where Active must be
// true, true where it must be false. Both leave an unset Active unsatisfied so
// the caller polls again. One backwards would either pass on an unreconciled
// config or never pass at all, so cover all three states.
func TestAttachedAndExcessAcrossEveryActiveState(t *testing.T) {
	yes, no := true, false

	for _, c := range []struct {
		name            string
		active          *bool
		attach, exclude bool
	}{
		{"active true", &yes, true, false},
		{"active false", &no, false, true},
		{"active unset", nil, false, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.attach, attached(nil, "lg", c.active, "lg"))
			assert.Equal(t, c.exclude, excess([]string{"boom"}, "", c.active))
		})
	}
}

// Each predicate needs all three of its fields, or a config halfway through
// reconciling reads as settled.
func TestAttachedAndExcessNeedEveryField(t *testing.T) {
	yes, no := true, false

	assert.False(t, attached([]string{"boom"}, "lg", &yes, "lg"), "problems reported")
	assert.False(t, attached(nil, "other", &yes, "lg"), "naming a different Logging")
	assert.True(t, attached(nil, "lg", &yes, "lg"))

	assert.False(t, excess(nil, "", &no), "no problems reported")
	assert.False(t, excess([]string{"boom"}, "lg", &no), "still naming a Logging")
	assert.True(t, excess([]string{"boom"}, "", &no))
}
