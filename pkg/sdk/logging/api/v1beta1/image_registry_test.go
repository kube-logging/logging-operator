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

package v1beta1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceImageRegistry(t *testing.T) {
	tests := []struct {
		registry string
		image    string
		expected string
	}{
		{"", "ghcr.io/fluent/fluent-bit", "ghcr.io/fluent/fluent-bit"},
		{"mirror.example.com", "ghcr.io/fluent/fluent-bit", "mirror.example.com/fluent/fluent-bit"},
		{"mirror.example.com/", "docker.io/library/busybox", "mirror.example.com/library/busybox"},
		{"mirror.example.com/ghcr", "ghcr.io/kube-logging/eventrouter:1.0.1", "mirror.example.com/ghcr/kube-logging/eventrouter:1.0.1"},
		{"localhost:5000", "registry.k8s.io/pause", "localhost:5000/pause"},
		{"mirror.example.com", "localhost/pause", "mirror.example.com/pause"},
		{"mirror.example.com", "fluent/fluent-bit", "mirror.example.com/fluent/fluent-bit"},
		{"mirror.example.com", "busybox", "mirror.example.com/busybox"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, ReplaceImageRegistry(test.registry, test.image), "registry %q, image %q", test.registry, test.image)
	}
}

func TestSetDefaultsWithDefaultImageRegistry(t *testing.T) {
	SetDefaultImageRegistry("mirror.example.com")
	t.Cleanup(func() { SetDefaultImageRegistry("") })

	fluentd := &FluentdSpec{Image: ImageSpec{Repository: "my.registry/custom/fluentd"}}
	require.NoError(t, fluentd.SetDefaults())
	assert.Equal(t, "my.registry/custom/fluentd", fluentd.Image.Repository)
	assert.Equal(t, "mirror.example.com/kube-logging/logging-operator/config-reloader", fluentd.ConfigReloaderImage.Repository)
	assert.Equal(t, "mirror.example.com/library/busybox", fluentd.VolumeModImage.Repository)

	syslogNG := &SyslogNGSpec{}
	syslogNG.SetDefaults()
	assert.Equal(t, "mirror.example.com/axoflow/axosyslog", syslogNG.SyslogNGImage.Repository)
}
