### ImageSpec
#### ImageSpec struct hold information about image specification

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| repository | string | No | - |  |
| tag | string | No | - |  |
| pullPolicy | string | No | - |  |
### Metrics
#### Metrics defines the service monitor endpoints

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| interval | string | No | - |  |
| timeout | string | No | - |  |
| port | int32 | No | - |  |
| path | string | No | - |  |
| serviceMonitor | bool | No | - |  |
| prometheusAnnotations | bool | No | - |  |
### Security
#### Security defines Fluentd, Fluentbit deployment security properties

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| serviceAccount | string | No | - |  |
| roleBasedAccessControlCreate | *bool | No | - |  |
| podSecurityPolicyCreate | bool | No | - |  |
| securityContext | *corev1.SecurityContext | No | - |  |
| podSecurityContext | *corev1.PodSecurityContext | No | - |  |
