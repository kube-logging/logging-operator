# Telemtry Controller multi-tenant routing

```bash
git clone https://github.com/kube-logging/telemetry-controller
cd telemetry-controller
make kind-cluster
make docker-build
kind load docker-image controller:local
helm upgrade --install --wait --create-namespace --namespace telemetry-controller-system telemetry-controller ./charts/telemetry-controller/ --set image.repository=controller --set image.tag=local
kubectl apply -f ../logging-operator/config/samples/multitenant-routing-tc/tc-simple
helm upgrade --install --namespace customer-a log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
helm upgrade --install --namespace customer-b log-generator oci://ghcr.io/kube-logging/helm-charts/log-generator
```

## Expected generated config

```yaml
connectors:
  count/output_metrics:
    logs:
      telemetry_controller_output_log_count:
        attributes:
          - key: tenant
          - key: subscription
          - key: exporter
        description: The number of logs sent out from each exporter.
        resource_attributes:
          - key: k8s.namespace.name
          - key: k8s.node.name
          - key: k8s.container.name
          - key: k8s.pod.name
          - key: k8s.pod.labels.app.kubernetes.io/name
          - key: k8s.pod.labels.app
  count/tenant_metrics:
    logs:
      telemetry_controller_tenant_log_count:
        attributes:
          - key: tenant
        description: The number of logs from each tenant pipeline.
        resource_attributes:
          - key: k8s.namespace.name
          - key: k8s.node.name
          - key: k8s.container.name
          - key: k8s.pod.name
          - key: k8s.pod.labels.app.kubernetes.io/name
          - key: k8s.pod.labels.app
  routing/subscription_customer-a_customer-a_outputs:
    table:
      - condition: "true"
        pipelines:
          - logs/output_customer-a_customer-a_customer-a_customer-a-receiver
  routing/subscription_customer-b_customer-b_outputs:
    table:
      - condition: "true"
        pipelines:
          - logs/output_customer-b_customer-b_customer-b_customer-b-receiver
  routing/subscription_infra_infra_outputs:
    table:
      - condition: "true"
        pipelines:
          - logs/output_infra_infra_infra_infra-all
  routing/tenant_customer-a_subscriptions:
    table:
      - condition: "true"
        pipelines:
          - logs/tenant_customer-a_subscription_customer-a_customer-a
  routing/tenant_customer-b_subscriptions:
    table:
      - condition: "true"
        pipelines:
          - logs/tenant_customer-b_subscription_customer-b_customer-b
  routing/tenant_infra_subscriptions:
    table:
      - condition: "true"
        pipelines:
          - logs/tenant_infra_subscription_infra_infra
exporters:
  debug:
    verbosity: detailed
  otlp/customer-a_customer-a-receiver:
    endpoint: receiver-a-collector.customer-a.svc.cluster.local:4317
    tls:
      insecure: true
  otlp/customer-b_customer-b-receiver:
    endpoint: receiver-b-collector.customer-b.svc.cluster.local:4317
    tls:
      insecure: true
  otlp/infra_infra-all:
    endpoint: receiver-infra-collector.infra.svc.cluster.local:4317
    tls:
      insecure: true
  prometheus/message_metrics_exporter:
    endpoint: :9999
extensions: {}
processors:
  attributes/exporter_name_customer-a-receiver:
    actions:
      - action: insert
        key: exporter
        value: otlp/customer-a_customer-a-receiver
  attributes/exporter_name_customer-b-receiver:
    actions:
      - action: insert
        key: exporter
        value: otlp/customer-b_customer-b-receiver
  attributes/exporter_name_infra-all:
    actions:
      - action: insert
        key: exporter
        value: otlp/infra_infra-all
  attributes/metricattributes:
    actions:
      - action: insert
        from_attribute: k8s.pod.labels.app
        key: app
      - action: insert
        from_attribute: k8s.node.name
        key: host
      - action: insert
        from_attribute: k8s.namespace.name
        key: namespace
      - action: insert
        from_attribute: k8s.container.name
        key: container
      - action: insert
        from_attribute: k8s.pod.name
        key: pod
  attributes/subscription_customer-a:
    actions:
      - action: insert
        key: subscription
        value: customer-a
  attributes/subscription_customer-b:
    actions:
      - action: insert
        key: subscription
        value: customer-b
  attributes/subscription_infra:
    actions:
      - action: insert
        key: subscription
        value: infra
  attributes/tenant_customer-a:
    actions:
      - action: insert
        key: tenant
        value: customer-a
  attributes/tenant_customer-b:
    actions:
      - action: insert
        key: tenant
        value: customer-b
  attributes/tenant_infra:
    actions:
      - action: insert
        key: tenant
        value: infra
  deltatocumulative: {}
  k8sattributes:
    auth_type: serviceAccount
    extract:
      labels:
        - from: pod
          key_regex: .*
          tag_name: all_labels
      metadata:
        - k8s.pod.name
        - k8s.pod.uid
        - k8s.deployment.name
        - k8s.namespace.name
        - k8s.node.name
        - k8s.pod.start_time
    passthrough: false
    pod_association:
      - sources:
          - from: resource_attribute
            name: k8s.namespace.name
          - from: resource_attribute
            name: k8s.pod.name
  memory_limiter:
    check_interval: 1s
    limit_percentage: 75
    spike_limit_mib: 25
receivers:
  filelog/customer-a:
    exclude:
      - /var/log/pods/*/otc-container/*.log
    include:
      - /var/log/pods/customer-a_*/*/*.log
    include_file_name: false
    include_file_path: true
    operators:
      - id: get-format
        routes:
          - expr: body matches "^\\{"
            output: parser-docker
          - expr: body matches "^[^ Z]+Z"
            output: parser-containerd
        type: router
      - id: parser-containerd
        output: extract_metadata_from_filepath
        regex: ^(?P<time>[^ ^Z]+Z) (?P<stream>stdout|stderr) (?P<logtag>[^ ]*) ?(?P<log>.*)$
        timestamp:
          layout: '%Y-%m-%dT%H:%M:%S.%LZ'
          parse_from: attributes.time
        type: regex_parser
      - id: parser-docker
        output: extract_metadata_from_filepath
        timestamp:
          layout: '%Y-%m-%dT%H:%M:%S.%LZ'
          parse_from: attributes.time
        type: json_parser
      - cache:
          size: 128
        id: extract_metadata_from_filepath
        parse_from: attributes["log.file.path"]
        regex: ^.*\/(?P<namespace>[^_]+)_(?P<pod_name>[^_]+)_(?P<uid>[a-f0-9-]+)\/(?P<container_name>[^\/]+)\/(?P<restart_count>\d+)\.log$
        type: regex_parser
      - from: attributes.log
        to: body
        type: move
      - from: attributes.stream
        to: attributes["log.iostream"]
        type: move
      - from: attributes.container_name
        to: resource["k8s.container.name"]
        type: move
      - from: attributes.namespace
        to: resource["k8s.namespace.name"]
        type: move
      - from: attributes.pod_name
        to: resource["k8s.pod.name"]
        type: move
      - from: attributes.restart_count
        to: resource["k8s.container.restart_count"]
        type: move
      - from: attributes.uid
        to: resource["k8s.pod.uid"]
        type: move
    retry_on_failure:
      enabled: true
      max_elapsed_time: 0
    start_at: end
  filelog/customer-b:
    exclude:
      - /var/log/pods/*/otc-container/*.log
    include:
      - /var/log/pods/customer-b_*/*/*.log
    include_file_name: false
    include_file_path: true
    operators:
      - id: get-format
        routes:
          - expr: body matches "^\\{"
            output: parser-docker
          - expr: body matches "^[^ Z]+Z"
            output: parser-containerd
        type: router
      - id: parser-containerd
        output: extract_metadata_from_filepath
        regex: ^(?P<time>[^ ^Z]+Z) (?P<stream>stdout|stderr) (?P<logtag>[^ ]*) ?(?P<log>.*)$
        timestamp:
          layout: '%Y-%m-%dT%H:%M:%S.%LZ'
          parse_from: attributes.time
        type: regex_parser
      - id: parser-docker
        output: extract_metadata_from_filepath
        timestamp:
          layout: '%Y-%m-%dT%H:%M:%S.%LZ'
          parse_from: attributes.time
        type: json_parser
      - cache:
          size: 128
        id: extract_metadata_from_filepath
        parse_from: attributes["log.file.path"]
        regex: ^.*\/(?P<namespace>[^_]+)_(?P<pod_name>[^_]+)_(?P<uid>[a-f0-9-]+)\/(?P<container_name>[^\/]+)\/(?P<restart_count>\d+)\.log$
        type: regex_parser
      - from: attributes.log
        to: body
        type: move
      - from: attributes.stream
        to: attributes["log.iostream"]
        type: move
      - from: attributes.container_name
        to: resource["k8s.container.name"]
        type: move
      - from: attributes.namespace
        to: resource["k8s.namespace.name"]
        type: move
      - from: attributes.pod_name
        to: resource["k8s.pod.name"]
        type: move
      - from: attributes.restart_count
        to: resource["k8s.container.restart_count"]
        type: move
      - from: attributes.uid
        to: resource["k8s.pod.uid"]
        type: move
    retry_on_failure:
      enabled: true
      max_elapsed_time: 0
    start_at: end
  filelog/infra:
    exclude:
      - /var/log/pods/*/otc-container/*.log
    include:
      - /var/log/pods/customer-a_*/*/*.log
      - /var/log/pods/customer-b_*/*/*.log
    include_file_name: false
    include_file_path: true
    operators:
      - id: get-format
        routes:
          - expr: body matches "^\\{"
            output: parser-docker
          - expr: body matches "^[^ Z]+Z"
            output: parser-containerd
        type: router
      - id: parser-containerd
        output: extract_metadata_from_filepath
        regex: ^(?P<time>[^ ^Z]+Z) (?P<stream>stdout|stderr) (?P<logtag>[^ ]*) ?(?P<log>.*)$
        timestamp:
          layout: '%Y-%m-%dT%H:%M:%S.%LZ'
          parse_from: attributes.time
        type: regex_parser
      - id: parser-docker
        output: extract_metadata_from_filepath
        timestamp:
          layout: '%Y-%m-%dT%H:%M:%S.%LZ'
          parse_from: attributes.time
        type: json_parser
      - cache:
          size: 128
        id: extract_metadata_from_filepath
        parse_from: attributes["log.file.path"]
        regex: ^.*\/(?P<namespace>[^_]+)_(?P<pod_name>[^_]+)_(?P<uid>[a-f0-9-]+)\/(?P<container_name>[^\/]+)\/(?P<restart_count>\d+)\.log$
        type: regex_parser
      - from: attributes.log
        to: body
        type: move
      - from: attributes.stream
        to: attributes["log.iostream"]
        type: move
      - from: attributes.container_name
        to: resource["k8s.container.name"]
        type: move
      - from: attributes.namespace
        to: resource["k8s.namespace.name"]
        type: move
      - from: attributes.pod_name
        to: resource["k8s.pod.name"]
        type: move
      - from: attributes.restart_count
        to: resource["k8s.container.restart_count"]
        type: move
      - from: attributes.uid
        to: resource["k8s.pod.uid"]
        type: move
    retry_on_failure:
      enabled: true
      max_elapsed_time: 0
    start_at: end
service:
  pipelines:
    logs/output_customer-a_customer-a_customer-a_customer-a-receiver:
      exporters:
        - otlp/customer-a_customer-a-receiver
        - count/output_metrics
      processors:
        - memory_limiter
        - attributes/exporter_name_customer-a-receiver
      receivers:
        - routing/subscription_customer-a_customer-a_outputs
    logs/output_customer-b_customer-b_customer-b_customer-b-receiver:
      exporters:
        - otlp/customer-b_customer-b-receiver
        - count/output_metrics
      processors:
        - memory_limiter
        - attributes/exporter_name_customer-b-receiver
      receivers:
        - routing/subscription_customer-b_customer-b_outputs
    logs/output_infra_infra_infra_infra-all:
      exporters:
        - otlp/infra_infra-all
        - count/output_metrics
      processors:
        - memory_limiter
        - attributes/exporter_name_infra-all
      receivers:
        - routing/subscription_infra_infra_outputs
    logs/tenant_customer-a:
      exporters:
        - routing/tenant_customer-a_subscriptions
        - count/tenant_metrics
      processors:
        - memory_limiter
        - k8sattributes
        - attributes/tenant_customer-a
      receivers:
        - filelog/customer-a
    logs/tenant_customer-a_subscription_customer-a_customer-a:
      exporters:
        - routing/subscription_customer-a_customer-a_outputs
      processors:
        - memory_limiter
        - attributes/subscription_customer-a
      receivers:
        - routing/tenant_customer-a_subscriptions
    logs/tenant_customer-b:
      exporters:
        - routing/tenant_customer-b_subscriptions
        - count/tenant_metrics
      processors:
        - memory_limiter
        - k8sattributes
        - attributes/tenant_customer-b
      receivers:
        - filelog/customer-b
    logs/tenant_customer-b_subscription_customer-b_customer-b:
      exporters:
        - routing/subscription_customer-b_customer-b_outputs
      processors:
        - memory_limiter
        - attributes/subscription_customer-b
      receivers:
        - routing/tenant_customer-b_subscriptions
    logs/tenant_infra:
      exporters:
        - routing/tenant_infra_subscriptions
        - count/tenant_metrics
      processors:
        - memory_limiter
        - k8sattributes
        - attributes/tenant_infra
      receivers:
        - filelog/infra
    logs/tenant_infra_subscription_infra_infra:
      exporters:
        - routing/subscription_infra_infra_outputs
      processors:
        - memory_limiter
        - attributes/subscription_infra
      receivers:
        - routing/tenant_infra_subscriptions
    metrics/output:
      exporters:
        - prometheus/message_metrics_exporter
      processors:
        - memory_limiter
        - deltatocumulative
        - attributes/metricattributes
      receivers:
        - count/output_metrics
    metrics/tenant:
      exporters:
        - prometheus/message_metrics_exporter
      processors:
        - memory_limiter
        - deltatocumulative
        - attributes/metricattributes
      receivers:
        - count/tenant_metrics
  telemetry:
    metrics:
      level: detailed
      readers:
        - pull:
            exporter:
              prometheus:
                host: ""
                port: 8888
```
