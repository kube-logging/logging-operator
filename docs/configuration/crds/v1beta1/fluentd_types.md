---
title: FluentdSpec
weight: 200
generated_file: true
---

## FluentdSpec

FluentdSpec defines the desired state of Fluentd

### affinity (*corev1.Affinity, optional) {#fluentdspec-affinity}


### annotations (map[string]string, optional) {#fluentdspec-annotations}


### bufferStorageVolume (volume.KubernetesVolume, optional) {#fluentdspec-bufferstoragevolume}

BufferStorageVolume is by default configured as PVC using FluentdPvcSpec [volume.KubernetesVolume](https://github.com/cisco-open/operator-tools/tree/master/docs/types) 


### bufferVolumeArgs ([]string, optional) {#fluentdspec-buffervolumeargs}


### bufferVolumeImage (ImageSpec, optional) {#fluentdspec-buffervolumeimage}


### bufferVolumeMetrics (*Metrics, optional) {#fluentdspec-buffervolumemetrics}


### bufferVolumeResources (corev1.ResourceRequirements, optional) {#fluentdspec-buffervolumeresources}


### compressConfigFile (bool, optional) {#fluentdspec-compressconfigfile}


### configCheckAnnotations (map[string]string, optional) {#fluentdspec-configcheckannotations}


### configCheckResources (corev1.ResourceRequirements, optional) {#fluentdspec-configcheckresources}


### configReloaderImage (ImageSpec, optional) {#fluentdspec-configreloaderimage}


### configReloaderResources (corev1.ResourceRequirements, optional) {#fluentdspec-configreloaderresources}


### dnsConfig (*corev1.PodDNSConfig, optional) {#fluentdspec-dnsconfig}


### dnsPolicy (corev1.DNSPolicy, optional) {#fluentdspec-dnspolicy}


### disablePvc (bool, optional) {#fluentdspec-disablepvc}


### enableMsgpackTimeSupport (bool, optional) {#fluentdspec-enablemsgpacktimesupport}

Allows Time object in buffer's MessagePack serde [more info]( https://docs.fluentd.org/deployment/system-config#enable_msgpack_time_support) 


### envVars ([]corev1.EnvVar, optional) {#fluentdspec-envvars}


### extraArgs ([]string, optional) {#fluentdspec-extraargs}


### extraVolumes ([]ExtraVolume, optional) {#fluentdspec-extravolumes}


### fluentLogDestination (string, optional) {#fluentdspec-fluentlogdestination}


### fluentOutLogrotate (*FluentOutLogrotate, optional) {#fluentdspec-fluentoutlogrotate}

FluentOutLogrotate sends fluent's stdout to file and rotates it 


### fluentdPvcSpec (*volume.KubernetesVolume, optional) {#fluentdspec-fluentdpvcspec}

Deprecated, use bufferStorageVolume 


### forwardInputConfig (*input.ForwardInputConfig, optional) {#fluentdspec-forwardinputconfig}


### ignoreRepeatedLogInterval (string, optional) {#fluentdspec-ignorerepeatedloginterval}

Ignore repeated log lines [more info]( https://docs.fluentd.org/deployment/logging#ignore_repeated_log_interval) 


### ignoreSameLogInterval (string, optional) {#fluentdspec-ignoresameloginterval}

Ignore same log lines [more info]( https://docs.fluentd.org/deployment/logging#ignore_same_log_interval) 


### image (ImageSpec, optional) {#fluentdspec-image}


### labels (map[string]string, optional) {#fluentdspec-labels}


### livenessDefaultCheck (bool, optional) {#fluentdspec-livenessdefaultcheck}


### livenessProbe (*corev1.Probe, optional) {#fluentdspec-livenessprobe}


### logLevel (string, optional) {#fluentdspec-loglevel}


### metrics (*Metrics, optional) {#fluentdspec-metrics}


### nodeSelector (map[string]string, optional) {#fluentdspec-nodeselector}


### pdb (*PdbInput, optional) {#fluentdspec-pdb}


### podPriorityClassName (string, optional) {#fluentdspec-podpriorityclassname}


### port (int32, optional) {#fluentdspec-port}

Fluentd port inside the container (24240 by default). The headless service port is controlled by this field as well. Note that the default ClusterIP service port is always 24240, regardless of this field. 


### readinessDefaultCheck (ReadinessDefaultCheck, optional) {#fluentdspec-readinessdefaultcheck}


### readinessProbe (*corev1.Probe, optional) {#fluentdspec-readinessprobe}


### resources (corev1.ResourceRequirements, optional) {#fluentdspec-resources}


### rootDir (string, optional) {#fluentdspec-rootdir}


### scaling (*FluentdScaling, optional) {#fluentdspec-scaling}


### security (*Security, optional) {#fluentdspec-security}


### serviceAccount (*typeoverride.ServiceAccount, optional) {#fluentdspec-serviceaccount}


### sidecarContainers ([]corev1.Container, optional) {#fluentdspec-sidecarcontainers}

Available in Logging operator version 4.5 and later. Configure sidecar container in Fluentd pods, for example: [https://github.com/kube-logging/logging-operator/config/samples/logging_logging_fluentd_sidecars.yaml](https://github.com/kube-logging/logging-operator/config/samples/logging_logging_fluentd_sidecars.yaml). 


### statefulsetAnnotations (map[string]string, optional) {#fluentdspec-statefulsetannotations}


### tls (FluentdTLS, optional) {#fluentdspec-tls}


### tolerations ([]corev1.Toleration, optional) {#fluentdspec-tolerations}


### topologySpreadConstraints ([]corev1.TopologySpreadConstraint, optional) {#fluentdspec-topologyspreadconstraints}


### volumeModImage (ImageSpec, optional) {#fluentdspec-volumemodimage}


### volumeMountChmod (bool, optional) {#fluentdspec-volumemountchmod}


### workers (int32, optional) {#fluentdspec-workers}



## FluentOutLogrotate

### age (string, optional) {#fluentoutlogrotate-age}


### enabled (bool, required) {#fluentoutlogrotate-enabled}


### path (string, optional) {#fluentoutlogrotate-path}


### size (string, optional) {#fluentoutlogrotate-size}



## ExtraVolume

ExtraVolume defines the fluentd extra volumes

### containerName (string, optional) {#extravolume-containername}


### path (string, optional) {#extravolume-path}


### volume (*volume.KubernetesVolume, optional) {#extravolume-volume}


### volumeName (string, optional) {#extravolume-volumename}



## FluentdScaling

FluentdScaling enables configuring the scaling behaviour of the fluentd statefulset

### drain (FluentdDrainConfig, optional) {#fluentdscaling-drain}


### podManagementPolicy (string, optional) {#fluentdscaling-podmanagementpolicy}


### replicas (int, optional) {#fluentdscaling-replicas}



## FluentdTLS

FluentdTLS defines the TLS configs

### enabled (bool, required) {#fluentdtls-enabled}


### secretName (string, optional) {#fluentdtls-secretname}


### sharedKey (string, optional) {#fluentdtls-sharedkey}



## FluentdDrainConfig

FluentdDrainConfig enables configuring the drain behavior when scaling down the fluentd statefulset

### annotations (map[string]string, optional) {#fluentddrainconfig-annotations}

Annotations to use for the drain watch sidecar 


### deleteVolume (bool, optional) {#fluentddrainconfig-deletevolume}

Should persistent volume claims be deleted after draining is done 


### enabled (bool, optional) {#fluentddrainconfig-enabled}

Should buffers on persistent volumes left after scaling down the statefulset be drained 


### image (ImageSpec, optional) {#fluentddrainconfig-image}


### labels (map[string]string, optional) {#fluentddrainconfig-labels}

Labels to use for the drain watch sidecar on top of labels added by the operator by default. Default values can be overwritten. 


### pauseImage (ImageSpec, optional) {#fluentddrainconfig-pauseimage}

Container image to use for the fluentd placeholder pod 


### resources (*corev1.ResourceRequirements, optional) {#fluentddrainconfig-resources}

Available in Logging operator version 4.4 and later. Configurable resource requirements for the drainer sidecar container. Default 20m cpu request, 20M memory limit 


### securityContext (*corev1.SecurityContext, optional) {#fluentddrainconfig-securitycontext}

Available in Logging operator version 4.4 and later. Configurable security context, uses fluentd pods' security context by default 



## PdbInput

### maxUnavailable (*intstr.IntOrString, optional) {#pdbinput-maxunavailable}


### minAvailable (*intstr.IntOrString, optional) {#pdbinput-minavailable}


### unhealthyPodEvictionPolicy (*policyv1.UnhealthyPodEvictionPolicyType, optional) {#pdbinput-unhealthypodevictionpolicy}



