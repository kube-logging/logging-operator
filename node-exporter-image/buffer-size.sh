#!/bin/sh

while true; do

echo "# HELP node_buffer_size_bytes Disk space used" > /prometheus/node_exporter/textfile_collector/buffer_size.prom 
echo "# TYPE node_buffer_size_bytes gauge" >> /prometheus/node_exporter/textfile_collector/buffer_size.prom 

[ -z "$BUFFER_PATH" ] && BUFFER_PATH=/buffers
du -ab ${BUFFER_PATH} | sed -ne 's/\\/\\\\/;s/"/\\"/g;s/^\([0-9]\+\)\t\(.*\)$/node_buffer_size_bytes{entity="\2"} \1/p' >> /prometheus/node_exporter/textfile_collector/buffer_size.prom 

sleep 60
done
