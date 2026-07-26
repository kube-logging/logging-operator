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

// Every image in SyslogNGSpec has to default its repository and its tag
// independently, so that overriding only the repository - to point at a
// registry mirror, for example - keeps the operator's own default tag.
func TestSyslogNGSpecSetDefaultsImages(t *testing.T) {
	images := []struct {
		name        string
		get         func(*SyslogNGSpec) *BasicImageSpec
		set         func(*SyslogNGSpec, *BasicImageSpec)
		defaultRepo string
		defaultTag  string
	}{
		{
			name:        "syslogNGImage",
			get:         func(s *SyslogNGSpec) *BasicImageSpec { return s.SyslogNGImage },
			set:         func(s *SyslogNGSpec, i *BasicImageSpec) { s.SyslogNGImage = i },
			defaultRepo: defaultSyslogngImageRepository,
			defaultTag:  defaultSyslogngImageTag,
		},
		{
			name:        "configReloadImage",
			get:         func(s *SyslogNGSpec) *BasicImageSpec { return s.ConfigReloadImage },
			set:         func(s *SyslogNGSpec, i *BasicImageSpec) { s.ConfigReloadImage = i },
			defaultRepo: defaultConfigReloaderImageRepository,
			defaultTag:  defaultConfigReloaderImageTag,
		},
		{
			name:        "metricsExporterImage",
			get:         func(s *SyslogNGSpec) *BasicImageSpec { return s.MetricsExporterImage },
			set:         func(s *SyslogNGSpec, i *BasicImageSpec) { s.MetricsExporterImage = i },
			defaultRepo: defaultPrometheusExporterImageRepository,
			defaultTag:  defaultPrometheusExporterImageTag,
		},
		{
			name:        "bufferVolumeMetricsImage",
			get:         func(s *SyslogNGSpec) *BasicImageSpec { return s.BufferVolumeMetricsImage },
			set:         func(s *SyslogNGSpec, i *BasicImageSpec) { s.BufferVolumeMetricsImage = i },
			defaultRepo: defaultBufferVolumeImageRepository,
			defaultTag:  defaultBufferVolumeImageTag,
		},
	}

	// Version is baked in at build time and overrides the fallback tag of the
	// operator's own images. Pin it so the expectations below stay stable.
	originalVersion := Version
	Version = ""
	t.Cleanup(func() { Version = originalVersion })

	const (
		customRepo = "registry.internal/mirror/some-image"
		customTag  = "1.2.3-custom"
	)

	for _, image := range images {
		t.Run(image.name, func(t *testing.T) {
			t.Run("unset defaults both fields", func(t *testing.T) {
				spec := &SyslogNGSpec{}
				spec.SetDefaults()

				result := image.get(spec)
				require.NotNil(t, result)
				assert.Equal(t, image.defaultRepo, result.Repository)
				assert.Equal(t, image.defaultTag, result.Tag)
			})

			t.Run("repository override keeps the default tag", func(t *testing.T) {
				spec := &SyslogNGSpec{}
				image.set(spec, &BasicImageSpec{Repository: customRepo})
				spec.SetDefaults()

				result := image.get(spec)
				require.NotNil(t, result)
				assert.Equal(t, customRepo, result.Repository)
				assert.Equal(t, image.defaultTag, result.Tag)
				assert.Equal(t, customRepo+":"+image.defaultTag, result.RepositoryWithTag())
			})

			t.Run("tag override keeps the default repository", func(t *testing.T) {
				spec := &SyslogNGSpec{}
				image.set(spec, &BasicImageSpec{Tag: customTag})
				spec.SetDefaults()

				result := image.get(spec)
				require.NotNil(t, result)
				assert.Equal(t, image.defaultRepo, result.Repository)
				assert.Equal(t, customTag, result.Tag)
			})

			t.Run("both overrides are preserved", func(t *testing.T) {
				spec := &SyslogNGSpec{}
				image.set(spec, &BasicImageSpec{Repository: customRepo, Tag: customTag})
				spec.SetDefaults()

				result := image.get(spec)
				require.NotNil(t, result)
				assert.Equal(t, customRepo, result.Repository)
				assert.Equal(t, customTag, result.Tag)
			})
		})
	}
}

// The operator's own images follow the build-time Version when the user does
// not pin a tag, while the third-party images keep their vendored defaults.
func TestSyslogNGSpecSetDefaultsImagesFollowVersion(t *testing.T) {
	originalVersion := Version
	Version = "9.9.9"
	t.Cleanup(func() { Version = originalVersion })

	spec := &SyslogNGSpec{
		ConfigReloadImage:        &BasicImageSpec{Repository: "registry.internal/mirror/reloader"},
		BufferVolumeMetricsImage: &BasicImageSpec{Repository: "registry.internal/mirror/node-exporter"},
		SyslogNGImage:            &BasicImageSpec{Repository: "registry.internal/mirror/axosyslog"},
		MetricsExporterImage:     &BasicImageSpec{Repository: "registry.internal/mirror/metrics-exporter"},
	}
	spec.SetDefaults()

	assert.Equal(t, "registry.internal/mirror/reloader:9.9.9", spec.ConfigReloadImage.RepositoryWithTag())
	assert.Equal(t, "registry.internal/mirror/node-exporter:9.9.9", spec.BufferVolumeMetricsImage.RepositoryWithTag())
	assert.Equal(t, "registry.internal/mirror/axosyslog:"+defaultSyslogngImageTag, spec.SyslogNGImage.RepositoryWithTag())
	assert.Equal(t, "registry.internal/mirror/metrics-exporter:"+defaultPrometheusExporterImageTag, spec.MetricsExporterImage.RepositoryWithTag())
}
