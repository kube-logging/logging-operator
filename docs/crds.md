# Custom Resource Definitions

This document contains the detailed information about the CRDs Logging operator uses.

Available CRDs:
- [loggings.logging.banzaicloud.io](/config/crd/bases/logging.banzaicloud.io_loggings.yaml)
- [outputs.logging.banzaicloud.io](/config/crd/bases/logging.banzaicloud.io_outputs.yaml)
- [flows.logging.banzaicloud.io](/config/crd/bases/logging.banzaicloud.io_flows.yaml)
- [clusteroutputs.logging.banzaicloud.io](/config/crd/bases/logging.banzaicloud.io_clusteroutputs.yaml)
- [clusterflows.logging.banzaicloud.io](/config/crd/bases/logging.banzaicloud.io_clusterflows.yaml)

> You can find example yamls  [here](/docs/examples)

## loggings

Logging resource define a logging infrastructure for your cluster. You can define **one** or **more** `logging` resource. This resource holds together a `logging pipeline`. It is responsible to deploy `fluentd` and `fluent-bit` on the cluster. It declares a `controlNamespace` and `watchNamespaces` if applicable.

> Note: The `logging` resources are referenced by `loggingRef`. If you setup multiple `logging flow` you have to reference other objects to this field. This can happen if you want to run multiple fluentd with separated configuration.

You can install `logging` resource via [Helm chart](/charts/logging-operator-logging) with built-in TLS generation.

### Namespace separation
A `logging pipeline` consist two type of resources.
- `Namespaced` resources: `Flow`, `Output`
- `Global` resources: `ClusterFlow`, `ClusterOutput`

The `namespaced` resources only effective in their **own** namespace. `Global` resources are operate **cluster wide**. 

> You can only create `ClusterFlow` and `ClusterOutput` in the `controlNamespace`. It **MUST** be a **protected** namespace that only **administrators** have access.

Create a namespace for logging
```bash
kubectl create ns logging
```

**`logging` plain example** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
  namespace: logging
spec:
  fluentd: {}
  fluentbit: {}
  controlNamespace: logging
```

**`logging` with filtered namespaces** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-namespaced
  namespace: logging
spec:
  fluentd: {}
  fluentbit: {}
  controlNamespace: logging
  watchNamespaces: ["prod", "test"]
```

