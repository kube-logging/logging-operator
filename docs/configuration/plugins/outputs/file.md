---
title: File
weight: 200
generated_file: true
---

# [File Output](https://docs.fluentd.org/output/file)
## Overview
 This plugin has been designed to output logs or metrics to File.

## Configuration
## FileOutputConfig

### path (string, required) {#fileoutputconfig-path}

The Path of the file. The actual path is path + time + ".log" by default. 

Default: -

### append (bool, optional) {#fileoutputconfig-append}

The flushed chunk is appended to existence file or not. The default is not appended. 

Default: -

### add_path_suffix (*bool, optional) {#fileoutputconfig-add_path_suffix}

Add path suffix(default: true) 

Default: true

### path_suffix (string, optional) {#fileoutputconfig-path_suffix}

The suffix of output result. 

Default:  ".log"

### symlink_path (bool, optional) {#fileoutputconfig-symlink_path}

Create symlink to temporary buffered file when buffer_type is file. This is useful for tailing file content to check logs. 

Default:  false

### compress (string, optional) {#fileoutputconfig-compress}

Compresses flushed files using gzip. No compression is performed by default. 

Default: -

### recompress (bool, optional) {#fileoutputconfig-recompress}

Performs compression again even if the buffer chunk is already compressed.  

Default:  false

### format (*Format, optional) {#fileoutputconfig-format}

[Format](../format/) 

Default: -

### buffer (*Buffer, optional) {#fileoutputconfig-buffer}

[Buffer](../buffer/) 

Default: -

### slow_flush_log_threshold (string, optional) {#fileoutputconfig-slow_flush_log_threshold}

The threshold for chunk flush performance check. Parameter type is float, not time, default: 20.0 (seconds) If chunk flush takes longer time than this threshold, fluentd logs warning message and increases metric fluentd_output_status_slow_flush_count. 

Default: -


 ## Example `File` output configurations
 ```yaml
 apiVersion: logging.banzaicloud.io/v1beta1
 kind: Output
 metadata:

	name: demo-output

 spec:

	file:
	  path: /tmp/logs/${tag}/%Y/%m/%d.%H.%M
	  append: true
	  buffer:
	    timekey: 1m
	    timekey_wait: 10s
	    timekey_use_utc: true

 ```

 #### Fluentd Config Result
 ```

	 <match **>
		@type file
		@id test_file
		add_path_suffix true
		append true
		path /tmp/logs/${tag}/%Y/%m/%d.%H.%M
		<buffer tag,time>
		  @type file
		  path /buffers/test_file.*.buffer
		  retry_forever true
		  timekey 1m
		  timekey_use_utc true
		  timekey_wait 30s
		</buffer>
	 </match>

 ```

---
