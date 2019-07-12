# Plugin stdout
## Variables
| Variable name | Default | Applied function |
|---|---|---|
| pattern | - |  |
## Plugin template
```
<match {{ .pattern }}.** >
  @type stdout
</match>
```