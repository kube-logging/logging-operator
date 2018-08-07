package plugins

// ParserFilter CRD name
const ParserFilter = "parser"

// ParserFilterDefaultValues for parser plugin
var ParserFilterDefaultValues = map[string]string{
	"keyName": "log",
}

// ParserFilterTemplate for parser plugin
const ParserFilterTemplate = `
<filter {{ .pattern }}.** >
  @type parser
  format {{ .format }}
  time_format {{ .timeFormat }}
  key_name {{ .keyName }}
</filter>
`
