# Developers documentation

THis documentation helps to set-up a developer environment and writing plugins for the operator.

## Setting up Kind

Install Kind on your computer
```
go get sigs.k8s.io/kind@v0.5.1
```

Create cluster
```
kind create cluster --name logging
```

Install prerequisites (this is a Kubebuilder makefile that will generate and install crds)
```
make install
```

Run the Operator
```
go run main.go
```

## Writing a plugin

To add a `plugin` to the logging operator you need to define the plugin struct.

> Note: Place your plugin in the corresponding directory `pkg/sdk/model/filter` or `pkg/sdk/model/output`

```go
type MyExampleOutput struct {
	// Path that is required for the plugin
	Path string `json:"path,omitempty"`
}
```

The plugin uses the **JSON** tags to parse and validate configuration. Without tags the configuration is not valid. The `fluent` parameter name must match with the JSON tag. Don't forget to use `omitempty` for non required parameters.

### Implement `ToDirective`

To render the configuration you have to implement the `ToDirective` function.
```go
func (c *S3OutputConfig) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	...
}
```
For simple Plugins you can use the `NewFlatDirective` function.
```go
func (c *ExampleOutput) ToDirective(secretLoader secret.SecretLoader) (types.Directive, error) {
	return types.NewFlatDirective(types.PluginMeta{
		Type:      "example",
		Directive: "output",
		Tags: "**",
	}, c, secretLoader)
}
```
For more example please check the available plugins.

### Reuse existing Plugin sections

You can embed existing configuration for your plugins. For example modern `Output` plugins have `Buffer` section.

```go
// +docLink:"Buffer,./buffer.md"
Buffer *Buffer `json:"buffer,omitempty"`
```

If you are using embedded sections you must call its `ToDirective` method manually and append it as a `SubDirective`

```go
if c.Buffer != nil {
	if buffer, err := c.Buffer.ToDirective(secretLoader); err != nil {
		return nil, err
	} else {
		s3.SubDirectives = append(s3.SubDirectives, buffer)
	}
}
```

### Special plugin tags
To document the plugins logging-operator uses the Go `tags` (like JSON tags). Logging operator uses `plugin` named tags for special instructions.

Special tag `default`
The default tag helps to give `default` values for parameters. These parameters are explicitly set in the generated fluentd configuration.
```go
RetryForever bool `json:"retry_forever" plugin:"default:true"`
```
Special tag `required`
The required tag ensures that the attribute can **not** be empty
```go
RetryForever bool `json:"retry_forever" plugin:"required"`
```

## Generate documentation for Plugin

The operator parse the `docstrings` for the documentation. 

```go
...
// AWS access key id
AwsAccessKey *secret.Secret `json:"aws_key_id,omitempty"`
...
```

Will generate the following Markdown

| Variable Name | Default | Applied function |
|---|---|---|
|AwsAccessKey| | AWS access key id|

You can *hint* default values in docstring via `(default: value)`. This is useful if you don't want to set default explicitly with `tag`. However during rendering defaults in `tags` have priority over docstring.
```go
...
// The format of S3 object keys (default: %{path}%{time_slice}_%{index}.%{file_extension})
S3ObjectKeyFormat string `json:"s3_object_key_format,omitempty"`
...
```

### Special docstrings

- `+docName:"Title for the plugin section"`
- `+docLink:"Buffer,./buffer.md"`

You can declare document **title** and **description** above the `type _doc* interface{}` variable declaration.

Example Document headings:
```go
// +docName:"Amazon S3 plugin for Fluentd"
// **s3** output plugin buffers event logs in local file and upload it to S3 periodically. This plugin splits files exactly by using the time of event logs (not the time when the logs are received). For example, a log '2011-01-02 message B' is reached, and then another log '2011-01-03 message B' is reached in this order, the former one is stored in "20110102.gz" file, and latter one in "20110103.gz" file.
type _docS3 interface{}
```

Example Plugin headings:
```go
// +kubebuilder:object:generate=true
// +docName:"Shared Credentials"
type S3SharedCredentials struct {
...
```

Example linking embedded sections
```go
// +docLink:"Buffer,./buffer.md"
Buffer *Buffer `json:"buffer,omitempty"`
```

### Generate docs for your Plugin

```
make docs
```
