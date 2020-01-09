# Scaling

In a large-scale infrastructure the logging components can get high load as well. The typical sign of this is when `fluentd` cannot handle its [buffer](./plugins/outputs/buffer.md) directory size growth for more than the configured or calculated (timekey + timekey_wait) flush interval. In this case, you can [scale the fluentd statefulset](./crds.md#Scaling).

>Note: When multiple instances send logs to the same output, the output can receive chunks of messages out of order. Some outputs tolerate this (for example, Elasticsearch), some do not, some require fine tuning (for example, Loki).