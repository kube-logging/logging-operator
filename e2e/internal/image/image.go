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

// Package image is the only declaration of the images the Makefile builds. It
// sits below the harness rather than on it because common/setup owns the loader
// and the harness imports that, so only a package below both can serve both.
package image

import (
	"os"
	"strings"

	"github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"
)

const localTag = "local"

type Image struct {
	Env        string
	Repository string
	Tag        string
}

func Operator() Image { return resolve("LOGGING_OPERATOR_IMAGE", "controller") }

func Fluentd() Image { return resolve("FLUENTD_IMAGE", "fluentd-full") }

func ConfigReloader() Image { return resolve("CONFIG_RELOADER_IMAGE", "config-reloader") }

func SyslogNGReloader() Image { return resolve("SYSLOG_NG_RELOADER_IMAGE", "syslog-ng-reloader") }

func DrainWatch() Image { return resolve("FLUENTD_DRAIN_WATCH_IMAGE", "fluentd-drain-watch") }

func NodeExporter() Image { return resolve("NODE_EXPORTER_IMAGE", "node-exporter") }

func All() []Image {
	return []Image{
		Operator(),
		ConfigReloader(),
		SyslogNGReloader(),
		DrainWatch(),
		NodeExporter(),
		Fluentd(),
	}
}

// resolve puts the override in Repository and Tag so the loader and the specs
// name the same thing. One without a tag is ignored: it would push a ref no
// spec names.
func resolve(env, repository string) Image {
	img := Image{Env: env, Repository: repository, Tag: localTag}
	if repo, tag, ok := splitRef(os.Getenv(env)); ok {
		img.Repository, img.Tag = repo, tag
	}
	return img
}

// splitRef cuts at the last colon no slash follows, since a registry port
// carries one too.
func splitRef(ref string) (repository, tag string, ok bool) {
	i := strings.LastIndex(ref, ":")
	if i <= 0 || i == len(ref)-1 || strings.Contains(ref[i+1:], "/") {
		return "", "", false
	}
	return ref[:i], ref[i+1:], true
}

func (i Image) Ref() string { return i.Repository + ":" + i.Tag }

func (i Image) Spec() v1beta1.ImageSpec {
	return v1beta1.ImageSpec{Repository: i.Repository, Tag: i.Tag}
}

func (i Image) Basic() *v1beta1.BasicImageSpec {
	return &v1beta1.BasicImageSpec{Repository: i.Repository, Tag: i.Tag}
}
