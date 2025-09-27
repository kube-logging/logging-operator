package output

// +name:"RabbitMQ"
// +weight:"200"
type _hugoRabbitMQ interface{} //nolint:deadcode,unused

// +docName:"RabbitMQ plugin for Fluentd"
/*
Sends logs to RabbitMQ Queues. For details, see [https://github.com/nttcom/fluent-plugin-rabbitmq](https://github.com/nttcom/fluent-plugin-rabbitmq).

## Example output configurations

```yaml
spec:
  rabbitmq:
    host: rabbitmq-master.prod.svc.cluster.local
    buffer:
      tags: "[]"
      flush_interval: 10s
```
*/
type _docRabbitMQ interface{} //nolint:deadcode,unused

// +name:"RabbitMQ"
// +url:"https://github.com/nttcom/fluent-plugin-rabbitmq"
// +version:"0.3.5"
// +description:"Sends logs to RabbitMQ Queues."
// +status:"GA"
type _metaRabbitMQ interface{} //nolint:deadcode,unused

// +kubebuilder:object:generate=true
// +docName:"Output Config"
type RabbitMQOutputConfig struct {
}
