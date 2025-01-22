# [Fluentd](https://www.fluentd.org/) container images for the [Logging operator](https://github.com/kube-logging/logging-operator)

You can find [Fluentd](https://www.fluentd.org/) container images to be used with the [Logging operator](https://github.com/kube-logging/logging-operator) in here.

## Usage

Pick an image type (`filters` contains filter plugins only, `full` has output plugins as well).
Image tags are constructed according to the following pattern:

```sh
ghcr.io/kube-logging/logging-operator/fluentd:{LOGGING-OPERATOR-VERSION}-{IMAGE-TYPE}
```

### Add new plugins

If you wish to add a new plugin, use this image as a base image in your `Dockerfile`:

```dockerfile
FROM ghcr.io/kube-logging/logging-operator/fluentd:{LOGGING-OPERATOR-VERSION}-{IMAGE-TYPE}
```

Then add your plugin:

```dockerfile
RUN fluent-gem install PLUGIN_NAME -v PLUGIN_VERSION
```

## Version Support Policy

According to the Logging Operators release-cycle (6 weeks) we maintain the corresponding fluentd image version, which we support for the last 3 releases.

## Maintenance

Whenever a new Fluentd version is released, check the supported versions and add/remove versions in this directory accordingly.

Based on the supported Fluentd versions, you may drop old versions from the repository.

The `Dockerfile` in this directory is not generated and it doesn't use build args to keep things simple.
We may revisit that decision in the future.
