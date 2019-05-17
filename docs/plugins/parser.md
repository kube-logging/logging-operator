# Plugin parser
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| format | - |  |
| timeFormat | - |  |
| keyName | log |  |
| reserveData | false |  |
## Plugin template
```
<filter {{ .pattern }}.** >
  @type parser
  format {{ .format }}
  time_format {{ .timeFormat }}
  key_name {{ .keyName }}
  reserve_data {{ .reserveData }}
</filter>

```