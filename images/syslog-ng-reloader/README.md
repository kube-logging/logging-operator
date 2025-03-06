# SyslogNG reload image

**SyslogNG reload** is a simple binary to trigger a reload when Kubernetes ConfigMaps or Secrets are updated.
It watches mounted volume dirs and notifies the target process changed files on dirs.
If changes exist - send webhook.

It is available as a Docker image at `ghcr.io/kube-logging/logging-operator/syslog-ng-reloader`

## License

The project is licensed under the [Apache License, Version 2.0](LICENSE).
