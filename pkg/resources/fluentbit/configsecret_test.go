// Copyright © 2025 Kube logging authors
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

package fluentbit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
	"github.com/kube-logging/logging-operator/pkg/sdk/logging/model/types"
)

func TestInvalidFilterGrepConfig(t *testing.T) {
	invalidFilterGrep := &v1beta1.FilterGrep{
		Match:     "*",
		Regex:     []string{"regex", "reg2"},
		Exclude:   []string{"exclude"},
		LogicalOp: "AND",
	}

	_, err := toFluentdFilterGrep(invalidFilterGrep)

	assert.EqualError(t, err, "failed to parse grep filter for fluentbit, LogicalOp is set, it's not possible to set both Regex and Exclude")
}

func TestValidFilterGrepConfig(t *testing.T) {
	filterGrep := &v1beta1.FilterGrep{
		Match:     "*",
		Regex:     []string{"regex1", "regex2"},
		LogicalOp: "AND",
	}

	expectedFluentFilterGrep := &FluentdFilterGrep{
		Match:     "*",
		Regex:     []string{"regex1", "regex2"},
		LogicalOp: "AND",
	}

	parserFluentdFilterGrep, err := toFluentdFilterGrep(filterGrep)

	assert.NoError(t, err)
	assert.EqualValues(t, parserFluentdFilterGrep, expectedFluentFilterGrep)
}

func TestBufferStorageServiceSection(t *testing.T) {
	testCases := []struct {
		name          string
		bufferStorage v1beta1.BufferStorage
		expected      []string
		notExpected   []string
	}{
		{
			name: "flush_on_shutdown enabled",
			bufferStorage: v1beta1.BufferStorage{
				StoragePath:                   "/buffers",
				StorageBacklogFlushOnShutdown: "On",
			},
			expected: []string{
				"storage.path  /buffers",
				"storage.backlog.flush_on_shutdown  On",
			},
		},
		{
			name: "flush_on_shutdown disabled explicitly",
			bufferStorage: v1beta1.BufferStorage{
				StorageBacklogFlushOnShutdown: "Off",
			},
			expected: []string{"storage.backlog.flush_on_shutdown  Off"},
		},
		{
			name:          "flush_on_shutdown unset stays out of the config",
			bufferStorage: v1beta1.BufferStorage{StoragePath: "/buffers"},
			expected:      []string{"storage.path  /buffers"},
			notExpected:   []string{"storage.backlog.flush_on_shutdown"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mapped, err := types.NewStructToStringMapper(nil).StringsMap(testCase.bufferStorage)
			require.NoError(t, err)

			rendered, err := generateConfig(fluentBitConfig{BufferStorage: mapped})
			require.NoError(t, err)

			for _, expected := range testCase.expected {
				assert.Contains(t, rendered, expected)
			}
			for _, notExpected := range testCase.notExpected {
				assert.NotContains(t, rendered, notExpected)
			}
		})
	}
}
