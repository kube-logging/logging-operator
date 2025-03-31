#!/bin/sh

[ -z "$BUFFER_PATH" ] && BUFFER_PATH=/buffers

while true; do
    # Deprecated: node_buffer_size_bytes will soon be replaced by logging_buffer_size_bytes
    # logging_buffer_size_bytes includes the host label
    echo "# HELP node_buffer_size_bytes Disk space used [deprecated]" > /prometheus/node_exporter/textfile_collector/buffer_size.prom.$$
    echo "# TYPE node_buffer_size_bytes gauge" >> /prometheus/node_exporter/textfile_collector/buffer_size.prom.$$
    du -sb ${BUFFER_PATH} | sed -ne 's/\\/\\\\/;s/"/\\"/g;s/^\([0-9]\+\)\t\(.*\)$/node_buffer_size_bytes{entity="\2"} \1/p' >> /prometheus/node_exporter/textfile_collector/buffer_size.prom.$$
    mv /prometheus/node_exporter/textfile_collector/buffer_size.prom.$$ /prometheus/node_exporter/textfile_collector/buffer_size.prom

    echo "# HELP logging_buffer_size_bytes Disk space used" > /prometheus/node_exporter/textfile_collector/logging_buffer_size_bytes.prom.$$
    echo "# TYPE logging_buffer_size_bytes gauge" >> /prometheus/node_exporter/textfile_collector/logging_buffer_size_bytes.prom.$$
    du -sb ${BUFFER_PATH} | sed -ne 's/\\/\\\\/;s/"/\\"/g;s/^\([0-9]\+\)\t\(.*\)$/logging_buffer_size_bytes{entity="\2", host="'$(hostname)'"} \1/p' >> /prometheus/node_exporter/textfile_collector/logging_buffer_size_bytes.prom.$$
    mv /prometheus/node_exporter/textfile_collector/logging_buffer_size_bytes.prom.$$ /prometheus/node_exporter/textfile_collector/logging_buffer_size_bytes.prom

    echo "# HELP logging_buffer_files File count" > /prometheus/node_exporter/textfile_collector/logging_buffer_files.prom.$$
    echo "# TYPE logging_buffer_files gauge" >> /prometheus/node_exporter/textfile_collector/logging_buffer_files.prom.$$
    echo -e "$(find "${BUFFER_PATH}" -type f 2>/dev/null | wc -l)\t${BUFFER_PATH}" | sed -ne 's/\\/\\\\/;s/"/\\"/g;s/^\([0-9]\+\)\t\(.*\)$/logging_buffer_files{entity="\2", host="'$(hostname)'"} \1/p' >> /prometheus/node_exporter/textfile_collector/logging_buffer_files.prom.$$
    mv /prometheus/node_exporter/textfile_collector/logging_buffer_files.prom.$$ /prometheus/node_exporter/textfile_collector/logging_buffer_files.prom

    sleep 15
done
