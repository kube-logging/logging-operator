<p align="center"><img src="./img/fluentbit.png" width="340"></p>

# Fluent Bit 

Fluent Bit is an open source and multi-platform Log Processor and Forwarder which allows you to collect data/logs from different sources, unify and send them to multiple destinations. It's fully compatible with Docker and Kubernetes environments.

Fluent Bit is written in C, have a pluggable architecture supporting around 30 extensions. It's fast and lightweight and provide the required security for network operations through TLS.

Current Version: [v1.3.6](https://github.com/fluent/fluent-bit/releases/tag/v1.3.6)

## Filters
### Kubernetes (filterKubernetes)
Fluent Bit Kubernetes Filter allows to enrich your log files with Kubernetes metadata.
[More info](https://github.com/fluent/fluent-bit-docs/blob/master/filter/kubernetes.md)

#### Example filter configurations
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit:
    filterKubernetes:
       Kube_URL: "https://kubernetes.default.svc:443"
       Match: "kube.*"
  controlNamespace: logging
```

#### Configuration Parameters

The plugin supports the following configuration parameters:

| Key | Description | Default |
| :--- | :--- | :--- |
| Buffer\_Size | Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the [Unit Size](../configuration/unit_sizes.md) specification. | 32k |
| Kube\_URL      | API Server end-point  | https://kubernetes.default.svc:443 |
| Kube\_CA\_File | CA certificate file   | /var/run/secrets/kubernetes.io/serviceaccount/ca.crt|
| Kube\_CA\_Path | Absolute path to scan for certificate files |  |
| Kube\_Token\_File | Token file | /var/run/secrets/kubernetes.io/serviceaccount/token |
| Kube_Tag_Prefix | When the source records comes from Tail input plugin, this option allows to specify what's the prefix used in Tail configuration. | kube.var.log.containers. |
| Merge\_Log | When enabled, it checks if the `log` field content is a JSON string map, if so, it append the map fields as part of the log structure. | Off |
| Merge\_Log\_Key | When `Merge_Log` is enabled, the filter tries to assume the `log` field from the incoming message is a JSON string message and make a structured representation of it at the same level of the `log` field in the map. Now if `Merge_Log_Key` is set \(a string name\), all the new structured fields taken from the original `log` content are inserted under the new key. |  |
| Merge\_Log\_Trim | When `Merge_Log` is enabled, trim (remove possible \n or \r) field values. | On |
| Merge\_Parser | Optional parser name to specify how to parse the data contained in the _log_ key. Recommended use is for developers or testing only.  | |
| Keep\_Log | When `Keep_Log` is disabled, the `log` field is removed from the incoming message once it has been successfully merged (`Merge_Log` must be enabled as well). | On |
| tls\_debug | Debug level between 0 \(nothing\) and 4 \(every detail\). | -1 |
| tls\_verify | When enabled, turns on certificate validation when connecting to the Kubernetes API server. | On |
| Use\_Journal | When enabled, the filter reads logs coming in Journald format. | Off |
| Regex\_Parser | Set an alternative Parser to process record Tag and extract pod\_name, namespace\_name, container\_name and docker\_id. The parser must be registered in a [parsers file](https://github.com/fluent/fluent-bit/blob/master/conf/parsers.conf) \(refer to parser _filter-kube-test_ as an example\). |  |
| K8S\_Logging\_Parser | Allow Kubernetes Pods to suggest a pre-defined Parser (read more about it in Kubernetes Annotations section) | Off |
| K8S\_Logging\_Exclude | Allow Kubernetes Pods to exclude their logs from the log processor (read more about it in Kubernetes Annotations section). | Off |
| Labels | Include Kubernetes resource labels in the extra metadata. | On |
| Annotations | Include Kubernetes resource annotations in the extra metadata. | On |
| Kube\_meta_preload_cache_dir | If set, Kubernetes meta-data can be cached/pre-loaded from files in JSON format in this directory, named as namespace-pod.meta | |
| Dummy\_Meta | If set, use dummy-meta data (for test/dev purposes) | Off |


## Inputs
### Tail (inputTail)
The tail input plugin allows to monitor one or several text files. It has a similar behavior like tail -f shell command.[More Info](https://github.com/fluent/fluent-bit-docs/blob/1.3/input/tail.md)

#### Example input configurations
```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit:
    inputTail:
       Refresh_Interval: "60"
       Rotate_Wait: "5"
  controlNamespace: logging
```

#### Configuration Parameters
The plugin supports the following configuration parameters:

| Key | Description | Default |
| :--- | :--- | :--- |
| storage.type | Specify the buffering mechanism to use. It can be memory or filesystem. | memory |
| Buffer\_Chunk\_Size | Set the initial buffer size to read files data. This value is used too to increase buffer size. The value must be according to the [Unit Size](../configuration/unit_sizes.md) specification. | 32k |
| Buffer\_Max\_Size | Set the limit of the buffer size per monitored file. When a buffer needs to be increased \(e.g: very long lines\), this value is used to restrict how much the memory buffer can grow. If reading a file exceed this limit, the file is removed from the monitored file list. The value must be according to the [Unit Size](../configuration/unit_sizes.md) specification. | Buffer\_Chunk\_Size |
| Path | Pattern specifying a specific log files or multiple ones through the use of common wildcards. |  |
| Path\_Key | If enabled, it appends the name of the monitored file as part of the record. The value assigned becomes the key in the map. |  |
| Exclude\_Path | Set one or multiple shell patterns separated by commas to exclude files matching a certain criteria, e.g: exclude\_path=\*.gz,\*.zip |  |
| Refresh\_Interval | The interval of refreshing the list of watched files in seconds. | 60 |
| Rotate\_Wait | Specify the number of extra time in seconds to monitor a file once is rotated in case some pending data is flushed. | 5 |
| Ignore\_Older | Ignores files that have been last modified before this time in seconds. Supports m,h,d \(minutes, hours,days\) syntax. Default behavior is to read all specified files. |  |
| Skip\_Long\_Lines | When a monitored file reach it buffer capacity due to a very long line \(Buffer\_Max\_Size\), the default behavior is to stop monitoring that file. Skip\_Long\_Lines alter that behavior and instruct Fluent Bit to skip long lines and continue processing other lines that fits into the buffer size. | Off |
| DB | Specify the database file to keep track of monitored files and offsets. |  |
| DB.Sync | Set a default synchronization \(I/O\) method. Values: Extra, Full, Normal, Off. This flag affects how the internal SQLite engine do synchronization to disk, for more details about each option please refer to [this section](https://www.sqlite.org/pragma.html#pragma_synchronous). | Full |
| Mem\_Buf\_Limit | Set a limit of memory that Tail plugin can use when appending data to the Engine. If the limit is reach, it will be paused; when the data is flushed it resumes. |  |
| Parser | Specify the name of a parser to interpret the entry as a structured message. |  |
| Key | When a message is unstructured \(no parser applied\), it's appended as a string under the key name _log_. This option allows to define an alternative name for that key. | log |
| Tag | Set a tag \(with regex-extract fields\) that will be placed on lines read. E.g. `kube.<namespace_name>.<pod_name>.<container_name>` |  |
| Tag\_Regex | Set a regex to exctract fields from the file. E.g. `(?<pod_name>[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*)_(?<namespace_name>[^_]+)_(?<container_name>.+)-` |  |
| Multiline | If enabled, the plugin will try to discover multiline messages and use the proper parsers to compose the outgoing messages. Note that when this option is enabled the Parser option is not used. | Off |
| Multiline\_Flush | Wait period time in seconds to process queued multiline messages | 4 |
| Parser\_Firstline | Name of the parser that matchs the beginning of a multiline message. Note that the regular expression defined in the parser must include a group name \(named capture\) |  |
| Parser\_N | Optional-extra parser to interpret and structure multiline entries. This option can be used to define multiple parsers, e.g: Parser\_1 ab1,  Parser\_2 ab2, Parser\_N abN. |  |
| Docker\_Mode | If enabled, the plugin will recombine split Docker log lines before passing them to any parser as configured above. This mode cannot be used at the same time as Multiline. | Off |
| Docker\_Mode\_Flush | Wait period time in seconds to flush queued unfinished split lines. | 4 |

## Buffering

### BufferStorage
A mechanism to place processed data into a temporal location until is ready to be shipped. [More Info](https://docs.fluentbit.io/manual/configuration/buffering)


| Key | Description | Default | 
| :--- | :--- | :--- |
| storage.path | Set an optional location in the file system to store streams and chunks of data. If this parameter is not set, Input plugins can only use in-memory buffering. | |
| storage.sync | Configure the synchronization mode used to store the data into the file system. It can take the values normal or full. | normal |
| storage.checksum | Enable the data integrity check when writing and reading data from the filesystem. The storage layer uses the CRC32 algorithm. | Off |
| storage.backlog.mem_limit | If storage.path is set, Fluent Bit will look for data chunks that were not delivered and are still in the storage layer, these are called backlog data. This option configure a hint of maximum value of memory to use when processing these records. | 5M |


#### Default configuration

If nothing is set, by default it configures the `storage.path` explicitly to use `/buffers` and leaves fluent-bit defaults for the other options. 

```
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit:
    bufferStorage:
       storage.path: /buffers
  controlNamespace: logging
```
