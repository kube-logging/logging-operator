# Fluentd Drain Watch

Fluentd Drain Watch is a monitoring script that ensures proper shutdown of Fluentd by waiting for its RPC endpoint to become available, monitoring buffer files, and terminating custom workers when no buffers remain.

## Features

- Waits for Fluentd's RPC endpoint to be available before proceeding.
- Monitors a custom-runner HTTP endpoint (if available).
- Handles cases where custom-runner is not deployed (e.g., when buffer metrics sidecar is disabled).
- Ensures all buffer files are processed before exiting.
- Triggers termination of custom workers upon completion.

## Usage

Set required environment variables before running:

```sh
export BUFFER_PATH=/path/to/buffers
export CHECK_INTERVAL=60  # Optional, default is 60 seconds
export RPC_ADDRESS=127.0.0.1:24444  # Optional, default is 127.0.0.1:24444
export CUSTOM_RUNNER_ADDRESS=127.0.0.1:7357  # Optional, default is 127.0.0.1:7357
export CUSTOM_RUNNER_TIMEOUT=30  # Optional, default is 30 seconds
export KILL_TIMEOUT=300  # Optional, default is 300 seconds
```

## Timeouts

### Custom Runner Timeout

The script waits for the custom-runner HTTP endpoint to become available, with a configurable timeout (default: 30 seconds). If the custom-runner is not available after the timeout, the script assumes it is not deployed (e.g., when buffer volume metrics sidecar is disabled) and continues without it. This prevents the drainer pod from hanging indefinitely when the custom-runner sidecar is not present.

### Kill Timeout

After sending the `killWorkers` signal to Fluentd, the script waits for the RPC endpoint to stop listening, confirming that Fluentd has shut down gracefully. The `KILL_TIMEOUT` (default: 300 seconds) sets the maximum time to wait for this shutdown. If Fluentd doesn't stop within this timeout, the script exits with an error code 1, preventing indefinite waiting in cases where Fluentd fails to terminate properly.
