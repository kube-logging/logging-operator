FROM ghcr.io/banzaicloud/custom-runner:0.1.4 as custom-runner
FROM quay.io/prometheus/node-exporter:master
WORKDIR /
COPY --from=custom-runner /runner .
USER root
RUN mkdir -p /prometheus/node_exporter/textfile_collector
ADD buffer-size.sh /prometheus/buffer-size.sh
RUN chmod 0744 /prometheus/buffer-size.sh
ENTRYPOINT ["/runner"]
