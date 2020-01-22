### FluentdSpec
#### FluentdSpec defines the desired state of Fluentd

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| annotations | map[string]string | No | - |  |
| labels | map[string]string | No | - |  |
| tls | FluentdTLS | No | - |  |
| image | ImageSpec | No | - |  |
| disablePvc | bool | No | - |  |
| bufferStorageVolume | storage.KubernetesStorage | No | - | BufferStorageVolume is by default configured as PVC using FluentdPvcSpec<br> |
| volumeMountChmod | bool | No | - |  |
| volumeModImage | ImageSpec | No | - |  |
| configReloaderImage | ImageSpec | No | - |  |
| resources | corev1.ResourceRequirements | No | - |  |
| livenessProbe | *corev1.Probe | No | - |  |
| livenessDefaultCheck | bool | No | - |  |
| readinessProbe | *corev1.Probe | No | - |  |
| port | int32 | No | - |  |
| tolerations | []corev1.Toleration | No | - |  |
| nodeSelector | map[string]string | No | - |  |
| metrics | *Metrics | No | - |  |
| security | *Security | No | - |  |
| scaling | *FluentdScaling | No | - |  |
| logLevel | string | No | - |  |
| podPriorityClassName | string | No | - |  |
| fluentLogDestination | string | No | - |  |
| fluentOutLogrotate | *FluentOutLogrotate | No | - | FluentOutLogrotate sends fluent's stdout to file and rotates it<br> |
### FluentOutLogrotate
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| enabled | bool | Yes | - |  |
| path | string | No | - |  |
| age | string | No | - |  |
| size | string | No | - |  |
### FluentdScaling
#### FluentdScaling enables configuring the scaling behaviour of the fluentd statefulset

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| replicas | int | Yes | - |  |
### FluentdTLS
#### FluentdTLS defines the TLS configs

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| enabled | bool | Yes | - |  |
| secretName | string | Yes | - |  |
| sharedKey | string | No | - |  |
