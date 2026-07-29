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

package fixture

import "github.com/kube-logging/logging-operator/pkg/sdk/logging/api/v1beta1"

// The locally built images the suites load into each cluster. These are
// duplicated from common while both exist; images_test.go asserts they stay
// equal, so the two cannot drift the way the syslog-ng reloader name did.
const (
	FluentdImageRepo      = "fluentd-full"
	ConfigReloaderRepo    = "config-reloader"
	SyslogNGReloaderRepo  = "syslog-ng-reloader"
	FluentdDrainWatchRepo = "fluentd-drain-watch"
	NodeExporterRepo      = "node-exporter"

	localTag = "local"
)

func image(repo string) v1beta1.ImageSpec {
	return v1beta1.ImageSpec{Repository: repo, Tag: localTag}
}
