# [Fluentd](https://www.fluentd.org/) container images for the [Logging operator](https://github.com/kube-logging/logging-operator)

You can find [Fluentd](https://www.fluentd.org/) container images to be used with the [Logging operator](https://github.com/kube-logging/logging-operator) in here.

## Usage

Pick an image type (`filters` contains filter plugins only, `full` has output plugins as well).
Image tags are constructed according to the following pattern:

```sh
ghcr.io/kube-logging/logging-operator/fluentd:{LOGGING-OPERATOR-VERSION}-{IMAGE-TYPE}
```

The images are based on the Debian variant of the official Ruby image (`ruby:<version>-slim-<debian-release>`).

### Memory allocator

Fluentd runs with [jemalloc](https://jemalloc.net/) preloaded through `LD_PRELOAD`, as recommended upstream, to keep heap fragmentation of long-running processes low.
To fall back to the glibc allocator, set the `LD_PRELOAD` environment variable of the fluentd container to an empty string.

### Add new plugins

If you wish to add a new plugin, use this image as a base image in your `Dockerfile`:

```dockerfile
FROM ghcr.io/kube-logging/logging-operator/fluentd:{LOGGING-OPERATOR-VERSION}-{IMAGE-TYPE}
```

Then add your plugin:

```dockerfile
RUN fluent-gem install PLUGIN_NAME -v PLUGIN_VERSION
```

Plugins with native extensions need a compiler and the matching `-dev` packages, which are installed with `apt-get` as `root`:

```dockerfile
USER root
RUN apt-get update \
 && apt-get install -y --no-install-recommends build-essential \
 && fluent-gem install PLUGIN_NAME -v PLUGIN_VERSION \
 && apt-get purge -y --auto-remove build-essential \
 && rm -rf /var/lib/apt/lists/*
USER fluent
```

## Version Support Policy

According to the Logging Operators release-cycle (6 weeks) we maintain the corresponding fluentd image version, which we support for the last 3 releases.

## Maintenance

Whenever a new Fluentd version is released, check the supported versions and add/remove versions in this directory accordingly.

Based on the supported Fluentd versions, you may drop old versions from the repository.

The `Dockerfile` in this directory is not generated and it doesn't use build args to keep things simple.
We may revisit that decision in the future.
