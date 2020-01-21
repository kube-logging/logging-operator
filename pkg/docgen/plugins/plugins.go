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

package plugins

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"emperror.dev/errors"
	"github.com/banzaicloud/logging-operator/pkg/docgen"
	"github.com/go-logr/logr"
)

type PluginLister struct {
	Logger         logr.Logger
	Sources        map[string]PluginDir
	IgnoredPlugins []string
}

type PluginDir struct {
	Path     string
	DestPath string
}

type Plugin struct {
	Item     docgen.DocItem
	Category string
}

func NewPluginLister(sources map[string]PluginDir, ignoredPlugins []string, logger logr.Logger) *PluginLister {
	return &PluginLister{
		Logger:         logger,
		Sources:        sources,
		IgnoredPlugins: ignoredPlugins,
	}
}

func (pd *PluginLister) GetPlugins() ([]Plugin, error) {
	pluginList := []Plugin{}
	for category, p := range pd.Sources {
		files, err := ioutil.ReadDir(p.Path)
		if err != nil {
			return nil, errors.WrapIff(err, "failed to read files from %s", p.Path)
		}
		for _, file := range files {
			pd.Logger.V(2).Info("fileListGenerator", "filename", "file")
			fname := strings.Replace(file.Name(), ".go", "", 1)
			if filepath.Ext(file.Name()) == ".go" && pd.getPluginWhiteList(fname) {
				fullPath := filepath.Join(p.Path, file.Name())
				pluginList = append(pluginList, Plugin{
					Category: category,
					Item: docgen.DocItem{
						Name: fname, SourcePath: fullPath, DestPath: p.DestPath},
				})
			}
		}
	}

	return pluginList, nil
}

func (pd PluginLister) getPluginWhiteList(pluginName string) bool {
	for _, p := range pd.IgnoredPlugins {
		r := regexp.MustCompile(p)
		if r.MatchString(pluginName) {
			pd.Logger.V(2).Info("fileListGenerator", "ignored plugin", pluginName)
			return false
		}
	}
	return true
}
