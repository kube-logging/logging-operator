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
	"fmt"
	"path/filepath"

	"emperror.dev/errors"
	"github.com/MakeNowJust/heredoc"
	"github.com/cisco-open/operator-tools/pkg/docgen"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var logger = zap.New(zap.UseDevMode(true))

func main() {
	plugins()
	crds()
}

func plugins() {
	lister := docgen.NewSourceLister(
		map[string]docgen.SourceDir{
			"filters":          {Path: "pkg/sdk/logging/model/filter", DestPath: "docs/configuration/plugins/filters"},
			"outputs":          {Path: "pkg/sdk/logging/model/output", DestPath: "docs/configuration/plugins/outputs"},
			"common":           {Path: "pkg/sdk/logging/model/common", DestPath: "docs/configuration/plugins/common"},
			"syslogng-outputs": {Path: "pkg/sdk/logging/model/syslogng/output", DestPath: "docs/configuration/plugins/syslogng-outputs"},
			"syslogng-filters": {Path: "pkg/sdk/logging/model/syslogng/filter", DestPath: "docs/configuration/plugins/syslogng-filters"},
		},
		logger.WithName("pluginlister"))

	lister.IgnoredSources = []string{
		"null",
		".*.deepcopy",
		".*_test",
	}

	lister.DefaultValueFromTagExtractor = func(tag string) string {
		return docgen.GetPrefixedValue(tag, `plugin:\"default:(.*)\"`)
	}

	lister.Index = docgen.NewDoc(docgen.DocItem{
		Name:     "_index",
		DestPath: "docs/configuration/plugins",
	}, logger.WithName("plugins"))

	lister.Header = heredoc.Doc(`
		---
		title: Supported Plugins
		generated_file: true
		---
		
		For more information please click on the plugin name
		<center>

		| Name | Profile | Description | Status |Version |
		|:---|---|:---|:---:|---:|`,
	)

	lister.Footer = heredoc.Doc(`
		</center>
	`)

	lister.DocGeneratedHook = func(document *docgen.Doc) error {
		relPath, err := filepath.Rel(lister.Index.Item.DestPath, document.Item.DestPath)
		if err != nil {
			return errors.WrapIff(err, "failed to determine relpath for %s", document.Item.DestPath)
		}

		lister.Index.Append(fmt.Sprintf("| **[%s](%s/)** | %s | %s | %s | [%s](%s) |",
			document.DisplayName,
			filepath.Join(relPath, document.Item.Name),
			document.Item.Category,
			document.Desc,
			document.Status,
			document.Version,
			document.Url))
		return nil
	}

	if err := lister.Generate(); err != nil {
		panic(err)
	}
}

func crds() {
	lister := docgen.NewSourceLister(
		map[string]docgen.SourceDir{
			"v1beta1":    {Path: "pkg/sdk/logging/api/v1beta1", DestPath: "docs/configuration/crds/v1beta1"},
			"extensions": {Path: "pkg/sdk/extensions/api/v1alpha1", DestPath: "docs/configuration/crds/extensions/v1alpha1"},
		},
		logger.WithName("crdlister"))

	lister.IgnoredSources = []string{
		".*.deepcopy",
		".*_test",
		".*_info",
	}

	lister.DefaultValueFromTagExtractor = func(tag string) string {
		return docgen.GetPrefixedValue(tag, `plugin:\"default:(.*)\"`)
	}

	lister.Index = docgen.NewDoc(docgen.DocItem{
		Name:     "_index",
		DestPath: "docs/configuration/crds/v1beta1",
	}, logger.WithName("crds"))

	lister.Header = heredoc.Doc(`
		---
		title: Available CRDs
		generated_file: true
		---
	
		For more information please click on the name
		<center>

		| Name | Description | Version |
		|---|---|---|`,
	)

	lister.Footer = heredoc.Doc(`
		</center>
	`)

	lister.DocGeneratedHook = func(document *docgen.Doc) error {
		relPath, err := filepath.Rel(lister.Index.Item.DestPath, document.Item.DestPath)
		if err != nil {
			return errors.WrapIff(err, "failed to determine relpath for %s", document.Item.DestPath)
		}
		lister.Index.Append(fmt.Sprintf("| **[%s](%s/)** | %s | %s |",
			document.DisplayName,
			filepath.Join(relPath, document.Item.Name),
			document.Desc,
			document.Item.Category))
		return nil
	}

	if err := lister.Generate(); err != nil {
		panic(err)
	}
}
