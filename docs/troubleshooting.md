<p align="center"><img src="./img/troubleshooting.svg" width="260"></p>
<p align="center">

# Logging operator troubleshooting

The following tips and commands can help you to troubleshoot your Logging operator installation.

## First things to do

1. Check that the necessary CRDs are installed. Issue the following command: `kubectl get crd`
   The output should include the following CRDs:
    ```bash
    clusterflows.logging.banzaicloud.io     2019-12-05T15:11:48Z
    clusteroutputs.logging.banzaicloud.io   2019-12-05T15:11:48Z
    flows.logging.banzaicloud.io            2019-12-05T15:11:48Z
    loggings.logging.banzaicloud.io         2019-12-05T15:11:48Z
    outputs.logging.banzaicloud.io          2019-12-05T15:11:48Z
    ```
1. Verify that the Logging operator pod is running. Issue the following command: `kubectl get pods |grep logging-operator`
   The output should include the a running pod, for example:
    ```bash
    NAME                                          READY   STATUS      RESTARTS   AGE
    logging-demo-log-generator-6448d45cd9-z7zk8   1/1     Running     0          24m
    ```

---

## Troubleshooting Fluent Bit

<p align="center"><img src="./img/fluentbit.png" height="100"></p>

The following sections help you troubleshoot the Fluent Bit component of the Logging operator.

### Check the Fluent Bit daemonset

Verify that the Fluent Bit daemonset is available. Issue the following command: `kubectl get daemonsets`
The output should include a Fluent Bit daemonset, for example:
```bash
NAME                     DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
logging-demo-fluentbit   1         1         1       1            1           <none>          110s
```

### Check the Fluent Bit configuration

You can display the current configuration of the Fluent Bit daemonset using the following command:
`kubectl get secret logging-demo-fluentbit -o jsonpath="{.data['fluent-bit\.conf']}" | base64 --decode`
The output looks like the following:
```yaml
[SERVICE]
    Flush        1
    Daemon       Off
    Log_Level    info
    Parsers_File parsers.conf
    storage.path  /buffers

[INPUT]
    Name         tail
    DB  /tail-db/tail-containers-state.db
    Mem_Buf_Limit  5MB
    Parser  docker
    Path  /var/log/containers/*.log
    Refresh_Interval  5
    Skip_Long_Lines  On
    Tag  kubernetes.*

[FILTER]
    Name        kubernetes
    Kube_CA_File  /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    Kube_Tag_Prefix  kubernetes.var.log.containers
    Kube_Token_File  /var/run/secrets/kubernetes.io/serviceaccount/token
    Kube_URL  https://kubernetes.default.svc:443
    Match  kubernetes.*
    Merge_Log  On

[OUTPUT]
    Name          forward
    Match         *
    Host          logging-demo-fluentd.logging.svc
    Port          24240

    tls           On
    tls.verify    Off
    tls.ca_file   /fluent-bit/tls/ca.crt
    tls.crt_file  /fluent-bit/tls/tls.crt
    tls.key_file  /fluent-bit/tls/tls.key
    Shared_Key    Kamk2_SukuWenk
    Retry_Limit   False
```

### Debug version of the fluentbit container

All Fluent Bit image tags have a debug version marked with the `-debug` suffix. You can install this debug version using the following command:
`kubectl edit loggings.logging.banzaicloud.io logging-demo`
```yam
fluentbit:
    image:
      pullPolicy: Always
      repository: fluent/fluent-bit
      tag: 1.3.2-debug
```


After deploying the debug version, you can kubectl exec into the pod using `sh` and look around. For example: `kubectl exec -it logging-demo-fluentbit-778zg sh`

### Check the queued log messages

You can check the buffer directory if Fluent Bit is configured to buffer queued log messages to disk instead of in memory. (You can configure it through the InputTail fluentbit config, by setting the `storage.type` field to `filesystem`.)

`kubectl exec -it logging-demo-fluentbit-9dpzg ls /buffers`

---
## Troubleshooting Fluentd

<p align="center"><img src="./img/fluentd.png" height="100"></p>

The following sections help you troubleshoot the Fluentd statefulset component of the Logging operator.

### Check Fluentd pod status (statefulset)

Verify that the Fluentd statefulset is available using the following command: `kubectl get statefulsets`

Expected output:
```bash
NAME                   READY   AGE
logging-demo-fluentd   1/1     1m
```

### ConfigCheck

The Logging operator has a builtin mechanism that validates the generated fluentd configuration before applying it to fluentd. You should be able to see the configcheck pod and it's log output. The result of the check is written into the `status` field of the corresponding `Logging` resource.

