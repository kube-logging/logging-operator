# Node Exporter Image

Node Exporter Image is a monitoring script that collects buffer size and file count metrics for Prometheus Node Exporter.

## Features

- Tracks disk usage of buffer files.
- Reports buffer file count.
- Generates Prometheus-compatible metrics.
- Supports a configurable buffer path.

## Usage

Set the required environment variable before running:

```sh
export BUFFER_PATH=/path/to/buffers  # Optional, default is /buffers
```

### Prometheus Integration

The script generates the following metrics for Prometheus Node Exporter:

- `node_buffer_size_bytes`: Deprecated metric for buffer disk usage.
- `logging_buffer_size_bytes`: New metric for buffer disk usage, including the host label.
- `logging_buffer_files`: Number of buffer files.

Metrics are stored in:

```sh
/prometheus/node_exporter/textfile_collector/
```

Ensure Node Exporter is configured to read from this directory.
