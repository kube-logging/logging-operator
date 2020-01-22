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

package main

import (
	"github.com/banzaicloud/logging-operator/pkg/docgen"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	verboseLogging := true
	rootLogger := zap.New(zap.UseDevMode(verboseLogging))

	lister := docgen.NewSourceLister(
		map[string]docgen.SourceDir{
			"filters": {Path: "pkg/sdk/model/filter", DestPath: "docs/plugins/filters"},
			"outputs": {Path: "pkg/sdk/model/output", DestPath: "docs/plugins/outputs"},
			"common":  {Path: "pkg/sdk/model/common", DestPath: "docs/plugins/common"},
		},
		rootLogger.WithName("pluginlister"))

	lister.IgnoredSources = []string{
		"null",
		".*.deepcopy",
		".*_test",
	}

	lister.DefaultValueFromTagExtractor = func(tag string) string {
		return docgen.GetPrefixedValue(tag, `plugin:\"default:(.*)\"`)
	}

	lister.Index = docgen.NewDoc(docgen.DocItem{
		Name:     "Readme",
		DestPath: "docs/plugins",
	}, rootLogger.WithName("plugins"))

	lister.Index.Append("# Supported Plugins\n\n")
	lister.Index.Append("For more information please click on the plugin name")

	if err := lister.Generate(rootLogger.WithName("plugins")); err != nil {
		panic(err)
	}
}
