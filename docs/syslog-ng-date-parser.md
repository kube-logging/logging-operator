## syslog-ng date parser

By default, the syslog-ng aggregator uses the time when a message has been received on its input source as the timestamp.
In case we want to use the timestamp written in the message metadata, we should use a [date-parser](https://axoflow.com/docs/axosyslog-core/chapter-parsers/date-parser/date-parser-options/).

To enable the timestamps written by the container runtime (_cri_ or _docker_) and parsed by fluentbit automatically, we just
have to define the `sourceDateParser` in the _syslog-ng_ spec.

```
kind: Logging
metadata:
  name: example
spec:
  syslogNG:
    sourceDateParser: {}
```

In case we want to define our own parser format and template we can also do so (these are the default values):

```
kind: Logging
metadata:
  name: example
spec:
  syslogNG:
    sourceDateParser:
      format: "%FT%T.%f%z"
      template: "${json.time}"
```
