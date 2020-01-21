// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.Apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/banzaicloud/logging-operator/pkg/docgen/plugins"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	verboseLogging := true
	rootLogger := zap.New(zap.UseDevMode(verboseLogging))

	lister := plugins.NewPluginLister(
		map[string]plugins.PluginDir{
			"filters": {"pkg/sdk/model/filter", "docs/plugins/filters"},
			"outputs": {"pkg/sdk/model/output", "docs/plugins/outputs"},
			"common":  {"pkg/sdk/model/common", "docs/plugins/common"},
		},
		[]string{
			"null",
			".*.deepcopy",
			".*_test",
		},
		rootLogger.WithName("pluginlister"))

	err := plugins.GenerateWithIndex(lister, rootLogger.WithName("plugins"))
	if err != nil {
		panic(err)
	}
}
