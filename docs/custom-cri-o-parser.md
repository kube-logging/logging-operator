CRI-O has a different date format then containerd:
https://opentelemetry.io/blog/2024/otel-collector-container-log-parser/#how-container-logs-look-like

- Containerd: `%Y-%m-%dT%H:%M:%S.%L%z`
- CRI-O: `%Y-%m-%dT%H:%M:%S.%9N%:z`

Until supported by the https://kube-logging.dev/docs/whats-new/#containerd-compatibility flag this custom parser can be applied like this through the chart: 
```
logging:
  fluentbit:
    inputTail:
      Parser: "cri-parser"
    customParsers: |
      [PARSER]
          Name cri-parser
          Format regex
          Regex ^(?<time>[^ ]+) (?<stream>stdout|stderr) (?<logtag>[^ ]*) (?<log>.*)$
          Time_Key time
          Time_Format %Y-%m-%dT%H:%M:%S.%9N%:z
```
