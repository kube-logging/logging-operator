package fluentd

var fluentdDefaultTemplate = `
# include other config files
@include /fluentd/etc/input.conf
@include /fluentd/app-config/*
@include /fluentd/etc/devnull.conf
`
var fluentdInputTemplate = `
# Enable RPC endpoint (this allows to trigger config reload without restart)
<system>
  rpc_endpoint 127.0.0.1:24444
</system>
# Prometheus monitoring
<source>
    @type prometheus
</source>
<source>
    @type prometheus_monitor
</source>
<source>
    @type prometheus_output_monitor
</source>

# Input plugin
<source>
    @type   forward
    port    24240
    @log_level debug
    {{ if .TLS.Enabled }}
    <security>
      self_hostname fluentd
      shared_key {{ .TLS.SharedKey }}
    </security>
    <transport tls>
      version                TLSv1_2
      ca_path                /fluentd/tls/caCert
      cert_path              /fluentd/tls/serverCert
      private_key_path       /fluentd/tls/serverKey
      client_cert_auth       true
    </transport>
    {{- end }}
</source>

# Prevent fluentd from handling records containing its own logs. Otherwise
# it can lead to an infinite loop, when error in sending one message generates
# another message which also fails to be sent and so on.
<match **.fluentd**>
    @type null
</match>

<match **.fluent-bit**>
    @type null
</match>

<match kubernetes.**>
  @type rewrite_tag_filter
  <rule>
    key $.kubernetes.namespace_name
    pattern ^(.+)$
    tag $1.${tag_parts[0]}
  </rule>
</match>

<match *.kubernetes.**>
  @type rewrite_tag_filter
  <rule>
    key $.kubernetes.labels.app
    pattern ^(.+)$
    tag $1.${tag_parts[0]}.${tag_parts[1]}
  </rule>
</match>
`
var fluentdOutputTemplate = `
<match **>
    @type null
</match>
`
