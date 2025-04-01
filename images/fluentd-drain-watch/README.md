# Fluentd Drain Watch

Fluentd Drain Watch is a monitoring script that ensures proper shutdown of Fluentd by waiting for its RPC endpoint to become available, monitoring buffer files, and terminating custom workers when no buffers remain.

## Features

- Waits for Fluentd's RPC endpoint to be available before proceeding.
- Monitors a custom-runner HTTP endpoint.
- Ensures all buffer files are processed before exiting.
- Triggers termination of custom workers upon completion.

## Usage

Set required environment variables before running:

```sh
export BUFFER_PATH=/path/to/buffers
export CHECK_INTERVAL=60  # Optional, default is 60 seconds
export RPC_ADDRESS=127.0.0.1:24444  # Optional, default is 127.0.0.1:24444
export CUSTOM_RUNNER_ADDRESS=127.0.0.1:7357  # Optional, default is 127.0.0.1:7357
```
