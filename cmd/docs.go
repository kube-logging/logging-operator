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
	"fmt"

	"github.com/banzaicloud/logging-operator/pkg/docgen"
	"github.com/banzaicloud/logging-operator/pkg/docgen/plugins"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	verboseLogging := true
	ctrl.SetLogger(zap.Logger(verboseLogging))
	var log = ctrl.Log.WithName("docs").WithName("main")

	fileList, err := plugins.PluginDirs{
		Sources: []plugins.PluginDir{
			{"filters", "./pkg/sdk/model/filter/"},
			{"outputs", "./pkg/sdk/model/output/"},
			{"common", "./pkg/sdk/model/common/"},
		},
		IgnoredPluginsList: []string{
			"null",
			".*.deepcopy",
			".*_test",
		},
	}.GetPlugins()

	if err != nil {
		log.Error(err, "Directory check error.")
	}
	index := docgen.Doc{
		Name: "Readme",
	}
	index.Append("# Supported Plugins\n\n")
	index.Append("For more information please click on the plugin name")
	index.Append("<center>\n")
	index.Append("| Name | Type | Description | Status |Version |")
	index.Append("|:---|---|:---|:---:|---:|")

	for _, file := range fileList {
		log.Info("plugin", "Name", file.SourcePath)
		document := docgen.GetDocumentParser(file)
		document.Generate("docs/plugins")
		index.Append(fmt.Sprintf("| **[%s](%s)** | %s | %s | %s | [%s](%s) |",
			document.DisplayName,
			file.DocumentationPath,
			document.Type,
			document.Desc,
			document.Status,
			document.Version,
			document.Url))
	}
	index.Append("</center>")
	index.Generate("docs/plugins")
}
