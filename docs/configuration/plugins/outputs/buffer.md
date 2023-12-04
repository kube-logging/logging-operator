---
title: Buffer
weight: 200
generated_file: true
---

## Buffer

### chunk_full_threshold (string, optional) {#buffer-chunk_full_threshold}

The percentage of chunk size threshold for flushing. output plugin will flush the chunk when actual size reaches chunk_limit_size * chunk_full_threshold (== 8MB * 0.95 in default) 


### chunk_limit_records (int, optional) {#buffer-chunk_limit_records}

The max number of events that each chunks can store in it 


### chunk_limit_size (string, optional) {#buffer-chunk_limit_size}

The max size of each chunks: events will be written into chunks until the size of chunks become this size (default: 8MB) 

Default: 8MB

### compress (string, optional) {#buffer-compress}

If you set this option to gzip, you can get Fluentd to compress data records before writing to buffer chunks. 


### delayed_commit_timeout (string, optional) {#buffer-delayed_commit_timeout}

The timeout seconds until output plugin decides that async write operation fails 


### disable_chunk_backup (bool, optional) {#buffer-disable_chunk_backup}

Instead of storing unrecoverable chunks in the backup directory, just discard them. This option is new in Fluentd v1.2.6. 


### disabled (bool, optional) {#buffer-disabled}

Disable buffer section (default: false) 

Default: false,hidden

### flush_at_shutdown (bool, optional) {#buffer-flush_at_shutdown}

The value to specify to flush/write all buffer chunks at shutdown, or not 


### flush_interval (string, optional) {#buffer-flush_interval}

Default: 60s 


### flush_mode (string, optional) {#buffer-flush_mode}

Default: default (equals to lazy if time is specified as chunk key, interval otherwise) lazy: flush/write chunks once per timekey interval: flush/write chunks per specified time via flush_interval immediate: flush/write chunks immediately after events are appended into chunks 


### flush_thread_burst_interval (string, optional) {#buffer-flush_thread_burst_interval}

The sleep interval seconds of threads between flushes when output plugin flushes waiting chunks next to next 


### flush_thread_count (int, optional) {#buffer-flush_thread_count}

The number of threads of output plugins, which is used to write chunks in parallel 


### flush_thread_interval (string, optional) {#buffer-flush_thread_interval}

The sleep interval seconds of threads to wait next flush trial (when no chunks are waiting) 


### overflow_action (string, optional) {#buffer-overflow_action}

How output plugin behaves when its buffer queue is full throw_exception: raise exception to show this error in log block: block processing of input plugin to emit events into that buffer drop_oldest_chunk: drop/purge oldest chunk to accept newly incoming chunk 


### path (string, optional) {#buffer-path}

The path where buffer chunks are stored. The '*' is replaced with random characters. It's highly recommended to leave this default.

Default: operator generated

### queue_limit_length (int, optional) {#buffer-queue_limit_length}

The queue length limitation of this buffer plugin instance 


### queued_chunks_limit_size (int, optional) {#buffer-queued_chunks_limit_size}

Limit the number of queued chunks. If you set smaller flush_interval, e.g. 1s, there are lots of small queued chunks in buffer. This is not good with file buffer because it consumes lots of fd resources when output destination has a problem. This parameter mitigates such situations. 


### retry_exponential_backoff_base (string, optional) {#buffer-retry_exponential_backoff_base}

The base number of exponential backoff for retries 


### retry_forever (*bool, optional) {#buffer-retry_forever}

If true, plugin will ignore retry_timeout and retry_max_times options and retry flushing forever 

Default: true

### retry_max_interval (string, optional) {#buffer-retry_max_interval}

The maximum interval seconds for exponential backoff between retries while failing 


### retry_max_times (int, optional) {#buffer-retry_max_times}

The maximum number of times to retry to flush while failing 


### retry_randomize (bool, optional) {#buffer-retry_randomize}

If true, output plugin will retry after randomized interval not to do burst retries 


### retry_secondary_threshold (string, optional) {#buffer-retry_secondary_threshold}

The ratio of retry_timeout to switch to use secondary while failing (Maximum valid value is 1.0) 


### retry_timeout (string, optional) {#buffer-retry_timeout}

The maximum seconds to retry to flush while failing, until plugin discards buffer chunks 


### retry_type (string, optional) {#buffer-retry_type}

exponential_backoff: wait seconds will become large exponentially per failures periodic: output plugin will retry periodically with fixed intervals (configured via retry_wait) 


### retry_wait (string, optional) {#buffer-retry_wait}

Seconds to wait before next retry to flush, or constant factor of exponential backoff 


### tags (*string, optional) {#buffer-tags}

When tag is specified as buffer chunk key, output plugin writes events into chunks separately per tags.

Default: tag,time

### timekey (string, required) {#buffer-timekey}

Output plugin will flush chunks per specified time (enabled when time is specified in chunk keys) 

Default: 10m

### timekey_use_utc (bool, optional) {#buffer-timekey_use_utc}

Output plugin decides to use UTC or not to format placeholders using timekey 


### timekey_wait (string, optional) {#buffer-timekey_wait}

Output plugin writes chunks after timekey_wait seconds later after timekey expiration 

Default: 1m

### timekey_zone (string, optional) {#buffer-timekey_zone}

The timezone (-0700 or Asia/Tokyo) string for formatting timekey placeholders 


### total_limit_size (string, optional) {#buffer-total_limit_size}

The size limitation of this buffer plugin instance. Once the total size of stored buffer reached this threshold, all append operations will fail with error (and data will be lost) 


### type (string, optional) {#buffer-type}

Fluentd core bundles memory and file plugins. 3rd party plugins are also available when installed. 



