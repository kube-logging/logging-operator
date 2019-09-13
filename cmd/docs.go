/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var log logr.Logger

type Doc struct {
	Name     string
	Content  string
	Type     string
	RootNode *ast.File
}

func (d *Doc) Append(line string) {
	d.Content = d.Content + line + "\n"
}

func (d *Doc) CheckNodes(n ast.Node) bool {
	generic, ok := n.(*ast.GenDecl)
	if ok {
		typeName, ok := generic.Specs[0].(*ast.TypeSpec)
		if ok {
			_, ok := typeName.Type.(*ast.InterfaceType)
			if ok && typeName.Name.Name == "_doc" {
				d.Append(fmt.Sprintf("# %s", getTypeName(generic, d.Name)))
				d.Append("## Overview")
				d.Append(getTypeDocs(generic))
				d.Append("## Configuration")
			}
			structure, ok := typeName.Type.(*ast.StructType)
			if ok {
				d.Append(fmt.Sprintf("### %s", getTypeName(generic, typeName.Name.Name)))
				if getTypeDocs(generic) != "" {
					d.Append(fmt.Sprintf("#### %s", getTypeDocs(generic)))
				}
				d.Append("| Variable Name | Type | Required | Default | Description |")
				d.Append("|---|---|---|---|---|")
				for _, item := range structure.Fields.List {
					name, com, def, required := getValuesFromItem(item)
					d.Append(fmt.Sprintf("| %s | %s | %s | %s | %s |", name, normaliseType(item.Type), required, def, com))
				}
			}

		}
	}

	return true
}

func normaliseType(fieldType ast.Expr) string {
	fset := token.NewFileSet()
	var typeNameBuf bytes.Buffer
	err := printer.Fprint(&typeNameBuf, fset, fieldType)
	if err != nil {
		log.Error(err, "error getting type")
	}
	return typeNameBuf.String()
}

func (d *Doc) Generate() {
	if d.RootNode != nil {
		ast.Inspect(d.RootNode, d.CheckNodes)
		log.Info("DocumentRoot not present skipping parse")
	}
	directory := fmt.Sprintf("./%s/%s/", DocsPath, d.Type)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		log.Error(err, "Md file create error %s", err.Error())
	}
	filepath := fmt.Sprintf("./%s/%s/%s.md", DocsPath, d.Type, d.Name)
	f, err := os.Create(filepath)
	if err != nil {
		log.Error(err, "Md file create error %s", err.Error())
	}
	defer closeFile(f)

	_, err = f.WriteString(d.Content)
	if err != nil {
		log.Error(err, "Md file write error %s", err.Error())
	}
}

var PluginDirs = map[string]string{
	"filters": "./pkg/model/filter/",
	"outputs": "./pkg/model/output/",
}

var DocsPath = "docs/plugins"

type Plugin struct {
	Name              string
	Type              string
	SourcePath        string
	DocumentationPath string
}

type Plugins []Plugin

var ignoredPluginsList = []string{
	"null",
	".*.deepcopy",
}

func main() {
	verboseLogging := true
	ctrl.SetLogger(zap.Logger(verboseLogging))
	log = ctrl.Log.WithName("docs").WithName("main")
	//log.Info("Plugin Directories:", "packageDir", packageDir)

	fileList, err := GetPlugins(PluginDirs)
	if err != nil {
		log.Error(err, "Directory check error.")
	}
	for _, file := range fileList {
		log.Info("Plugin", "Name", file.SourcePath)
		document := GetDocumentParser(file)
		document.Generate()
	}

	index := Doc{
		Name: "index",
	}
	index.Append("## Table of Contents\n\n")
	for pluginType := range PluginDirs {
		index.Append(fmt.Sprintf("### %s\n", pluginType))
		for _, plugin := range fileList {
			if plugin.Type == pluginType {
				index.Append(fmt.Sprintf("- [%s](%s)", plugin.Name, plugin.DocumentationPath))
			}
		}
		index.Append("\n")
	}

	index.Generate()

}

