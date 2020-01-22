// Copyright © 2020 Banzai Cloud
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

package docgen

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/go-logr/logr"
)

type SourceLister struct {
	Logger                       logr.Logger
	Sources                      map[string]SourceDir
	IgnoredSources               []string
	DefaultValueFromTagExtractor func(string) string
	Index                        *Doc
}

type DocIndex struct {
	Path string
}

type SourceDir struct {
	Path     string
	DestPath string
}

type Source struct {
	Item     DocItem
	Category string
}

func NewSourceLister(sources map[string]SourceDir, logger logr.Logger) *SourceLister {
	return &SourceLister{
		Logger:  logger,
		Sources: sources,
	}
}

func (sl *SourceLister) ListSources() ([]Source, error) {
	sourceList := []Source{}
	for category, p := range sl.Sources {
		files, err := ioutil.ReadDir(p.Path)
		if err != nil {
			return nil, errors.WrapIff(err, "failed to read files from %s", p.Path)
		}
		for _, file := range files {
			fname := strings.Replace(file.Name(), ".go", "", 1)
			if filepath.Ext(file.Name()) == ".go" && sl.IsWhiteListed(fname) {
				fullPath := filepath.Join(p.Path, file.Name())
				sourceList = append(sourceList, Source{
					Category: category,
					Item: DocItem{
						Name:                         fname,
						SourcePath:                   fullPath,
						DestPath:                     p.DestPath,
						DefaultValueFromTagExtractor: sl.DefaultValueFromTagExtractor,
					},
				})
			}
		}
	}

	return sourceList, nil
}

func (sl *SourceLister) IsWhiteListed(source string) bool {
	for _, p := range sl.IgnoredSources {
		r := regexp.MustCompile(p)
		if r.MatchString(source) {
			sl.Logger.V(2).Info("ignored source", "source", source)
			return false
		}
	}
	sl.Logger.V(2).Info("included source", "source", source)
	return true
}

func (lister *SourceLister) Generate(log logr.Logger) error {
	lister.Index.Append("<center>\n")
	lister.Index.Append("| Name | Type | Description | Status |Version |")
	lister.Index.Append("|:---|---|:---|:---:|---:|")

	sources, err := lister.ListSources()
	if err != nil {
		return errors.WrapIf(err, "failed to get plugin list")
	}

	for _, source := range sources {
		document := GetDocumentParser(source.Item, log.WithName("docgen"))
		if err := document.Generate(); err != nil {
			return err
		}

		if lister.Index != nil {
			relPath, err := filepath.Rel(lister.Index.Item.DestPath, document.Item.DestPath)
			if err != nil {
				return errors.WrapIff(err, "failed to determine relpath for %s", document.Item.DestPath)
			}

			lister.Index.Append(fmt.Sprintf("| **[%s](%s)** | %s | %s | %s | [%s](%s) |",
				document.DisplayName,
				filepath.Join(relPath, document.Item.Name+".md"),
				source.Category,
				document.Desc,
				document.Status,
				document.Version,
				document.Url))
		}

	}

	lister.Index.Append("</center>")

	if err := lister.Index.Generate(); err != nil {
		return err
	}

	return nil
}
