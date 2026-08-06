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

package fixture

import "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

// Duplicated from common while the unmigrated suites still reference its copies,
// with TestImageNamesMatchCommon holding the two equal in the meantime. Both the
// duplicates and that test go when the last suite stops using common's, leaving
// this the only source.
const (
	FluentdImageRepo      = "fluentd-full"
	ConfigReloaderRepo    = "config-reloader"
	SyslogNGReloaderRepo  = "syslog-ng-reloader"
	FluentdDrainWatchRepo = "fluentd-drain-watch"
	NodeExporterRepo      = "node-exporter"

	localTag = "local"
)

// The images are the one part of a suite's spec that never varies: no suite
// overrides a repository or a tag.
func FluentdImage() v1beta1.ImageSpec { return image(FluentdImageRepo) }

func ConfigReloaderImage() v1beta1.ImageSpec { return image(ConfigReloaderRepo) }

func SyslogNGReloaderImage() v1beta1.ImageSpec { return image(SyslogNGReloaderRepo) }

func DrainWatchImage() v1beta1.ImageSpec { return image(FluentdDrainWatchRepo) }

func NodeExporterImage() v1beta1.ImageSpec { return image(NodeExporterRepo) }

// Basic is the shape the syslog-ng spec takes for the same images.
func Basic(spec v1beta1.ImageSpec) *v1beta1.BasicImageSpec {
	return &v1beta1.BasicImageSpec{Repository: spec.Repository, Tag: spec.Tag}
}

func image(repo string) v1beta1.ImageSpec {
	return v1beta1.ImageSpec{Repository: repo, Tag: localTag}
}
