## Flow control with durability in a multi tenant setup

Resources:
- https://docs.fluentbit.io/manual/administration/backpressure
- https://docs.fluentbit.io/manual/administration/buffering-and-storage
- https://docs.fluentbit.io/manual/pipeline/inputs/tail#sqlite-and-write-ahead-logging
- https://docs.fluentbit.io/manual/administration/monitoring
- https://docs.fluentbit.io/manual/administration/troubleshooting#dump-internals-signal

### Context

Let's consider we have multiple separate inputs, each sending data to their respective dedicated outputs (using tenant ids in the tags).

### Durability

According to the referenced resources we need `storage.type  filesystem` for *every input* 
where we want to avoid losing data. If we just enable this option, there will be no limit 
on how many data fluent-bit should keep on disk.

> Note: we also have to configure the position db to avoid fluent-bit 
> reading the same files from the beginning after a restart

### Memory limit

The limit that is applied by default is `storage.max_chunks_up 128` on the *service* which is a global limit.
But this only means, that even if fluent-bit writes all chunks to disk, there is a limit on how many
chunks it can read up and handle in memory at the same time.
Without any further configuration fluent-bit will write chunks to disk indefinitely and this setting will only
affect the overall throughput.

### Disk usage limit

In case we want to limit the actual disk usage we need to set `storage.total_limit_size` for
every *output* individually. This sounds good, but the problem with this option is that it doesn't
cause any backpressure, rather just starts to discard the oldest data, which obviously results in data loss,
so this option should be used with care.

### Backpressure

Backpressure can be enabled using `storage.pause_on_chunks_overlimit on` on the *input* which is great, but one important
caveat again: the limit this setting considers as the trigger event is `storage.max_chunks_up` which is a global limit. 

Going back to our main scenario, when one of the outputs is down (tenant is down), chunks for that output start to pile up
on disk and in memory. When there are more than `storage.max_chunks_up` chunks in memory globally, fluent-bit pauses inputs that
tries to load additional chunks. It's not clear how fluent-bit decides which output should be paused, but based on our 
observations (using `config/samples/multitenant-routing` for example) this works as expected as only the input that belongs
to the faulty output is paused and when the output gets back online, the input is resumed immediately.

Also based on fluent-bit's metrics, if an output is permanently down, the chunks that are waiting for that output to be sent
are not kept in memory, so other input/output pairs are not limited by the throughput. 

In case we configure `storage.pause_on_chunks_overlimit` in the inputs we can make sure the disk usage is bounded.

As long as pods are not restarting, the backpressure can prevent log loss, but keep in mind, that since the input is paused,
data in log files that gets deleted by the container runtime during the output's downtime will get lost.
