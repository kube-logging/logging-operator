# Plugin parser
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
| format | - |  |
| timeFormat | - |  |
| keyName | log |  |
| reserveData | true |  |
| removeKeyNameField | log |  |
## Plugin template
```
<filter {{ .pattern }}.** >
  @type parser
  format {{ .format }}
  time_format {{ .timeFormat }}
  key_name {{ .keyName }}
  reserve_data {{ .reserveData }}
  remove_key_name_field {{ .removeKeyNameField }}
</filter>

```