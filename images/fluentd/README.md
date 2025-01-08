# [Fluentd](https://www.fluentd.org/) container images for the [Logging operator](https://github.com/kube-logging/logging-operator)

This repository publishes [Fluentd](https://www.fluentd.org/) container images to be used with the [Logging operator](https://github.com/kube-logging/logging-operator).

## Usage

Pick Fluentd version (either full semver or a shortened major-minor version).
Pick an image type (`filters` contains filter plugins only, `full` has output plugins as well).
Image tags are constructed according to the following pattern:

```sh
ghcr.io/kube-logging/fluentd:VERSION-IMAGE-TYPE
```

To ensure that subsequent builds don't break your production environment,
you may want to pin your image to a specific build:

```sh
ghcr.io/kube-logging/fluentd:VERSION-IMAGE-TYPE-build.BUILD_NUMBER
```

While the tag in the first example is a moving tag (subsequent builds of the same versions produce the same tags),
build number annotated tags are immutable.

### Add new plugins

If you wish to add a new plugin, use this image as a base image in your `Dockerfile`:

```dockerfile
FROM ghcr.io/kube-logging/fluentd:VERSION-IMAGE-TYPE-build.BUILD_NUMBER
```

Then add your plugin:

```dockerfile
RUN fluent-gem install PLUGIN_NAME -v PLUGIN_VERSION
```

## Version Support Policy

According to the Logging Operators release-cycle (6 weeks) we maintain the corresponding fluentd image version, which we support for the last 3 releases.

## Maintenance

Whenever a new Fluentd version is released, check the supported versions and add/remove versions in this repository accordingly.

The versioned directories in the repository root are Fluentd versions.

When a new Fluentd minor version is released, copy the directory of an earlier version and update the version numbers in it.
Based on the supported Fluentd versions, you may drop old versions from the repository.

When adding and deleting versions from this repository, don't forget to update the build matrix in [.github/workflows/depdency-images.yaml](../../.github/workflows/dependency-images.yaml).

Dockerfiles in this repository are not generated and they don't use build args to keep things simple.
We may revisit that decision in the future.

Patch versions are automatically updated by renovate.
