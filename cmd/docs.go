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

type doc struct {
	Name        string
	DisplayName string
	Content     string
	Type        string
	Version     string
	Url         string
	Desc        string
	Status      string

	RootNode *ast.File
}

func (d *doc) append(line string) {
	d.Content = d.Content + line + "\n"
}

func (d *doc) checkNodes(n ast.Node) bool {
	generic, ok := n.(*ast.GenDecl)
	if ok {
		typeName, ok := generic.Specs[0].(*ast.TypeSpec)
		if ok {
			_, ok := typeName.Type.(*ast.InterfaceType)
			if ok && strings.HasPrefix(typeName.Name.Name, "_doc") {
				d.append(fmt.Sprintf("# %s", getTypeName(generic, d.Name)))
				d.append("## Overview")
				d.append(getTypeDocs(generic, false))
				d.append("## Configuration")
			}
			if ok && strings.HasPrefix(typeName.Name.Name, "_meta") {
				d.DisplayName = getPrefixedLine(getTypeDocs(generic, true), `\+name:\"(.*)\"`)
				d.Url = getPrefixedLine(getTypeDocs(generic, true), `\+url:\"(.*)\"`)
				d.Version = getPrefixedLine(getTypeDocs(generic, true), `\+version:\"(.*)\"`)
				d.Desc = getPrefixedLine(getTypeDocs(generic, true), `\+description:\"(.*)\"`)
				d.Status = getPrefixedLine(getTypeDocs(generic, true), `\+status:\"(.*)\"`)
			}
			if ok && strings.HasPrefix(typeName.Name.Name, "_exp") {
				d.append(getTypeDocs(generic, false))
				d.append("---")
			}
			structure, ok := typeName.Type.(*ast.StructType)
			if ok {
				d.append(fmt.Sprintf("### %s", getTypeName(generic, typeName.Name.Name)))
				if getTypeDocs(generic, true) != "" {
					d.append(fmt.Sprintf("#### %s", getTypeDocs(generic, true)))
				}
				d.append("| Variable Name | Type | Required | Default | Description |")
				d.append("|---|---|---|---|---|")
				for _, item := range structure.Fields.List {
					name, com, def, required := getValuesFromItem(item)
					d.append(fmt.Sprintf("| %s | %s | %s | %s | %s |", name, normaliseType(item.Type), required, def, com))
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

func (d *doc) generate() {
	if d.RootNode != nil {
		ast.Inspect(d.RootNode, d.checkNodes)
		log.Info("DocumentRoot not present skipping parse")
	}
	directory := fmt.Sprintf("./%s/%s/", docsPath, d.Type)
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		log.Error(err, "Md file create error %s", err.Error())
	}
	filepath := fmt.Sprintf("./%s/%s/%s.md", docsPath, d.Type, d.Name)
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

type PluginDir struct {
	Type string
	Path string
}

var pluginDirs = []PluginDir{
	{"filters", "./pkg/sdk/model/filter/"},
	{"outputs", "./pkg/sdk/model/output/"},
	{"common", "./pkg/sdk/model/common/"},
}

var docsPath = "docs/plugins"

type plugin struct {
	Name              string
	Type              string
	SourcePath        string
	DocumentationPath string
}

type plugins []plugin

var ignoredPluginsList = []string{
	"null",
	".*.deepcopy",
	".*_test",
}

func main() {
	verboseLogging := true
	ctrl.SetLogger(zap.Logger(verboseLogging))
	log = ctrl.Log.WithName("docs").WithName("main")

	fileList, err := getPlugins(pluginDirs)
	if err != nil {
		log.Error(err, "Directory check error.")
	}
	index := doc{
		Name: "Readme",
	}
	index.append("# Supported Plugins\n\n")
	index.append("For more information please click on the plugin name")
	index.append("<center>\n")
	index.append("| Name | Type | Description | Status |Version |")
	index.append("|:---|---|:---|:---:|---:|")

	for _, file := range fileList {
		log.Info("plugin", "Name", file.SourcePath)
		document := getDocumentParser(file)
		document.generate()
		index.append(fmt.Sprintf("| **[%s](%s)** | %s | %s | %s | [%s](%s) |",
			document.DisplayName,
			file.DocumentationPath,
			document.Type,
			document.Desc,
			document.Status,
			document.Version,
			document.Url))
	}
	index.append("</center>")
	index.generate()
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

func getTypeDocs(generic *ast.GenDecl, trimSpace bool) string {
	comment := ""
	if generic.Doc != nil {
		for _, line := range generic.Doc.List {
			newLine := strings.TrimPrefix(line.Text, "//")
			if trimSpace {
				newLine = strings.TrimSpace(newLine)
			}
			if !strings.HasPrefix(strings.TrimSpace(newLine), "+kubebuilder") &&
				!strings.HasPrefix(strings.TrimSpace(newLine), "+docName") {
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
	}
	result := getPrefixedLine(commentWithDefault, `\(default:(.*)\)`)
	if result != "" {
		ignore := fmt.Sprintf("(default:%s)", result)
		comment = strings.Replace(commentWithDefault, ignore, "", 1)
		return nameResult, comment, getLink(result), required
	}

	return nameResult, getLink(commentWithDefault), "-", required
}

func getDocumentParser(file plugin) *doc {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, file.SourcePath, nil, parser.ParseComments)
	if err != nil {
		log.Error(err, "Error!")
	}
	newDoc := &doc{
		Name:     file.Name,
		RootNode: node,
		Type:     file.Type,
	}
	return newDoc
}

func getPlugins(PluginDirs []PluginDir) (plugins, error) {
	var PluginList plugins
	for _, p := range PluginDirs {
		files, err := ioutil.ReadDir(p.Path)
		if err != nil {
			log.Error(err, err.Error())
			return nil, err
		}
		for _, file := range files {
			log.V(2).Info("fileListGenerator", "filename", "file")
			fname := strings.Replace(file.Name(), ".go", "", 1)
			if filepath.Ext(file.Name()) == ".go" && getPluginWhiteList(fname) {
				fullPath := p.Path + file.Name()
				filepath := fmt.Sprintf("./%s/%s.md", p.Type, fname)
				PluginList = append(PluginList, plugin{
					Name: fname, SourcePath: fullPath, DocumentationPath: filepath, Type: p.Type})
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
