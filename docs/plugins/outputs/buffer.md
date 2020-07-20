---
title: Buffer
weight: 200
generated_file: true
---

### Buffer
| Variable Name | Type | Required | Default | Description |
|---|---|---|---|---|
| type | string | No | - | Fluentd core bundles memory and file plugins. 3rd party plugins are also available when installed.<br> |
| tags | string | No |  tag,time | When tag is specified as buffer chunk key, output plugin writes events into chunks separately per tags. <br> |
| path | string | No |  operator generated | The path where buffer chunks are stored. The '*' is replaced with random characters. It's highly recommended to leave this default. <br> |
| chunk_limit_size | string | No | - | The max size of each chunks: events will be written into chunks until the size of chunks become this size<br> |
| chunk_limit_records | int | No | - | The max number of events that each chunks can store in it<br> |
| total_limit_size | string | No | - | The size limitation of this buffer plugin instance. Once the total size of stored buffer reached this threshold, all append operations will fail with error (and data will be lost)<br> |
| queue_limit_length | int | No | - | The queue length limitation of this buffer plugin instance<br> |
| chunk_full_threshold | string | No | - | The percentage of chunk size threshold for flushing. output plugin will flush the chunk when actual size reaches chunk_limit_size * chunk_full_threshold (== 8MB * 0.95 in default)<br> |
| queued_chunks_limit_size | int | No | - | Limit the number of queued chunks. If you set smaller flush_interval, e.g. 1s, there are lots of small queued chunks in buffer. This is not good with file buffer because it consumes lots of fd resources when output destination has a problem. This parameter mitigates such situations.<br> |
| compress | string | No | - | If you set this option to gzip, you can get Fluentd to compress data records before writing to buffer chunks.<br> |
| flush_at_shutdown | bool | No | - | The value to specify to flush/write all buffer chunks at shutdown, or not<br> |
| flush_mode | string | No | - | Default: default (equals to lazy if time is specified as chunk key, interval otherwise)<br>lazy: flush/write chunks once per timekey<br>interval: flush/write chunks per specified time via flush_interval<br>immediate: flush/write chunks immediately after events are appended into chunks<br> |
| flush_interval | string | No | - | Default: 60s<br> |
| flush_thread_count | int | No | - | The number of threads of output plugins, which is used to write chunks in parallel<br> |
| flush_thread_interval | string | No | - | The sleep interval seconds of threads to wait next flush trial (when no chunks are waiting)<br> |
| flush_thread_burst_interval | string | No | - | The sleep interval seconds of threads between flushes when output plugin flushes waiting chunks next to next<br> |
| delayed_commit_timeout | string | No | - | The timeout seconds until output plugin decides that async write operation fails<br> |
| overflow_action | string | No | - | How output plugin behaves when its buffer queue is full<br>throw_exception: raise exception to show this error in log<br>block: block processing of input plugin to emit events into that buffer<br>drop_oldest_chunk: drop/purge oldest chunk to accept newly incoming chunk<br> |
| retry_timeout | string | No | - | The maximum seconds to retry to flush while failing, until plugin discards buffer chunks<br> |
| retry_forever | *bool | No | true | If true, plugin will ignore retry_timeout and retry_max_times options and retry flushing forever<br> |
| retry_max_times | int | No | - | The maximum number of times to retry to flush while failing<br> |
| retry_secondary_threshold | string | No | - | The ratio of retry_timeout to switch to use secondary while failing (Maximum valid value is 1.0)<br> |
| retry_type | string | No | - | exponential_backoff: wait seconds will become large exponentially per failures<br>periodic: output plugin will retry periodically with fixed intervals (configured via retry_wait)<br> |
| retry_wait | string | No | - | Seconds to wait before next retry to flush, or constant factor of exponential backoff<br> |
| retry_exponential_backoff_base | string | No | - | The base number of exponential backoff for retries<br> |
| retry_max_interval | string | No | - | The maximum interval seconds for exponential backoff between retries while failing<br> |
| retry_randomize | bool | No | - | If true, output plugin will retry after randomized interval not to do burst retries<br> |
| disable_chunk_backup | bool | No | - | Instead of storing unrecoverable chunks in the backup directory, just discard them. This option is new in Fluentd v1.2.6.<br> |
| timekey | string | Yes | 10m | Output plugin will flush chunks per specified time (enabled when time is specified in chunk keys)<br> |
| timekey_wait | string | No | 10m | Output plugin writes chunks after timekey_wait seconds later after timekey expiration<br> |
| timekey_use_utc | bool | No | - | Output plugin decides to use UTC or not to format placeholders using timekey<br> |
| timekey_zone | string | No | - | The timezone (-0700 or Asia/Tokyo) string for formatting timekey placeholders<br> |
