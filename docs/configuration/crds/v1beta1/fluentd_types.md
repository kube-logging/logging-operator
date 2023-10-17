---
title: FluentdSpec
weight: 200
generated_file: true
---

## FluentdSpec

FluentdSpec defines the desired state of Fluentd

### affinity (*corev1.Affinity, optional) {#fluentdspec-affinity}

Default: -

### annotations (map[string]string, optional) {#fluentdspec-annotations}

Default: -

### bufferStorageVolume (volume.KubernetesVolume, optional) {#fluentdspec-bufferstoragevolume}

BufferStorageVolume is by default configured as PVC using FluentdPvcSpec [volume.KubernetesVolume](https://github.com/cisco-open/operator-tools/tree/master/docs/types) 

Default: -

### bufferVolumeArgs ([]string, optional) {#fluentdspec-buffervolumeargs}

Default: -

### bufferVolumeImage (ImageSpec, optional) {#fluentdspec-buffervolumeimage}

Default: -

### bufferVolumeMetrics (*Metrics, optional) {#fluentdspec-buffervolumemetrics}

Default: -

### bufferVolumeResources (corev1.ResourceRequirements, optional) {#fluentdspec-buffervolumeresources}

Default: -

### compressConfigFile (bool, optional) {#fluentdspec-compressconfigfile}

Default: -

### configCheckAnnotations (map[string]string, optional) {#fluentdspec-configcheckannotations}

Default: -

### configCheckResources (corev1.ResourceRequirements, optional) {#fluentdspec-configcheckresources}

Default: -

### configReloaderImage (ImageSpec, optional) {#fluentdspec-configreloaderimage}

Default: -

### configReloaderResources (corev1.ResourceRequirements, optional) {#fluentdspec-configreloaderresources}

Default: -

### dnsConfig (*corev1.PodDNSConfig, optional) {#fluentdspec-dnsconfig}

Default: -

### dnsPolicy (corev1.DNSPolicy, optional) {#fluentdspec-dnspolicy}

Default: -

### disablePvc (bool, optional) {#fluentdspec-disablepvc}

Default: -

### enableMsgpackTimeSupport (bool, optional) {#fluentdspec-enablemsgpacktimesupport}

Allows Time object in buffer's MessagePack serde [more info]( https://docs.fluentd.org/deployment/system-config#enable_msgpack_time_support) 

Default: -

### envVars ([]corev1.EnvVar, optional) {#fluentdspec-envvars}

Default: -

### extraArgs ([]string, optional) {#fluentdspec-extraargs}

Default: -

### extraVolumes ([]ExtraVolume, optional) {#fluentdspec-extravolumes}

Default: -

### fluentLogDestination (string, optional) {#fluentdspec-fluentlogdestination}

Default: -

### fluentOutLogrotate (*FluentOutLogrotate, optional) {#fluentdspec-fluentoutlogrotate}

FluentOutLogrotate sends fluent's stdout to file and rotates it 

Default: -

### fluentdPvcSpec (*volume.KubernetesVolume, optional) {#fluentdspec-fluentdpvcspec}

Deprecated, use bufferStorageVolume 

Default: -

### forwardInputConfig (*input.ForwardInputConfig, optional) {#fluentdspec-forwardinputconfig}

Default: -

### ignoreRepeatedLogInterval (string, optional) {#fluentdspec-ignorerepeatedloginterval}

Ignore repeated log lines [more info]( https://docs.fluentd.org/deployment/logging#ignore_repeated_log_interval) 

Default: -

### ignoreSameLogInterval (string, optional) {#fluentdspec-ignoresameloginterval}

Ignore same log lines [more info]( https://docs.fluentd.org/deployment/logging#ignore_same_log_interval) 

Default: -

### image (ImageSpec, optional) {#fluentdspec-image}

Default: -

### labels (map[string]string, optional) {#fluentdspec-labels}

Default: -

### livenessDefaultCheck (bool, optional) {#fluentdspec-livenessdefaultcheck}

Default: -

### livenessProbe (*corev1.Probe, optional) {#fluentdspec-livenessprobe}

Default: -

### logLevel (string, optional) {#fluentdspec-loglevel}

Default: -

### metrics (*Metrics, optional) {#fluentdspec-metrics}

Default: -

### nodeSelector (map[string]string, optional) {#fluentdspec-nodeselector}

Default: -

### podPriorityClassName (string, optional) {#fluentdspec-podpriorityclassname}

Default: -

### port (int32, optional) {#fluentdspec-port}

Fluentd port inside the container (24240 by default) The headless service port is controlled by this field as well Note, that the default ClusterIP service port is always 24240 regardless of this field 

Default: -

### readinessDefaultCheck (ReadinessDefaultCheck, optional) {#fluentdspec-readinessdefaultcheck}

Default: -

### readinessProbe (*corev1.Probe, optional) {#fluentdspec-readinessprobe}

Default: -

### resources (corev1.ResourceRequirements, optional) {#fluentdspec-resources}

Default: -

### rootDir (string, optional) {#fluentdspec-rootdir}

Default: -

### scaling (*FluentdScaling, optional) {#fluentdspec-scaling}

Default: -

### security (*Security, optional) {#fluentdspec-security}

Default: -

### serviceAccount (*typeoverride.ServiceAccount, optional) {#fluentdspec-serviceaccount}

Default: -

### statefulsetAnnotations (map[string]string, optional) {#fluentdspec-statefulsetannotations}

Default: -

### tls (FluentdTLS, optional) {#fluentdspec-tls}

Default: -

### tolerations ([]corev1.Toleration, optional) {#fluentdspec-tolerations}

Default: -

### topologySpreadConstraints ([]corev1.TopologySpreadConstraint, optional) {#fluentdspec-topologyspreadconstraints}

Default: -

### volumeModImage (ImageSpec, optional) {#fluentdspec-volumemodimage}

Default: -

### volumeMountChmod (bool, optional) {#fluentdspec-volumemountchmod}

Default: -

### workers (int32, optional) {#fluentdspec-workers}

Default: -


## FluentOutLogrotate

### age (string, optional) {#fluentoutlogrotate-age}

Default: -

### enabled (bool, required) {#fluentoutlogrotate-enabled}

Default: -

### path (string, optional) {#fluentoutlogrotate-path}

Default: -

### size (string, optional) {#fluentoutlogrotate-size}

Default: -


## ExtraVolume

ExtraVolume defines the fluentd extra volumes

### containerName (string, optional) {#extravolume-containername}

Default: -

### path (string, optional) {#extravolume-path}

Default: -

### volume (*volume.KubernetesVolume, optional) {#extravolume-volume}

Default: -

### volumeName (string, optional) {#extravolume-volumename}

Default: -


## FluentdScaling

FluentdScaling enables configuring the scaling behaviour of the fluentd statefulset

### drain (FluentdDrainConfig, optional) {#fluentdscaling-drain}

Default: -

### podManagementPolicy (string, optional) {#fluentdscaling-podmanagementpolicy}

Default: -

### replicas (int, optional) {#fluentdscaling-replicas}

Default: -


## FluentdTLS

FluentdTLS defines the TLS configs

### enabled (bool, required) {#fluentdtls-enabled}

Default: -

### secretName (string, optional) {#fluentdtls-secretname}

Default: -

### sharedKey (string, optional) {#fluentdtls-sharedkey}

Default: -


## FluentdDrainConfig

FluentdDrainConfig enables configuring the drain behavior when scaling down the fluentd statefulset

### annotations (map[string]string, optional) {#fluentddrainconfig-annotations}

Annotations to use for the drain watch sidecar 

Default: -

### deleteVolume (bool, optional) {#fluentddrainconfig-deletevolume}

Should persistent volume claims be deleted after draining is done 

Default: -

### enabled (bool, optional) {#fluentddrainconfig-enabled}

Should buffers on persistent volumes left after scaling down the statefulset be drained 

Default: -

### image (ImageSpec, optional) {#fluentddrainconfig-image}

Default: -

### labels (map[string]string, optional) {#fluentddrainconfig-labels}

Labels to use for the drain watch sidecar on top of labels added by the operator by default. Default values can be overwritten. 

Default: -

### pauseImage (ImageSpec, optional) {#fluentddrainconfig-pauseimage}

Container image to use for the fluentd placeholder pod 

Default: -

### resources (*corev1.ResourceRequirements, optional) {#fluentddrainconfig-resources}

Configurable resource requirements for the drainer sidecar container. Default 20m cpu request, 20M memory limit 

Default: -

### securityContext (*corev1.SecurityContext, optional) {#fluentddrainconfig-securitycontext}

Configurable security context, uses fluentd pods' security context by default 

Default: -


