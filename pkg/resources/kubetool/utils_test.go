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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixQualifiedNameIfInvalid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid name within limits",
			input:    "valid-name",
			expected: "valid-name",
		},
		{
			name:     "Valid name at max length",
			input:    "a-very-long-but-valid-name",
			expected: "a-very-long-but-valid-name",
		},
		{
			name:     "Name with invalid characters",
			input:    "^invalid$name%",
			expected: "invalid-name",
		},
		{
			name:     "Name with uppercase letters",
			input:    "UpperCaseName",
			expected: "uppercasename",
		},
		{
			name:     "Name with leading and trailing spaces",
			input:    "  spaced-name  ",
			expected: "spaced-name",
		},
		{
			name:     "Name too long should be truncated",
			input:    "a-really-long-name-that-exceeds-the-max-length-of-the-dns-label-standard-by-far",
			expected: "a-really-long-name-that-exceeds-the-max-l-", // Additionally some hash suffix will be added
		},
		{
			name:     "Name with consecutive invalid characters",
			input:    "invalid$$$name",
			expected: "invalid-name",
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(ttp.name, func(t *testing.T) {
			result := FixQualifiedNameIfInvalid(ttp.input)
			if len(ttp.input) > maxNameLength {
				assert.Equal(t, ttp.expected, result[:maxNameLength-hashSuffixLen])
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Name with special characters",
			input:    "invalid$name",
			expected: "invalid-name",
		},
		{
			name:     "Name with spaces and uppercases",
			input:    "  Upper Case Name  ",
			expected: "upper-case-name",
		},
		{
			name:     "Name with consecutive invalid characters",
			input:    "name$$with***symbols",
			expected: "name-with-symbols",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, sanitizeName(tt.input))
		})
	}
}