In case the operator is stuck in an error state caused by a failed configcheck, restore the previous configuration by modifying or removing the invalid resources to the point where the configcheck pod is finally able to complete successfully.

### Check Fluentd configuration

Use the following command to display the configuration of Fluentd:
`kubectl get secret logging-demo-fluentd-app -o jsonpath="{.data['fluentd\.conf']}" | base64 --decode`

The output should be similar to the following:
```yaml
<source>
  @type forward
  @id main_forward
  bind 0.0.0.0
  port 24240
  <transport tls>
    ca_path /fluentd/tls/ca.crt
    cert_path /fluentd/tls/tls.crt
    client_cert_auth true
    private_key_path /fluentd/tls/tls.key
    version TLSv1_2
  </transport>
  <security>
    self_hostname fluentd
    shared_key Kamk2_SukuWenk
  </security>
</source>
<match **>
  @type label_router
  @id main_label_router
  <route>
    @label @427b3e18f3a3bc3f37643c54e9fc960b
    labels app.kubernetes.io/instance:logging-demo,app.kubernetes.io/name:log-generator
    namespace logging
  </route>
</match>
<label @427b3e18f3a3bc3f37643c54e9fc960b>
  <match kubernetes.**>
    @type tag_normaliser
    @id logging-demo-flow_0_tag_normaliser
    format ${namespace_name}.${pod_name}.${container_name}
  </match>
  <filter **>
    @type parser
    @id logging-demo-flow_1_parser
    key_name log
    remove_key_name_field true
    reserve_data true
    <parse>
      @type nginx
    </parse>
  </filter>
  <match **>
    @type s3
    @id logging_logging-demo-flow_logging-demo-output-minio_s3
    aws_key_id WVKblQelkDTSKTn4aaef
    aws_sec_key LAmjIah4MTKTM3XGrDxuD2dTLLmysVHvZrtxpzK6
    force_path_style true
    path logs/${tag}/%Y/%m/%d/
    s3_bucket demo
    s3_endpoint http://logging-demo-minio.logging.svc.cluster.local:9000
    s3_region test_region
    <buffer tag,time>
      @type file
      path /buffers/logging_logging-demo-flow_logging-demo-output-minio_s3.*.buffer
      retry_forever true
      timekey 10s
      timekey_use_utc true
      timekey_wait 0s
    </buffer>
  </match>
</label>
```

### Set Fluentd log Level

Use the following command to change the log level of Fluentd.
`kubectl edit loggings.logging.banzaicloud.io logging-demo`
```yaml 
fluentd:
  logLevel: debug
```

### Get Fluentd logs

The following command displays the logs of the Fluentd container.
`kubectl exec -it logging-demo-fluentd-0 cat /fluentd/log/out`

### Set stdout as an output

You can use an stdout filter at any point in the flow to dump the log messages to the stdout of the Fluentd container. For example:
`kubectl edit loggings.logging.banzaicloud.io logging-demo`
```yaml
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: exchange
  namespace: logging
spec:
  filters:
    - stdout: {}
  outputRefs:
    - exchange
  selectors:
    application: exchange
```

### Check the buffer path in the fluentd container

`kubectl exec -it logging-demo-fluentd-0 ls  /buffers`
```bash
Defaulting container name to fluentd.
Use 'kubectl describe pod/logging-demo-fluentd-0 -n logging' to see all of the containers in this pod.
logging_logging-demo-flow_logging-demo-output-minio_s3.b598f7eb0b2b34076b6da13a996ff2671.buffer
logging_logging-demo-flow_logging-demo-output-minio_s3.b598f7eb0b2b34076b6da13a996ff2671.buffer.meta
```

### Other problems and getting support

If you encounter any problems that the documentation does not address, [file an issue](https://github.com/banzaicloud/logging-operator/issues) or talk to us on the Banzai Cloud Slack channel [#logging-operator](https://slack.banzaicloud.io/).

[Commercial support](https://banzaicloud.com/products/logging-operator/) is also available for the Logging operator.

Before asking for help, prepare the following information to make troubleshooting faster:

- Logging operator version
- kubernetes version
- helm/chart version (if you installed the Logging operator with helm)
- Logging operator logs
- [fluentd configuration](#check-fluentd-configuration)
- [fluentd logs](#get-fluentd-logs)
- [fluentbit configuration](#check-the-fluent-bit-configuration)
- fluentbit logs

Do not forget to remove any sensitive information (for example, passwords and private keys) before sharing.