### Logging parameters
| Name                    | Type           | Default | Description                                                             |
|-------------------------|----------------|---------|-------------------------------------------------------------------------|
| loggingRef              | string         | ""      | Reference name of the logging deployment                                |
| flowConfigCheckDisabled | bool           | False   | Disable configuration check before deploy                               |
| flowConfigOverride      | string         | ""      | Use static configuration instead of generated config.                   |  
| fluentbit               | [FluentbitSpec](#Fluent-bit-Spec) | {}      | Fluent-bit configurations                                               |
| fluentd                 | [FluentdSpec](#Fluentd-Spec)   | {}      | Fluentd configurations                                                  |
| watchNamespaces         | []string       | ""      | Limit namespaces from where to read Flow and Output specs               |
| controlNamespace        | string         | ""      | Control namespace that contains ClusterOutput and ClusterFlow resources |
| enableRecreateWorkloadOnImmutableFieldChange | bool | false | Recreate workloads that cannot be updated, see details below |

**enableRecreateWorkloadOnImmutableFieldChange**

Not all fields can be updated on Kubernetes objects. This is especially true for Statefulsets and Daemonsets.
In case there is a change that requires recreating the fluentd/fluentbit workloads use this field  
to move on but make sure to understand the consequences:
 - As of fluentd, to avoid data loss, make sure to use a persistent volume for buffers `logging.spec.fluentd.`, 
 which is the default, unless explicitly disabled or configured differently.
 - As of fluent-bit, to avoid duplicated logs, make sure to configure a hostPath volume for 
 the positions through `logging.spec.fluentbit.spec.positiondb`.

#### Fluentd Spec

You can customize the `fluentd` statefulset with the following parameters.

| Name                    | Type           | Default | Description                                                             |
|-------------------------|----------------|---------|-------------------------------------------------------------------------|
| annotations | map[string]string | {} | Extra annotations to Kubernetes resource|
| labels | map[string]string | {} | Extra labels for fluentd and it's related resources |
| tls | [TLS](#TLS-Spec) | {} | Configure TLS settings|
| image | [ImageSpec](#Image-Spec) | {} | Fluentd image override |
| fluentdPvcSpec | [PersistentVolumeClaimSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#persistentvolumeclaimspec-v1-core) | {} | Deprecated, use BufferStorageVolume |
| bufferStorageVolume | [KubernetesStorage](#KubernetesStorage) | nil | Fluentd PVC spec to mount persistent volume for Buffer |
| disablePvc | bool | false | Disable PVC binding |
| volumeModImage | [ImageSpec](#Image-Spec) | {} | Volume modifier image override |
| configReloaderImage | [ImageSpec](#Image-Spec) | {} | Config reloader image override |
| resources | [ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#resourcerequirements-v1-core) | {} | Resource requirements and limits |
| port | int | 24240 | Fluentd target port |
| tolerations | [Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#toleration-v1-core) | {} | Pod toleration |
| nodeSelector | [NodeSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#nodeselector-v1-core) | {} | A node selector represents the union of the results of one or more label queries over a set of nodes |
| metrics | [Metrics](./logging-operator-monitoring.md#metrics-variables) | {} | Metrics defines the service monitor endpoints |
| security | [Security](./security#security-variables) | {} | Security defines Fluentd, Fluentbit deployment security properties |
| podPriorityClassName | string | "" | Name of a priority class to launch fluentd with |
| scaling | [scaling](#scaling)] | "" | Fluentd scaling preferences |
| fluentLogDestination | string | "null" | Send internal fluentd logs to stdout, or use "null" to omit them, see: https://docs.fluentd.org/deployment/logging#capture-fluentd-logs |
| fluentOutLogrotate | [FluentOutLogrotate](#FluentOutLogrotate) | nil | Write to file instead of stdout and configure logrotate params. The operator configures it by default to write to /fluentd/log/out. https://docs.fluentd.org/deployment/logging#output-to-log-file |
| livenessProbe | [Probe](#Probe) | {} | Periodic probe of fluentd container liveness. Container will be restarted if the probe fails. |
| LivenessDefaultCheck | bool | false | Enable default liveness probe of fluentd container. |
| readinessProbe | [Probe](#Probe) | {} | Periodic probe of fluentd container service readiness. Container will be removed from service endpoints if the probe fails. |
| scaling | [Scaling](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#deploymentspec-v1-apps) | {replicas: 1} | Fluentd scaling configuration i.e replica count

**`logging` with custom pvc volume for buffers** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: 
    bufferStorageVolume:
      pvc:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 40Gi
          storageClassName: fast
          volumeMode: Filesystem
  fluentbit: {}
  controlNamespace: logging
```

**`logging` with custom hostPath volume for buffers** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: 
    disablePvc: true
    bufferStorageVolume:
      hostPath:
        path: "" # leave it empty to automatically generate: /opt/logging-operator/default-logging-simple/default-logging-simple-fluentd-buffer
  fluentbit: {}
  controlNamespace: logging
```

#### Fluent-bit Spec
| Name                    | Type           | Default | Description                                                             |
|-------------------------|----------------|---------|-------------------------------------------------------------------------|
| annotations | map[string]string | {} | Extra annotations to Kubernetes resource|
| labels | map[string]string | {} | Extra labels for fluent-bit and it's related resources |
| tls | [TLS](#TLS-Spec) | {} | Configure TLS settings|
| image | [ImageSpec](#Image-Spec) | {} | Fluentd image override |
| resources | [ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#resourcerequirements-v1-core) | {} | Resource requirements and limits |
| targetHost | string | *Fluentd host* | Hostname to send the logs forward |
| targetPort | int | *Fluentd port* |  Port to send the logs forward |
| parser | string | cri | Change fluent-bit input parse configuration. [Available parsers](https://github.com/fluent/fluent-bit/blob/master/conf/parsers.conf)  |
| tolerations | [Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#toleration-v1-core) | {} | Pod toleration |
| metrics | [Metrics](./logging-operator-monitoring.md#metrics-variables) | {} | Metrics defines the service monitor endpoints |
| security | [Security](./security#security-variables) | {} | Security defines Fluentd, Fluentbit deployment security properties |
| position_db | | | Deprecated, use positiondb instead |
| positiondb |  [KubernetesStorage](#KubernetesStorage) | nil | Add position db storage support. If nothing is configured an emptyDir volume will be used. |
| inputTail | [InputTail](./fluentbit.md#tail-inputtail) | {} | Preconfigured tailer for container logs on the host. Container runtime (containerd vs. docker) is automatically detected for convenience. |
| filterKubernetes | [FilterKubernetes](./fluentbit.md#kubernetes-filterkubernetes) | {} | Fluent Bit Kubernetes Filter allows to enrich your log files with Kubernetes metadata. |
| bufferStorage | [BufferStorage](./fluentbit.md#bufferstorage) |  | Buffer Storage configures persistent buffer to avoid losing data in case of a failure |
| bufferStorageVolume | [KubernetesStorage](#KubernetesStorage) | nil | Volume definition for the Buffer Storage. If nothing is configured an emptydir volume will be used. |
| customConfigSecret | string | "" | Custom secret to use as fluent-bit config.<br /> It must include all the config files necessary to run fluent-bit (_fluent-bit.conf_, _parsers*.conf_) |
| podPriorityClassName    | string         | ""      | Name of a priority class to launch fluentbit with                       |
| livenessProbe | [Probe](#Probe) | {} | Periodic probe of fluentbit container liveness. Container will be restarted if the probe fails. |
| readinessProbe | [Probe](#Probe) | {} | Periodic probe of fluentbit container service readiness. Container will be removed from service endpoints if the probe fails. |

**`logging` with custom fluent-bit annotations** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit:
    annotations:
      my-annotations/enable: true
  controlNamespace: logging
```

**`logging` with hostPath volumes for buffers and positions** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit:
    bufferStorageVolume:
      hostPath:
        path: "" # leave it empty to automatically generate
    positiondb:
      hostPath:
        path: "" # leave it empty to automatically generate
  controlNamespace: logging
```

#### Image Spec

Override default images

| Name                    | Type           | Default | Description |
|-------------------------|----------------|---------|-------------|
| repository | string | "" | Image repository |
| tag | string | "" | Image tag |
| pullPolicy | string | "" | Always, IfNotPresent, Never |

**`logging` with custom fluentd image** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: 
    image:
      repository: banzaicloud/fluentd
      tag: v1.7.4-alpine-12
      pullPolicy: IfNotPresent
  fluentbit: {}
  controlNamespace: logging
```

#### TLS Spec	

Define TLS certificate secret

| Name                    | Type           | Default | Description |
|-------------------------|----------------|---------|-------------|
| enabled | string | "" | Image repository |
| secretName | string | "" | Kubernetes secret that contains: **tls.crt, tls.key, ca.crt** |
| sharedKey | string | "" | Shared secret for fluentd authentication |


**`logging` setup with TLS**
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-tls
spec:
  fluentd:
    tls:
      enabled: true
      secretName: fluentd-tls
      sharedKey: asdadas
  fluentbit:
    tls:
      enabled: true
      secretName: fluentbit-tls
      sharedKey: asdadas
  controlNamespace: logging

```

#### KubernetesStorage

Define Kubernetes storage

| Name      | Type | Default | Description |
|-----------|------|---------|-------------|
| host_path | | | deprecated, use hostPath instead |
| hostPath | [HostPathVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#hostpathvolumesource-v1-core) | - | Represents a host path mapped into a pod. If path is empty, it will automatically be set to "/opt/logging-operator/<name of the logging CR>/<name of the volume>" |
| emptyDir | [EmptyDirVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#emptydirvolumesource-v1-core) | - | Represents an empty directory for a pod. |
| pvc | [PersistentVolumeClaim](#Persistent Volume Claim) | - | A PersistentVolumeClaim (PVC) is a request for storage by a user. |

#### Persistent Volume Claim

| Name      | Type | Default | Description |
|-----------|------|---------|-------------|
| spec | [PersistentVolumeClaimSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#persistentvolumeclaimspec-v1-core) | - | Spec defines the desired characteristics of a volume requested by a pod author. |
| source | [PersistentVolumeClaimVolumeSource](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.12/#persistentvolumeclaimvolumesource-v1-core) | - | PersistentVolumeClaimVolumeSource references the user's PVC in the same namespace.  |

The Persistent Volume Claim should be created with the given `spec` and with the `name` defined in the `source`'s `claimName`.

#### FluentOutLogrotate

Redirect fluentd's stdout to file and configure rotation settings.

This is important to avoid fluentd getting into a ripple effect when there is an error and the error message get's
back to the system as a log message, which generates another error, etc... 

Default settings configured by the operator
```
spec:
  fluentd:
    fluentOutLogrotate:
      enabled: true
      path: /fluentd/log/out
      age: 10
      size: 10485760
```

Disabling it and write to stdout (not recommended)
```
spec:
  fluentd:
    fluentOutLogrotate:
      enabled: false
```


#### Scaling

Scaling components 

| Name                    | Type           | Default | Description |
|-------------------------|----------------|---------|-------------|
| replicas | int | 1 | number of pod replicas |

**`logging` with custom fluentd replica number** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: 
    scaling:
      replicas: 3
  fluentbit: {}
  controlNamespace: logging
```

## Outputs, Clusteroutputs

Outputs are the final stage for a `logging flow`. You can define multiple `outputs` and attach them to multiple `flows`.

> Note: `Flow` can be connected to `Output` and `ClusterOutput` but `ClusterFlow` is only attachable to `ClusterOutput`.

### Defining outputs

The supported `Output` plugins are documented [here](./plugins/outputs)

| Name                    | Type              | Default | Description |
|-------------------------|-------------------|---------|-------------|
| **Output Definitions** | [Output](./plugins/outputs) | nil | Named output definitions |
| loggingRef | string | "" | Specified `logging` resource reference to connect `Output` and `ClusterOutput` to |


**`output` s3 example**
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: s3-output-sample
spec:
  s3:
    aws_key_id:
      valueFrom:
        secretKeyRef:
          name: s3-secret
          key: awsAccessKeyId
          namespace: default
    aws_sec_key:
      valueFrom:
        secretKeyRef:
          name: s3-secret
          key: awsSecretAccesKey
          namespace: default
    s3_bucket: example-logging-bucket
    s3_region: eu-west-1
    path: logs/${tag}/%Y/%m/%d/
    buffer:
      timekey: 1m
      timekey_wait: 10s
      timekey_use_utc: true
```

## flows, clusterflows

Flows define a `logging flow` that defines the `filters` and `outputs`.

> `Flow` resources are `namespaced`, the `selector` only select `Pod` logs within namespace.
> `ClusterFlow` select logs from **ALL** namespace.

### Parameters
| Name                    | Type              | Default | Description |
|-------------------------|-------------------|---------|-------------|
| selectors               | map[string]string | {}      | Kubernetes label selectors for the log. |
| filters                 | [][Filter](./plugins/filters)          | []      | List of applied [filter](./plugins/filters).  |
| loggingRef              | string | "" | Specified `logging` resource reference to connect `FLow` and `ClusterFlow` to |
| outputRefs              | []string | [] | List of [Outputs](#Defining-outputs) or [ClusterOutputs](#Defining-outputs) names |

*`flow` example with filters and output in the `default` namespace*
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: flow-sample
  namespace: default
spec:
  filters:
    - parser:
        remove_key_name_field: true
        parse:
          type: nginx
    - tag_normaliser:
        format: ${namespace_name}.${pod_name}.${container_name}
  outputRefs:
    - s3-output
  selectors:
    app: nginx
```


#### Probe
A Probe is a diagnostic performed periodically by the kubelet on a Container. To perform a diagnostic, the kubelet calls a Handler implemented by the Container. [More info](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)


| Name                    | Type           | Default | Description |
|-------------------------|----------------|---------|-------------|
| initialDelaySeconds | int | 0 | Number of seconds after the container has started before liveness probes are initiated. |
| timeoutSeconds | int | 1 | Number of seconds after which the probe times out. |
| periodSeconds | int | 10 | How often (in seconds) to perform the probe. |
| successThreshold | int | 1 | Minimum consecutive successes for the probe to be considered successful after having failed. |
| failureThreshold | int | 3 |  Minimum consecutive failures for the probe to be considered failed after having succeeded. |
| exec | array | {} |  Exec specifies the action to take. [More info](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#execaction-v1-core) |
| httpGet | array | {} |  HTTPGet specifies the http request to perform. [More info](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#httpgetaction-v1-core) |
| tcpSocket | array | {} |  TCPSocket specifies an action involving a TCP port. [More info](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#tcpsocketaction-v1-core) |

**`logging` with custom liveness config** 
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd:
    livenessProbe:
      periodSeconds: 60
      initialDelaySeconds: 600
      exec:
        command: 
        - "/bin/sh"
        - "-c"
        - >
          LIVENESS_THRESHOLD_SECONDS=${LIVENESS_THRESHOLD_SECONDS:-300};
          if [ ! -e /buffers ];
          then
            exit 1;
          fi;
          touch -d "${LIVENESS_THRESHOLD_SECONDS} seconds ago" /tmp/marker-liveness;
          if [ -z "$(find /buffers -type d -newer /tmp/marker-liveness -print -quit)" ];
          then
            exit 1;
          fi;
  fluentbit: {}
  controlNamespace: logging
```
