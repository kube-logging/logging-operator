## Fluent Bit config reload

It is now possible to configure Fluent Bit to reload its configuration on the fly.

This behaviour is disabled by default, but can be enabled with a single configuration option under 
the Logging's `spec.fluentbit.configHotReload` (legacy method) or the new FluentbitAgent's `spec.configHotReload`:

```
apiVersion: logging.banzaicloud.io/v1beta1
kind: FluentbitAgent
metadata:
  name: reload-example
spec:
  configHotReload: {}
```

Currently `resources` and `image` is configurable:
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: FluentbitAgent
metadata:
  name: reload-example
spec:
  configHotReload:
    resources: ...
    image:
      repository: ghcr.io/kube-logging/config-reloader
      tag: v0.0.5
```

For all the available configuration options please check the [API docs](https://kube-logging.dev/docs/configuration/crds/v1beta1/fluentbit_types/)
