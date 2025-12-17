FROM ghcr.io/kube-logging/custom-runner:v0.14.0 AS custom-runner

FROM quay.io/prometheus/node-exporter:v1.10.2

COPY --from=custom-runner /runner /

USER root

RUN mkdir -p /prometheus/node_exporter/textfile_collector && \
    chown -R nobody:nobody /prometheus

COPY buffer-size.sh /prometheus/buffer-size.sh
RUN chmod 0755 /prometheus/buffer-size.sh && \
    chown nobody:nobody /prometheus/buffer-size.sh

USER nobody:nobody

WORKDIR /

ENTRYPOINT ["/runner"]