func getPrefixedLine(origin, expression string) string {
	r := regexp.MustCompile(expression)
	result := r.FindStringSubmatch(origin)
	if len(result) > 1 {
		return fmt.Sprintf("%s", result[1])
	}
	return ""
}

func getTypeName(generic *ast.GenDecl, defaultName string) string {
	structName := generic.Doc.Text()
	result := getPrefixedLine(structName, `\+docName:\"(.*)\"`)
	if result != "" {
		return result
	}
	return defaultName
}

func getTypeDocs(generic *ast.GenDecl) string {
	comment := ""
	if generic.Doc != nil {
		for _, line := range generic.Doc.List {
			newLine := strings.TrimPrefix(line.Text, "//")
			newLine = strings.TrimSpace(newLine)
			if !strings.HasPrefix(newLine, "+kubebuilder") &&
				!strings.HasPrefix(newLine, "+docName") {
				comment += newLine + "\n"
			}
		}
	}
	return comment
}

func getLink(def string) string {
	result := getPrefixedLine(def, `\+docLink:\"(.*)\"`)
	if result != "" {
		url := strings.Split(result, ",")
		def = strings.Replace(def, fmt.Sprintf("+docLink:\"%s\"", result), fmt.Sprintf("[%s](%s)", url[0], url[1]), 1)
	}
	return def
}

func formatRequired(r bool) string {
	if r {
		return "Yes"
	}
	return "No"
}

func getValuesFromItem(item *ast.Field) (name, comment, def, required string) {
	commentWithDefault := ""
	if item.Doc != nil {
		for _, line := range item.Doc.List {
			newLine := strings.TrimPrefix(line.Text, "//")
			newLine = strings.TrimSpace(newLine)
			if !strings.HasPrefix(newLine, "+kubebuilder") {
				commentWithDefault += newLine + "<br>"
			}
		}
	}
	tag := item.Tag.Value
	tagResult := getPrefixedLine(tag, `plugin:\"default:(.*)\"`)
	nameResult := getPrefixedLine(tag, `json:\"([^,\"]*).*\"`)
	required = formatRequired(!strings.Contains(getPrefixedLine(tag, `json:\"(.*)\"`), "omitempty"))
	if tagResult != "" {
		return nameResult, getLink(commentWithDefault), tagResult, required
	} else {
		result := getPrefixedLine(commentWithDefault, `\(default:(.*)\)`)
		if result != "" {
			ignore := fmt.Sprintf("(default:%s)", result)
			comment = strings.Replace(commentWithDefault, ignore, "", 1)
			return nameResult, comment, getLink(result), required
		}

		return nameResult, getLink(commentWithDefault), "-", required
	}
}

func GetDocumentParser(file Plugin) *Doc {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, file.SourcePath, nil, parser.ParseComments)
	if err != nil {
		log.Error(err, "Error!")
	}
	newDoc := &Doc{
		Name:     file.Name,
		RootNode: node,
		Type:     file.Type,
	}
	return newDoc
}

func GetPlugins(PluginDirs map[string]string) (Plugins, error) {
	var PluginList Plugins
	for pluginType, path := range PluginDirs {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Error(err, err.Error())
			return nil, err
		}

		for _, file := range files {
			log.V(2).Info("fileListGenerator", "filename", "file")
			fname := strings.Replace(file.Name(), ".go", "", 1)
			if filepath.Ext(file.Name()) == ".go" && getPluginWhiteList(fname) {
				fullPath := path + file.Name()
				filepath := fmt.Sprintf("./%s/%s.md", pluginType, fname)
				PluginList = append(PluginList, Plugin{Name: fname, SourcePath: fullPath, DocumentationPath: filepath, Type: pluginType})
			}
		}
	}

	return PluginList, nil
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Error(err, "File Close Error: %s", err.Error())
	}
}

func getPluginWhiteList(pluginName string) bool {
	for _, p := range ignoredPluginsList {
		r := regexp.MustCompile(p)
		if r.MatchString(pluginName) {
			log.Info("fileListGenerator", "ignored plugin", pluginName)
			return false
		}
	}
	return true
}
