---
title: FluentbitSpec
weight: 200
generated_file: true
---

## FluentbitSpec

FluentbitSpec defines the desired state of Fluentbit

### daemonsetAnnotations (map[string]string, optional) {#fluentbitspec-daemonsetannotations}

Default: -

### annotations (map[string]string, optional) {#fluentbitspec-annotations}

Default: -

### labels (map[string]string, optional) {#fluentbitspec-labels}

Default: -

### envVars ([]corev1.EnvVar, optional) {#fluentbitspec-envvars}

Default: -

### image (ImageSpec, optional) {#fluentbitspec-image}

Default: -

### tls (*FluentbitTLS, optional) {#fluentbitspec-tls}

Default: -

### targetHost (string, optional) {#fluentbitspec-targethost}

Default: -

### targetPort (int32, optional) {#fluentbitspec-targetport}

Default: -

### flush (int32, optional) {#fluentbitspec-flush}

Set the flush time in seconds.nanoseconds. The engine loop uses a Flush timeout to define when is required to flush the records ingested by input plugins through the defined output plugins. (default: 1) 

Default: 1

### grace (int32, optional) {#fluentbitspec-grace}

Set the grace time in seconds as Integer value. The engine loop uses a Grace timeout to define wait time on exit (default: 5) 

Default: 5

### logLevel (string, optional) {#fluentbitspec-loglevel}

Set the logging verbosity level. Allowed values are: error, warn, info, debug and trace. Values are accumulative, e.g: if 'debug' is set, it will include error, warning, info and debug.  Note that trace mode is only available if Fluent Bit was built with the WITH_TRACE option enabled. (default: info) 

Default: info

### coroStackSize (int32, optional) {#fluentbitspec-corostacksize}

Set the coroutines stack size in bytes. The value must be greater than the page size of the running system. Don't set too small value (say 4096), or coroutine threads can overrun the stack buffer. Do not change the default value of this parameter unless you know what you are doing. (default: 24576) 

Default: 24576

### resources (corev1.ResourceRequirements, optional) {#fluentbitspec-resources}

Default: -

### tolerations ([]corev1.Toleration, optional) {#fluentbitspec-tolerations}

Default: -

### nodeSelector (map[string]string, optional) {#fluentbitspec-nodeselector}

Default: -

### affinity (*corev1.Affinity, optional) {#fluentbitspec-affinity}

Default: -

### metrics (*Metrics, optional) {#fluentbitspec-metrics}

Default: -

### security (*Security, optional) {#fluentbitspec-security}

Default: -

### positiondb (volume.KubernetesVolume, optional) {#fluentbitspec-positiondb}

