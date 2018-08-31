package plugins

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Plugin register map
var pluginRegister = map[string]Plugin{}

// Plugin struct to store plugin informations
type Plugin struct {
	Template      string
	DefaultValues map[string]string
}

// RegisterPlugin to use in CRD file
func RegisterPlugin(name string, template string, values map[string]string) {
	logrus.Infof("Registering plugin: %s", name)
	pluginRegister[name] = Plugin{Template: template, DefaultValues: values}
}

// GetDefaultValues get default values by name
func GetDefaultValues(name string) (map[string]string, error) {
	var err error
	value, ok := pluginRegister[name]
	if !ok {
		err = fmt.Errorf("plugin %q not found", name)
	}
	return value.DefaultValues, err
}

// GetTemplate get template string by name
func GetTemplate(name string) (string, error) {
	var err error
	value, ok := pluginRegister[name]
	if !ok {
		err = fmt.Errorf("plugin %q not found", name)
	}
	return value.Template, err
}

// Register plugins
func init() {
	RegisterPlugin(S3Output, S3Template, S3DefaultValues)
	RegisterPlugin(GCSOutput, GCSTemplate, GCSDefaultValues)
	RegisterPlugin(ParserFilter, ParserFilterTemplate, ParserFilterDefaultValues)
}
