# Scaling

In a large-scale infrastructure the logging components can get high load as well. The typical sign of this is when `fluentd` cannot handle its buffer directory size growth for more then the configured or calculated (timekey + timekey_wait) flush interval. In this case, you can [scale the fluentd statefulset](./crds.md#Scaling).

>Note: One caveat when scaling fluentd is that when there are more then one instances sending logs to the same output, then the output will more likely receive chunks of messages out of order. Some outputs tolerate this (e.g. Elasticsearch), however some do not, or requires fine tuning (e.g. Loki).