[volume.KubernetesVolume](https://github.com/banzaicloud/operator-tools/tree/master/docs/types) 

Default: -

### position_db (*volume.KubernetesVolume, optional) {#fluentbitspec-position_db}

Deprecated, use positiondb 

Default: -

### mountPath (string, optional) {#fluentbitspec-mountpath}

Default: -

### extraVolumeMounts ([]*VolumeMount, optional) {#fluentbitspec-extravolumemounts}

Default: -

### inputTail (InputTail, optional) {#fluentbitspec-inputtail}

Default: -

### filterAws (*FilterAws, optional) {#fluentbitspec-filteraws}

Default: -

### filterModify ([]FilterModify, optional) {#fluentbitspec-filtermodify}

Default: -

### parser (string, optional) {#fluentbitspec-parser}

Deprecated, use inputTail.parser 

Default: -

### filterKubernetes (FilterKubernetes, optional) {#fluentbitspec-filterkubernetes}

Parameters for Kubernetes metadata filter 

Default: -

### disableKubernetesFilter (*bool, optional) {#fluentbitspec-disablekubernetesfilter}

Disable Kubernetes metadata filter 

Default: -

### bufferStorage (BufferStorage, optional) {#fluentbitspec-bufferstorage}

Default: -

### bufferStorageVolume (volume.KubernetesVolume, optional) {#fluentbitspec-bufferstoragevolume}

[volume.KubernetesVolume](https://github.com/banzaicloud/operator-tools/tree/master/docs/types) 

Default: -

### bufferVolumeMetrics (*Metrics, optional) {#fluentbitspec-buffervolumemetrics}

Default: -

### bufferVolumeImage (ImageSpec, optional) {#fluentbitspec-buffervolumeimage}

Default: -

### bufferVolumeArgs ([]string, optional) {#fluentbitspec-buffervolumeargs}

Default: -

### customConfigSecret (string, optional) {#fluentbitspec-customconfigsecret}

Default: -

### podPriorityClassName (string, optional) {#fluentbitspec-podpriorityclassname}

Default: -

### livenessProbe (*corev1.Probe, optional) {#fluentbitspec-livenessprobe}

Default: -

### livenessDefaultCheck (bool, optional) {#fluentbitspec-livenessdefaultcheck}

Default: -

### readinessProbe (*corev1.Probe, optional) {#fluentbitspec-readinessprobe}

Default: -

### network (*FluentbitNetwork, optional) {#fluentbitspec-network}

Default: -

### forwardOptions (*ForwardOptions, optional) {#fluentbitspec-forwardoptions}

Default: -

### enableUpstream (bool, optional) {#fluentbitspec-enableupstream}

Default: -

### serviceAccount (*typeoverride.ServiceAccount, optional) {#fluentbitspec-serviceaccount}

Default: -

### dnsPolicy (corev1.DNSPolicy, optional) {#fluentbitspec-dnspolicy}

Default: -

### dnsConfig (*corev1.PodDNSConfig, optional) {#fluentbitspec-dnsconfig}

Default: -

### HostNetwork (bool, optional) {#fluentbitspec-hostnetwork}

Default: -

### syslogng_output (*FluentbitTCPOutput, optional) {#fluentbitspec-syslogng_output}

Default: -


## FluentbitTLS

FluentbitTLS defines the TLS configs

### enabled (*bool, required) {#fluentbittls-enabled}

Default: -

### secretName (string, optional) {#fluentbittls-secretname}

Default: -

### sharedKey (string, optional) {#fluentbittls-sharedkey}

Default: -


## FluentbitTCPOutput

FluentbitTCPOutput defines the TLS configs

### json_date_key (string, optional) {#fluentbittcpoutput-json_date_key}

Default: ts

### json_date_format (string, optional) {#fluentbittcpoutput-json_date_format}

Default: iso8601


## FluentbitNetwork

FluentbitNetwork defines network configuration for fluentbit

### connectTimeout (*uint32, optional) {#fluentbitnetwork-connecttimeout}

Sets the timeout for connecting to an upstream  

Default:  10

### connectTimeoutLogError (*bool, optional) {#fluentbitnetwork-connecttimeoutlogerror}

On connection timeout, specify if it should log an error. When disabled, the timeout is logged as a debug message  

Default:  true

### dnsMode (string, optional) {#fluentbitnetwork-dnsmode}

Sets the primary transport layer protocol used by the asynchronous DNS resolver for connections established  

Default:  UDP, UDP or TCP

### dnsPreferIpv4 (*bool, optional) {#fluentbitnetwork-dnspreferipv4}

Prioritize IPv4 DNS results when trying to establish a connection  

Default:  false

### dnsResolver (string, optional) {#fluentbitnetwork-dnsresolver}

Select the primary DNS resolver type  

Default:  ASYNC, LEGACY or ASYNC

### keepalive (*bool, optional) {#fluentbitnetwork-keepalive}

Whether or not TCP keepalive is used for the upstream connection  

Default:  true

### keepaliveIdleTimeout (*uint32, optional) {#fluentbitnetwork-keepaliveidletimeout}

How long in seconds a TCP keepalive connection can be idle before being recycled  

Default:  30

### keepaliveMaxRecycle (*uint32, optional) {#fluentbitnetwork-keepalivemaxrecycle}

How many times a TCP keepalive connection can be used before being recycled  

Default:  0, disabled

### sourceAddress (string, optional) {#fluentbitnetwork-sourceaddress}

Specify network address (interface) to use for connection and data traffic.  

Default:  disabled


## BufferStorage

BufferStorage is the Service Section Configuration of fluent-bit

### storage.path (string, optional) {#bufferstorage-storage.path}

Set an optional location in the file system to store streams and chunks of data. If this parameter is not set, Input plugins can only use in-memory buffering. 

Default: -

### storage.sync (string, optional) {#bufferstorage-storage.sync}

Configure the synchronization mode used to store the data into the file system. It can take the values normal or full.  

Default: normal

### storage.checksum (string, optional) {#bufferstorage-storage.checksum}

Enable the data integrity check when writing and reading data from the filesystem. The storage layer uses the CRC32 algorithm.  

Default: Off

### storage.backlog.mem_limit (string, optional) {#bufferstorage-storage.backlog.mem_limit}

If storage.path is set, Fluent Bit will look for data chunks that were not delivered and are still in the storage layer, these are called backlog data. This option configure a hint of maximum value of memory to use when processing these records.  

Default: 5M


## InputTail

InputTail defines Fluentbit tail input configuration The tail input plugin allows to monitor one or several text files. It has a similar behavior like tail -f shell command.

### storage.type (string, optional) {#inputtail-storage.type}

Specify the buffering mechanism to use. It can be memory or filesystem.  

Default: memory

### Buffer_Chunk_Size (string, optional) {#inputtail-buffer_chunk_size}

Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the Unit Size specification.  

Default: 32k

### Buffer_Max_Size (string, optional) {#inputtail-buffer_max_size}

Set the limit of the buffer size per monitored file. When a buffer needs to be increased (e.g: very long lines), this value is used to restrict how much the memory buffer can grow. If reading a file exceed this limit, the file is removed from the monitored file list. The value must be according to the Unit Size specification.  

Default: Buffer_Chunk_Size

### Path (string, optional) {#inputtail-path}

Pattern specifying a specific log files or multiple ones through the use of common wildcards. 

Default: -

### Path_Key (string, optional) {#inputtail-path_key}

If enabled, it appends the name of the monitored file as part of the record. The value assigned becomes the key in the map. 

Default: -

### Exclude_Path (string, optional) {#inputtail-exclude_path}

Set one or multiple shell patterns separated by commas to exclude files matching a certain criteria, e.g: exclude_path=*.gz,*.zip 

Default: -

### Read_From_Head (bool, optional) {#inputtail-read_from_head}

For new discovered files on start (without a database offset/position), read the content from the head of the file, not tail. 

Default: -

### Refresh_Interval (string, optional) {#inputtail-refresh_interval}

The interval of refreshing the list of watched files in seconds.  

Default: 60

### Rotate_Wait (string, optional) {#inputtail-rotate_wait}

Specify the number of extra time in seconds to monitor a file once is rotated in case some pending data is flushed.  

Default: 5

### Ignore_Older (string, optional) {#inputtail-ignore_older}

Ignores files that have been last modified before this time in seconds. Supports m,h,d (minutes, hours,days) syntax. Default behavior is to read all specified files. 

Default: -

### Skip_Long_Lines (string, optional) {#inputtail-skip_long_lines}

When a monitored file reach it buffer capacity due to a very long line (Buffer_Max_Size), the default behavior is to stop monitoring that file. Skip_Long_Lines alter that behavior and instruct Fluent Bit to skip long lines and continue processing other lines that fits into the buffer size.  

Default: Off

### DB (*string, optional) {#inputtail-db}

Specify the database file to keep track of monitored files and offsets. 

Default: -

### DB_Sync (string, optional) {#inputtail-db_sync}

Set a default synchronization (I/O) method. Values: Extra, Full, Normal, Off. This flag affects how the internal SQLite engine do synchronization to disk, for more details about each option please refer to this section.  

Default: Full

### DB.locking (*bool, optional) {#inputtail-db.locking}

Specify that the database will be accessed only by Fluent Bit. Enabling this feature helps to increase performance when accessing the database but it restrict any external tool to query the content.  

Default:  true

### DB.journal_mode (string, optional) {#inputtail-db.journal_mode}

sets the journal mode for databases (WAL). Enabling WAL provides higher performance. Note that WAL is not compatible with shared network file systems.  

Default:  WAL

### Mem_Buf_Limit (string, optional) {#inputtail-mem_buf_limit}

Set a limit of memory that Tail plugin can use when appending data to the Engine. If the limit is reach, it will be paused; when the data is flushed it resumes. 

Default: -

### Parser (string, optional) {#inputtail-parser}

Specify the name of a parser to interpret the entry as a structured message. 

Default: -

### Key (string, optional) {#inputtail-key}

When a message is unstructured (no parser applied), it's appended as a string under the key name log. This option allows to define an alternative name for that key.  

Default: log

### Tag (string, optional) {#inputtail-tag}

Set a tag (with regex-extract fields) that will be placed on lines read. 

Default: -

### Tag_Regex (string, optional) {#inputtail-tag_regex}

Set a regex to extract fields from the file. 

Default: -

### Multiline (string, optional) {#inputtail-multiline}

If enabled, the plugin will try to discover multiline messages and use the proper parsers to compose the outgoing messages. Note that when this option is enabled the Parser option is not used.  

Default: Off

### Multiline_Flush (string, optional) {#inputtail-multiline_flush}

Wait period time in seconds to process queued multiline messages  

Default: 4

### Parser_Firstline (string, optional) {#inputtail-parser_firstline}

Name of the parser that machs the beginning of a multiline message. Note that the regular expression defined in the parser must include a group name (named capture) 

Default: -

### Parser_N ([]string, optional) {#inputtail-parser_n}

Optional-extra parser to interpret and structure multiline entries. This option can be used to define multiple parsers, e.g: Parser_1 ab1,  Parser_2 ab2, Parser_N abN. 

Default: -

### Docker_Mode (string, optional) {#inputtail-docker_mode}

If enabled, the plugin will recombine split Docker log lines before passing them to any parser as configured above. This mode cannot be used at the same time as Multiline.  

Default: Off

### Docker_Mode_Parser (string, optional) {#inputtail-docker_mode_parser}

Specify an optional parser for the first line of the docker multiline mode. 

Default: -

### Docker_Mode_Flush (string, optional) {#inputtail-docker_mode_flush}

Wait period time in seconds to flush queued unfinished split lines.  

Default: 4

### multiline.parser ([]string, optional) {#inputtail-multiline.parser}

Specify one or multiple parser definitions to apply to the content. Part of the new Multiline Core support in 1.8  

Default:  ""


## FilterKubernetes

FilterKubernetes Fluent Bit Kubernetes Filter allows to enrich your log files with Kubernetes metadata.

### Match (string, optional) {#filterkubernetes-match}

Match filtered records (default:kube.*) 

Default: kubernetes.*

### Buffer_Size (string, optional) {#filterkubernetes-buffer_size}

Set the buffer size for HTTP client when reading responses from Kubernetes API server. The value must be according to the Unit Size specification. A value of 0 results in no limit, and the buffer will expand as-needed. Note that if pod specifications exceed the buffer limit, the API response will be discarded when retrieving metadata, and some kubernetes metadata will fail to be injected to the logs. If this value is empty we will set it "0".  

Default: "0"

### Kube_URL (string, optional) {#filterkubernetes-kube_url}

API Server end-point (default:https://kubernetes.default.svc:443) 

Default: https://kubernetes.default.svc:443

### Kube_CA_File (string, optional) {#filterkubernetes-kube_ca_file}

CA certificate file (default:/var/run/secrets/kubernetes.io/serviceaccount/ca.crt) 

Default: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt

### Kube_CA_Path (string, optional) {#filterkubernetes-kube_ca_path}

Absolute path to scan for certificate files 

Default: -

### Kube_Token_File (string, optional) {#filterkubernetes-kube_token_file}

Token file  (default:/var/run/secrets/kubernetes.io/serviceaccount/token) 

Default: /var/run/secrets/kubernetes.io/serviceaccount/token

### Kube_Tag_Prefix (string, optional) {#filterkubernetes-kube_tag_prefix}

When the source records comes from Tail input plugin, this option allows to specify what's the prefix used in Tail configuration. (default:kube.var.log.containers.) 

Default: kubernetes.var.log.containers

### Merge_Log (string, optional) {#filterkubernetes-merge_log}

When enabled, it checks if the log field content is a JSON string map, if so, it append the map fields as part of the log structure. (default:Off) 

Default: On

### Merge_Log_Key (string, optional) {#filterkubernetes-merge_log_key}

When Merge_Log is enabled, the filter tries to assume the log field from the incoming message is a JSON string message and make a structured representation of it at the same level of the log field in the map. Now if Merge_Log_Key is set (a string name), all the new structured fields taken from the original log content are inserted under the new key. 

Default: -

### Merge_Log_Trim (string, optional) {#filterkubernetes-merge_log_trim}

When Merge_Log is enabled, trim (remove possible \n or \r) field values.   

Default: On

### Merge_Parser (string, optional) {#filterkubernetes-merge_parser}

Optional parser name to specify how to parse the data contained in the log key. Recommended use is for developers or testing only. 

Default: -

### Keep_Log (string, optional) {#filterkubernetes-keep_log}

When Keep_Log is disabled, the log field is removed from the incoming message once it has been successfully merged (Merge_Log must be enabled as well).  

Default: On

### tls.debug (string, optional) {#filterkubernetes-tls.debug}

Debug level between 0 (nothing) and 4 (every detail).  

Default: -1

### tls.verify (string, optional) {#filterkubernetes-tls.verify}

When enabled, turns on certificate validation when connecting to the Kubernetes API server.  

Default: On

### Use_Journal (string, optional) {#filterkubernetes-use_journal}

When enabled, the filter reads logs coming in Journald format.  

Default: Off

### Cache_Use_Docker_Id (string, optional) {#filterkubernetes-cache_use_docker_id}

When enabled, metadata will be fetched from K8s when docker_id is changed.  

Default: Off

### Regex_Parser (string, optional) {#filterkubernetes-regex_parser}

Set an alternative Parser to process record Tag and extract pod_name, namespace_name, container_name and docker_id. The parser must be registered in a parsers file (refer to parser filter-kube-test as an example). 

Default: -

### K8S-Logging.Parser (string, optional) {#filterkubernetes-k8s-logging.parser}

Allow Kubernetes Pods to suggest a pre-defined Parser (read more about it in Kubernetes Annotations section)  

Default: Off

### K8S-Logging.Exclude (string, optional) {#filterkubernetes-k8s-logging.exclude}

Allow Kubernetes Pods to exclude their logs from the log processor (read more about it in Kubernetes Annotations section).  

Default: Off

### Labels (string, optional) {#filterkubernetes-labels}

Include Kubernetes resource labels in the extra metadata.  

Default: On

### Annotations (string, optional) {#filterkubernetes-annotations}

Include Kubernetes resource annotations in the extra metadata.  

Default: On

### Kube_meta_preload_cache_dir (string, optional) {#filterkubernetes-kube_meta_preload_cache_dir}

If set, Kubernetes meta-data can be cached/pre-loaded from files in JSON format in this directory, named as namespace-pod.meta 

Default: -

### Dummy_Meta (string, optional) {#filterkubernetes-dummy_meta}

If set, use dummy-meta data (for test/dev purposes)  

Default: Off

### DNS_Retries (string, optional) {#filterkubernetes-dns_retries}

DNS lookup retries N times until the network start working  

Default: 6

### DNS_Wait_Time (string, optional) {#filterkubernetes-dns_wait_time}

DNS lookup interval between network status checks  

Default: 30

### Use_Kubelet (string, optional) {#filterkubernetes-use_kubelet}

This is an optional feature flag to get metadata information from kubelet instead of calling Kube Server API to enhance the log.  

Default: Off

### Kubelet_Port (string, optional) {#filterkubernetes-kubelet_port}

kubelet port using for HTTP request, this only works when Use_Kubelet  set to On  

Default: 10250


## FilterAws

FilterAws The AWS Filter Enriches logs with AWS Metadata.

### imds_version (string, optional) {#filteraws-imds_version}

Specify which version of the instance metadata service to use. Valid values are 'v1' or 'v2' (default). 

Default: v2

### az (*bool, optional) {#filteraws-az}

The availability zone (default:true). 

Default: true

### ec2_instance_id (*bool, optional) {#filteraws-ec2_instance_id}

The EC2 instance ID. (default:true) 

Default: true

### ec2_instance_type (*bool, optional) {#filteraws-ec2_instance_type}

The EC2 instance type. (default:false) 

Default: false

### private_ip (*bool, optional) {#filteraws-private_ip}

The EC2 instance private ip. (default:false) 

Default: false

### ami_id (*bool, optional) {#filteraws-ami_id}

The EC2 instance image id. (default:false) 

Default: false

### account_id (*bool, optional) {#filteraws-account_id}

The account ID for current EC2 instance. (default:false) 

Default: false

### hostname (*bool, optional) {#filteraws-hostname}

The hostname for current EC2 instance. (default:false) 

Default: false

### vpc_id (*bool, optional) {#filteraws-vpc_id}

The VPC ID for current EC2 instance. (default:false) 

Default: false

### Match (string, optional) {#filteraws-match}

Match filtered records (default:*) 

Default: *


## FilterModify

FilterModify The Modify Filter plugin allows you to change records using rules and conditions.

### rules ([]FilterModifyRule, optional) {#filtermodify-rules}

Fluentbit Filter Modification Rule 

Default: -

### conditions ([]FilterModifyCondition, optional) {#filtermodify-conditions}

Fluentbit Filter Modification Condition 

Default: -


## FilterModifyRule

FilterModifyRule The Modify Filter plugin allows you to change records using rules and conditions.

### Set (*FilterKeyValue, optional) {#filtermodifyrule-set}

Add a key/value pair with key KEY and value VALUE. If KEY already exists, this field is overwritten 

Default: -

### Add (*FilterKeyValue, optional) {#filtermodifyrule-add}

Add a key/value pair with key KEY and value VALUE if KEY does not exist 

Default: -

### Remove (*FilterKey, optional) {#filtermodifyrule-remove}

Remove a key/value pair with key KEY if it exists 

Default: -

### Remove_wildcard (*FilterKey, optional) {#filtermodifyrule-remove_wildcard}

Remove all key/value pairs with key matching wildcard KEY 

Default: -

### Remove_regex (*FilterKey, optional) {#filtermodifyrule-remove_regex}

Remove all key/value pairs with key matching regexp KEY 

Default: -

### Rename (*FilterKeyValue, optional) {#filtermodifyrule-rename}

Rename a key/value pair with key KEY to RENAMED_KEY if KEY exists AND RENAMED_KEY does not exist 

Default: -

### Hard_rename (*FilterKeyValue, optional) {#filtermodifyrule-hard_rename}

Rename a key/value pair with key KEY to RENAMED_KEY if KEY exists. If RENAMED_KEY already exists, this field is overwritten 

Default: -

### Copy (*FilterKeyValue, optional) {#filtermodifyrule-copy}

Copy a key/value pair with key KEY to COPIED_KEY if KEY exists AND COPIED_KEY does not exist 

Default: -

### Hard_copy (*FilterKeyValue, optional) {#filtermodifyrule-hard_copy}

Copy a key/value pair with key KEY to COPIED_KEY if KEY exists. If COPIED_KEY already exists, this field is overwritten 

Default: -


## FilterModifyCondition

FilterModifyCondition The Modify Filter plugin allows you to change records using rules and conditions.

### Key_exists (*FilterKey, optional) {#filtermodifycondition-key_exists}

Is true if KEY exists 

Default: -

### Key_does_not_exist (*FilterKeyValue, optional) {#filtermodifycondition-key_does_not_exist}

Is true if KEY does not exist 

Default: -

### A_key_matches (*FilterKey, optional) {#filtermodifycondition-a_key_matches}

Is true if a key matches regex KEY 

Default: -

### No_key_matches (*FilterKey, optional) {#filtermodifycondition-no_key_matches}

Is true if no key matches regex KEY 

Default: -

### Key_value_equals (*FilterKeyValue, optional) {#filtermodifycondition-key_value_equals}

Is true if KEY exists and its value is VALUE 

Default: -

### Key_value_does_not_equal (*FilterKeyValue, optional) {#filtermodifycondition-key_value_does_not_equal}

Is true if KEY exists and its value is not VALUE 

Default: -

### Key_value_matches (*FilterKeyValue, optional) {#filtermodifycondition-key_value_matches}

Is true if key KEY exists and its value matches VALUE 

Default: -

### Key_value_does_not_match (*FilterKeyValue, optional) {#filtermodifycondition-key_value_does_not_match}

Is true if key KEY exists and its value does not match VALUE 

Default: -

### Matching_keys_have_matching_values (*FilterKeyValue, optional) {#filtermodifycondition-matching_keys_have_matching_values}

Is true if all keys matching KEY have values that match VALUE 

Default: -

### Matching_keys_do_not_have_matching_values (*FilterKeyValue, optional) {#filtermodifycondition-matching_keys_do_not_have_matching_values}

Is true if all keys matching KEY have values that do not match VALUE 

Default: -


## Operation

Operation Doc stub

### Op (string, optional) {#operation-op}

Default: -

### Key (string, optional) {#operation-key}

Default: -

### Value (string, optional) {#operation-value}

Default: -


## FilterKey

### key (string, optional) {#filterkey-key}

Default: -


## FilterKeyValue

### key (string, optional) {#filterkeyvalue-key}

Default: -

### value (string, optional) {#filterkeyvalue-value}

Default: -


## VolumeMount

VolumeMount defines source and destination folders of a hostPath type pod mount

### source (string, required) {#volumemount-source}

Source folder 

Default: -

### destination (string, required) {#volumemount-destination}

Destination Folder 

Default: -

### readOnly (*bool, optional) {#volumemount-readonly}

Mount Mode 

Default: -


## ForwardOptions

ForwardOptions defines custom forward output plugin options, see https://docs.fluentbit.io/manual/pipeline/outputs/forward

### Time_as_Integer (bool, optional) {#forwardoptions-time_as_integer}

Default: -

### Send_options (bool, optional) {#forwardoptions-send_options}

Default: -

### Require_ack_response (bool, optional) {#forwardoptions-require_ack_response}

Default: -

### Tag (string, optional) {#forwardoptions-tag}

Default: -

### Retry_Limit (string, optional) {#forwardoptions-retry_limit}

Default: -

### storage.total_limit_size (string, optional) {#forwardoptions-storage.total_limit_size}

`storage.total_limit_size` Limit the maximum number of Chunks in the filesystem for the current output logical destination. 

Default: -


