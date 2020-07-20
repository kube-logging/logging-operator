---
title: FluentbitSpec
weight: 200
generated_file: true
---

### FluentbitSpec
#### FluentbitSpec defines the desired state of Fluentbit

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| annotations | map[string]string | No | - |  |
| labels | map[string]string | No | - |  |
| image | ImageSpec | No | - |  |
| tls | FluentbitTLS | No | - |  |
| targetHost | string | No | - |  |
| targetPort | int32 | No | - |  |
| resources | corev1.ResourceRequirements | No | - |  |
| tolerations | []corev1.Toleration | No | - |  |
| nodeSelector | map[string]string | No | - |  |
| affinity | *corev1.Affinity | No | - |  |
| metrics | *Metrics | No | - |  |
| security | *Security | No | - |  |
| positiondb | volume.KubernetesVolume | No | - | [volume.KubernetesVolume](https://github.com/banzaicloud/operator-tools/tree/master/docs/types)<br> |
| position_db | *volume.KubernetesVolume | No | - | Deprecated, use positiondb<br> |
| mountPath | string | No | - |  |
| extraVolumeMounts | []VolumeMount | No | - |  |
| inputTail | InputTail | No | - |  |
| filterAws | *FilterAws | No | - |  |
| parser | string | No | - | Deprecated, use inputTail.parser<br> |
| filterKubernetes | FilterKubernetes | No | - |  |
| bufferStorage | BufferStorage | No | - |  |
| bufferStorageVolume | volume.KubernetesVolume | No | - | [volume.KubernetesVolume](https://github.com/banzaicloud/operator-tools/tree/master/docs/types)<br> |
| customConfigSecret | string | No | - |  |
| podPriorityClassName | string | No | - |  |
| livenessProbe | *corev1.Probe | No | - |  |
| livenessDefaultCheck | bool | No | - |  |
| readinessProbe | *corev1.Probe | No | - |  |
### FluentbitTLS
#### FluentbitTLS defines the TLS configs

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| enabled | bool | Yes | - |  |
| secretName | string | Yes | - |  |
| sharedKey | string | No | - |  |
### BufferStorage
#### BufferStorage is the Service Section Configuration of fluent-bit

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| storage.path | string | No | - | Set an optional location in the file system to store streams and chunks of data. If this parameter is not set, Input plugins can only use in-memory buffering.<br> |
| storage.sync | string | No | normal | Configure the synchronization mode used to store the data into the file system. It can take the values normal or full. <br> |
| storage.checksum | string | No | Off | Enable the data integrity check when writing and reading data from the filesystem. The storage layer uses the CRC32 algorithm. <br> |
| storage.backlog.mem_limit | string | No | 5M | If storage.path is set, Fluent Bit will look for data chunks that were not delivered and are still in the storage layer, these are called backlog data. This option configure a hint of maximum value of memory to use when processing these records. <br> |
### InputTail
#### InputTail defines Fluentbit tail input configuration The tail input plugin allows to monitor one or several text files. It has a similar behavior like tail -f shell command.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| storage.type | string | No | memory | Specify the buffering mechanism to use. It can be memory or filesystem. <br> |
| Buffer_Chunk_Size | string | No | 32k | Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the Unit Size specification. <br> |
| Buffer_Max_Size | string | No | Buffer_Chunk_Size | Set the limit of the buffer size per monitored file. When a buffer needs to be increased (e.g: very long lines), this value is used to restrict how much the memory buffer can grow. If reading a file exceed this limit, the file is removed from the monitored file list. The value must be according to the Unit Size specification. <br> |
| Path | string | No | - | Pattern specifying a specific log files or multiple ones through the use of common wildcards.<br> |
| Path_Key | string | No | - | If enabled, it appends the name of the monitored file as part of the record. The value assigned becomes the key in the map.<br> |
| Exclude_Path | string | No | - | Set one or multiple shell patterns separated by commas to exclude files matching a certain criteria, e.g: exclude_path=*.gz,*.zip<br> |
| Refresh_Interval | string | No | 60 | The interval of refreshing the list of watched files in seconds. <br> |
| Rotate_Wait | string | No | 5 | Specify the number of extra time in seconds to monitor a file once is rotated in case some pending data is flushed. <br> |
| Ignore_Older | string | No | - | Ignores files that have been last modified before this time in seconds. Supports m,h,d (minutes, hours,days) syntax. Default behavior is to read all specified files.<br> |
| Skip_Long_Lines | string | No | Off | When a monitored file reach it buffer capacity due to a very long line (Buffer_Max_Size), the default behavior is to stop monitoring that file. Skip_Long_Lines alter that behavior and instruct Fluent Bit to skip long lines and continue processing other lines that fits into the buffer size. <br> |
| DB | *string | No | - | Specify the database file to keep track of monitored files and offsets.<br> |
| DB_Sync | string | No | Full | Set a default synchronization (I/O) method. Values: Extra, Full, Normal, Off. This flag affects how the internal SQLite engine do synchronization to disk, for more details about each option please refer to this section. <br> |
| Mem_Buf_Limit | string | No | - | Set a limit of memory that Tail plugin can use when appending data to the Engine. If the limit is reach, it will be paused; when the data is flushed it resumes.<br> |
| Parser | string | No | - | Specify the name of a parser to interpret the entry as a structured message.<br> |
| Key | string | No | log | When a message is unstructured (no parser applied), it's appended as a string under the key name log. This option allows to define an alternative name for that key. <br> |
| Tag | string | No | - | Set a tag (with regex-extract fields) that will be placed on lines read.<br> |
| Tag_Regex | string | No | - | Set a regex to extract fields from the file.<br> |
| Multiline | string | No | Off | If enabled, the plugin will try to discover multiline messages and use the proper parsers to compose the outgoing messages. Note that when this option is enabled the Parser option is not used. <br> |
| Multiline_Flush | string | No | 4 | Wait period time in seconds to process queued multiline messages <br> |
| Parser_Firstline | string | No | - | Name of the parser that machs the beginning of a multiline message. Note that the regular expression defined in the parser must include a group name (named capture)<br> |
| Parser_N | []string | No | - | Optional-extra parser to interpret and structure multiline entries. This option can be used to define multiple parsers, e.g: Parser_1 ab1,  Parser_2 ab2, Parser_N abN.<br> |
| Docker_Mode | string | No | Off | If enabled, the plugin will recombine split Docker log lines before passing them to any parser as configured above. This mode cannot be used at the same time as Multiline. <br> |
| Docker_Mode_Flush | string | No | 4 | Wait period time in seconds to flush queued unfinished split lines. <br> |
### FilterKubernetes
#### FilterKubernetes Fluent Bit Kubernetes Filter allows to enrich your log files with Kubernetes metadata.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| Match | string | No | kubernetes.* | Match filtered records (default:kube.*)<br> |
| Buffer_Size | string | No | 32k | Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the Unit Size specification. <br> |
| Kube_URL | string | No | https://kubernetes.default.svc:443 | API Server end-point (default:https://kubernetes.default.svc:443)<br> |
| Kube_CA_File | string | No | /var/run/secrets/kubernetes.io/serviceaccount/ca.crt | CA certificate file (default:/var/run/secrets/kubernetes.io/serviceaccount/ca.crt)<br> |
| Kube_CA_Path | string | No | - | Absolute path to scan for certificate files<br> |
| Kube_Token_File | string | No | /var/run/secrets/kubernetes.io/serviceaccount/token | Token file  (default:/var/run/secrets/kubernetes.io/serviceaccount/token)<br> |
| Kube_Tag_Prefix | string | No | kubernetes.var.log.containers | When the source records comes from Tail input plugin, this option allows to specify what's the prefix used in Tail configuration. (default:kube.var.log.containers.)<br> |
| Merge_Log | string | No | On | When enabled, it checks if the log field content is a JSON string map, if so, it append the map fields as part of the log structure. (default:Off)<br> |
| Merge_Log_Key | string | No | - | When Merge_Log is enabled, the filter tries to assume the log field from the incoming message is a JSON string message and make a structured representation of it at the same level of the log field in the map. Now if Merge_Log_Key is set (a string name), all the new structured fields taken from the original log content are inserted under the new key.<br> |
| Merge_Log_Trim | string | No | On | When Merge_Log is enabled, trim (remove possible \n or \r) field values.  <br> |
| Merge_Parser | string | No | - | Optional parser name to specify how to parse the data contained in the log key. Recommended use is for developers or testing only.<br> |
| Keep_Log | string | No | On | When Keep_Log is disabled, the log field is removed from the incoming message once it has been successfully merged (Merge_Log must be enabled as well). <br> |
| tls.debug | string | No | -1 | Debug level between 0 (nothing) and 4 (every detail). <br> |
| tls.verify | string | No | On | When enabled, turns on certificate validation when connecting to the Kubernetes API server. <br> |
| Use_Journal | string | No | Off | When enabled, the filter reads logs coming in Journald format. <br> |
| Regex_Parser | string | No | - | Set an alternative Parser to process record Tag and extract pod_name, namespace_name, container_name and docker_id. The parser must be registered in a parsers file (refer to parser filter-kube-test as an example).<br> |
| K8S-Logging.Parser | string | No | Off | Allow Kubernetes Pods to suggest a pre-defined Parser (read more about it in Kubernetes Annotations section) <br> |
| K8S-Logging.Exclude | string | No | Off | Allow Kubernetes Pods to exclude their logs from the log processor (read more about it in Kubernetes Annotations section). <br> |
| Labels | string | No | On | Include Kubernetes resource labels in the extra metadata. <br> |
| Annotations | string | No | On | Include Kubernetes resource annotations in the extra metadata. <br> |
| Kube_meta_preload_cache_dir | string | No | - | If set, Kubernetes meta-data can be cached/pre-loaded from files in JSON format in this directory, named as namespace-pod.meta<br> |
| Dummy_Meta | string | No | Off | If set, use dummy-meta data (for test/dev purposes) <br> |
### FilterAws
#### FilterAws The AWS Filter Enriches logs with AWS Metadata.

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| imds_version | string | No | v2 | Specify which version of the instance metadata service to use. Valid values are 'v1' or 'v2' (default).<br> |
| Match | string | No | * | Match filtered records (default:*)<br> |
### VolumeMount
#### VolumeMount defines source and destination folders of a hostPath type pod mount

| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| source | string | Yes | - | Source folder<br> |
| destination | string | Yes | - | Destination Folder<br> |
| readOnly | bool | No | - | Mount Mode<br> |
