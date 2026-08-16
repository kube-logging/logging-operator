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

package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

// defaults clears the overrides so an ambient value cannot decide the result.
func defaults(t *testing.T) {
	t.Helper()
	for _, i := range All() {
		t.Setenv(i.Env, "")
	}
}

func TestAllMatchesTheLoaderList(t *testing.T) {
	defaults(t)

	require.Equal(t, []Image{
		{Env: "LOGGING_OPERATOR_IMAGE", Repository: "controller", Tag: "local"},
		{Env: "CONFIG_RELOADER_IMAGE", Repository: "config-reloader", Tag: "local"},
		{Env: "SYSLOG_NG_RELOADER_IMAGE", Repository: "syslog-ng-reloader", Tag: "local"},
		{Env: "FLUENTD_DRAIN_WATCH_IMAGE", Repository: "fluentd-drain-watch", Tag: "local"},
		{Env: "NODE_EXPORTER_IMAGE", Repository: "node-exporter", Tag: "local"},
		{Env: "FLUENTD_IMAGE", Repository: "fluentd-full", Tag: "local"},
	}, All())
}

func TestSpecImagesAreLoaded(t *testing.T) {
	defaults(t)

	for _, i := range []Image{Fluentd(), ConfigReloader(), SyslogNGReloader(), DrainWatch(), NodeExporter()} {
		assert.Contains(t, All(), i)
	}
}

func TestRef(t *testing.T) {
	defaults(t)

	assert.Equal(t, "fluentd-full:local", Fluentd().Ref())
}

func TestSpecAndBasicAgree(t *testing.T) {
	defaults(t)

	i := SyslogNGReloader()
	assert.Equal(t, v1beta1.ImageSpec{Repository: "syslog-ng-reloader", Tag: "local"}, i.Spec())
	assert.Equal(t, &v1beta1.BasicImageSpec{Repository: "syslog-ng-reloader", Tag: "local"}, i.Basic())
}

// The override has to reach the specs, not just the loader: naming one ref and
// pushing another leaves the aggregator unable to pull.
func TestOverrideReachesTheSpecs(t *testing.T) {
	t.Run("repository and tag both follow it", func(t *testing.T) {
		defaults(t)
		t.Setenv("FLUENTD_IMAGE", "myrepo/fluentd:v2")

		i := Fluentd()

		assert.Equal(t, "myrepo/fluentd", i.Repository)
		assert.Equal(t, "v2", i.Tag)
		assert.Equal(t, "myrepo/fluentd:v2", i.Ref())
		assert.Equal(t, v1beta1.ImageSpec{Repository: "myrepo/fluentd", Tag: "v2"}, i.Spec())
	})

	t.Run("a registry port is not the tag separator", func(t *testing.T) {
		defaults(t)
		t.Setenv("FLUENTD_IMAGE", "registry:5000/fluentd:v2")

		assert.Equal(t, "registry:5000/fluentd:v2", Fluentd().Ref())
		assert.Equal(t, "v2", Fluentd().Tag)
	})

	t.Run("one without a tag is ignored", func(t *testing.T) {
		defaults(t)
		t.Setenv("FLUENTD_IMAGE", "registry:5000/fluentd")

		assert.Equal(t, "fluentd-full:local", Fluentd().Ref())
	})

	t.Run("the loader list carries it too", func(t *testing.T) {
		defaults(t)
		t.Setenv("NODE_EXPORTER_IMAGE", "ghcr.io/kube-logging/node-exporter:v1")

		assert.Contains(t, All(), Image{
			Env: "NODE_EXPORTER_IMAGE", Repository: "ghcr.io/kube-logging/node-exporter", Tag: "v1",
		})
	})
}
