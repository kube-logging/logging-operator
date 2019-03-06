# Plugin parser
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| format | - |  |
| timeFormat | - |  |
| keyName | log |  |
## Plugin template
```
<filter {{ .pattern }}.** >
  @type parser
  format {{ .format }}
  time_format {{ .timeFormat }}
  key_name {{ .keyName }}
</filter>

```