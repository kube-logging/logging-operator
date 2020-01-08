# Scaling

In a large-scale infrastructure the logging components can get high load as well. The typical sign of this is when `fluentd` cannot handle its buffer directory size growth for more then the configured or calculated (timekey + timekey_wait) flush interval. In this case, you can [scale the fluentd statefulset](./crds.md#Scaling).