package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"text/template"
	"text/template/parse"

	"github.com/Masterminds/sprig"
	"github.com/banzaicloud/logging-operator/pkg/resources/plugins"
)

//TODO handle parameters
func main() {
	pluginMap := plugins.GetAll()
	var indexPage bytes.Buffer
	indexPage.WriteString("# List of ")
	for name, plugin := range pluginMap {
		var data bytes.Buffer
		data.WriteString(fmt.Sprintf("# Plugin %s\n", name))
		t := template.New("PluginTemplate").Funcs(sprig.TxtFuncMap())
		t, err := t.Parse(plugin.Template)
		if err != nil {
			panic(err)
		}
		data.WriteString("## Variables\n")
		data.WriteString("| Variable name | Default | Applied function |\n")
		data.WriteString(fmt.Sprintf("|---|---|---|\n"))
		for _, item := range listTemplateFields(t) {
			regExp, err := regexp.Compile(`{{(?P<Function>\w*)?\s*.(?P<Variable>.*)}}`)
			if err != nil {
				panic(err)
			}
			matches := regExp.FindStringSubmatch(item)
			vairableName := matches[2]
			variableFunc := matches[1]
			defaultValue, ok := plugin.DefaultValues[matches[2]]
			if !ok {
				defaultValue = "-"
			}
			data.WriteString(fmt.Sprintf("| %s | %s | %s |\n", vairableName, defaultValue, variableFunc))

		}
		data.WriteString("## Plugin template\n")
		data.WriteString("```" + plugin.Template + "\n```")
		err = ioutil.WriteFile("docs/plugins/"+name+".md", data.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func listTemplateFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}
	return res
}
