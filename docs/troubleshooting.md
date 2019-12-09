<p align="center"><img src="./img/troubleshooting.svg" width="260"></p>
<p align="center">


# Logging Operator Troubleshooting

## Check to ensure the necessary CRD are installed.
`kubectl get crd`
```bash
clusterflows.logging.banzaicloud.io     2019-12-05T15:11:48Z
clusteroutputs.logging.banzaicloud.io   2019-12-05T15:11:48Z
flows.logging.banzaicloud.io            2019-12-05T15:11:48Z
loggings.logging.banzaicloud.io         2019-12-05T15:11:48Z
outputs.logging.banzaicloud.io          2019-12-05T15:11:48Z
```

## Running Operator POD
`kubectl get pods |grep logging-operator`
```bash
NAME                                          READY   STATUS      RESTARTS   AGE
logging-demo-log-generator-6448d45cd9-z7zk8   1/1     Running     0          24m
```

---

<p align="center"><img src="./img/fluentbit.png" height="100"></p>

## Check Fluentbit daemonset 
`kubectl get daemonsets`
```bash
NAME                     DESIRED   CURRENT   READY   UP-TO-DATE   AVAILABLE   NODE SELECTOR   AGE
logging-demo-fluentbit   1         1         1       1            1           <none>          110s
```

### Configuration secret
`kubectl get secret logging-demo-fluentbit -o jsonpath="{.data['fluent-bit\.conf']}" | base64 --decode`
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
### Debug fluentbit Container
All fluentbit image has debug tag version you can use like this:
`kubectl edit loggings.logging.banzaicloud.io logging-demo`
```yam
fluentbit:
    image:
      pullPolicy: Always
      repository: fluent/fluent-bit
      tag: 1.3.2-debug
``` 
or update directly the daemonset
`kubectl edit daemonsets logging-demo-fluentbit`

With this you can kubectl exec into the pod using `sh` and look around. <br>
e.g.: `kubectl exec -it logging-demo-fluentbit-778zg sh`


### temp dir

`kubectl exec -it logging-demo-fluentbit-9dpzg ls  /buffers`

---

<p align="center"><img src="./img/fluentd.png" height="100"></p>

### Pod status (statefulset)
`kubectl get statefulsets`
```bash
NAME                   READY   AGE
logging-demo-fluentd   1/1     1m
```
### ConfigCheck How its works
### Check Configuration
`kubectl get secret logging-demo-fluentd-app -o jsonpath="{.data['fluentd\.conf']}" | base64 --decode`
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
`kubectl edit loggings.logging.banzaicloud.io logging-demo`
```yaml 
fluentd:
  logLevel: debug
```

### Get fluentd logs
`kubectl exec -it logging-demo-fluentd-0 cat /fluentd/log/out`


### Set stdout as an output 
You can use an stdout filter for this at any point in the flow to dump logs to the fluentd container’s stdout:<br>
e.g.: `kubectl edit loggings.logging.banzaicloud.io logging-demo`
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

### check the buffer path in the fluentd container
`kubectl exec -it logging-demo-fluentd-0 ls  /buffers`
```bash
Defaulting container name to fluentd.
Use 'kubectl describe pod/logging-demo-fluentd-0 -n logging' to see all of the containers in this pod.
logging_logging-demo-flow_logging-demo-output-minio_s3.b598f7eb0b2b34076b6da13a996ff2671.buffer
logging_logging-demo-flow_logging-demo-output-minio_s3.b598f7eb0b2b34076b6da13a996ff2671.buffer.meta
```

If you encounter any problems that the documentation does not address, please [file an issue](https://github.com/banzaicloud/logging-operator/issues) or talk to us on the Banzai Cloud Slack channel [#logging-operator](https://slack.banzaicloud.io/).

## Still In Trouble
If you are still in trouble with the operator you can contact us in the #logging-operator slack channel [here](https://slack.banzaicloud.io/).
The following information gives a huge help for us:
- logging-operator version
- kubernetes version
- helm/chart version(if you are install with helm)
- operator logs
- fluentd config
- fluentd logs
- fluentbit config
- fluentbit logs
Please be sure you cleaned up all the sensitive information before sharing.
