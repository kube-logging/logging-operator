<p align="center">
  <a href="https://github.com/vlzemtsov/config-reloader/">
    <img src="https://img.shields.io/badge/license-Apache%20v2-green.svg" alt="license">
  </a>
</p>

This directory is a fork of @zvlb's repository you can find the origin here -> <https://github.com/zvlb/config-reloader/>

# Kubernetes Config (ConfigMap and Secret) Reloader

This progect - based on <https://github.com/jimmidyson/configmap-reload> and <https://github.com/prometheus-operator/prometheus-operator/pkgs/container/prometheus-config-reloader>

**config-reloader** is a simple binary to trigger a reload when Kubernetes ConfigMaps or Secrets are updated.
It watches mounted volume dirs and notifies the target process changed files on dirs.
If changes exist - send webhook.

## Features

- Send webook if files in dirs changed (if configmap or secret have been changed)
- Control many dirs
- Unarchive .gz archive to file and update file, if .gz has been changed
- Init mode (stop after unarchive)
- Prometheus metrics

It is available as a Docker image at ghcr.io/banzaicloud/config-reloader:latest

### Usage

```sh
Usage of ./out/config-reloader:
  -dir-for-unarchive string
        Directory where the archives will be unpacked (default "/tmp/unarchive")
  -init-mode
        Init mode for unarchive files. Works only if volume-dir-archive exist. Default - false
  -volume-dir value
        the config map volume directory to watch for updates; may be used multiple times
  -volume-dir-archive value
        the config map volume directory to watch for updates and unarchiving; may be used multiple times
  -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9533")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
  -webhook-method string
        the HTTP method url to use to send the webhook (default "POST")
  -webhook-retries int
        the amount of times to retry the webhook reload request (default 1)
  -webhook-status-code int
        the HTTP status code indicating successful triggering of reload (default 200)
  -webhook-url value
        the url to send a request to when the specified config map volume directory has been updated
```
