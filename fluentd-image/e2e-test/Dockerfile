# Workaround for https://gitlab.alpinelinux.org/alpine/aports/-/issues/12821
FROM alpine:3.17

ENV FLUENTD_VERSION 1.15.3

LABEL Description="Fluentd docker image" Vendor="Kube logging" Version="${FLUENTD_VERSION}"

# Do not split this into multiple RUN!
# Docker creates a layer for every RUN-Statement
# therefore an 'apk delete' has no effect
RUN apk update \
 && apk add --no-cache \
        ca-certificates \
        ruby ruby-irb ruby-etc \
        tini libmaxminddb geoip \
 && buildDeps=" \
      make gcc g++ libc-dev \
      wget bzip2 zlib-dev git linux-headers \
      automake autoconf libtool build-base \
      ruby-dev libc6-compat geoip-dev \
    " \
 && apk add --no-cache --virtual .build-deps \
        build-base \
        ruby-dev gnupg \
 && apk add $buildDeps \
 && echo 'gem: --no-document' >> /etc/gemrc \
 && gem install oj -v 3.14.2 \
 && gem install http_parser.rb -v 0.8.0 \
 && gem install tzinfo -v 1.2.11 \
 && gem install resolv -v 0.2.2 \
 && gem install json -v 2.6.3 \
 && gem install yajl-ruby -v 1.4.1 \
 && gem install async-http -v 0.60.1 \
  # removed as it is not available for ruby 3
  # && gem install ext_monitor -v 0.1.2 \
 && gem install fluentd -v ${FLUENTD_VERSION} \
 && gem install prometheus-client -v 4.1.0 \
 && gem install bigdecimal -v 3.1.4 \
 && gem install webrick -v 1.8.1 \
 && gem install gelf -v 3.1.0 \
 && gem install elasticsearch -v 8.6.0 \
 && gem install elasticsearch-xpack -v 7.17.7 \
 && gem install fluent-plugin-detect-exceptions -v 0.0.15 \
 && gem install specific_install -v 0.3.8 \
 && gem install fluent-plugin-remote-syslog -v 1.1.0 \
 && gem install fluent-plugin-prometheus -v 2.0.3 \
 && gem install fluent-plugin-s3 -v 1.7.2 \
 && gem install fluent-plugin-rewrite-tag-filter -v 2.4.0 \
 && gem install fluent-plugin-dedot_filter -v 1.0.0 \
 && gem install fluent-plugin-sumologic_output -v 1.8.0 \
 && gem install fluent-plugin-kafka -v 0.18.1 \
 && gem install fluent-plugin-geoip -v 1.3.2 \
 && gem install fluent-plugin-tag-normaliser -v 0.1.2 \
 && gem install fluent-plugin-grafana-loki -v 1.2.20 \
 && gem install fluent-plugin-concat -v 2.5.0 \
 && gem install fluent-plugin-kinesis -v 3.4.2 \
 && gem install fluent-plugin-parser-logfmt -v 0.0.0 \
 && gem install fluent-plugin-multi-format-parser -v 1.0.0 \
 && gem install fluent-plugin-record-modifier -v 2.1.1 \
 && gem install fluent-plugin-splunk-hec -v 1.3.2 \
 && gem install fluent-plugin-elasticsearch -v 5.2.5 \
 && gem install fluent-plugin-opensearch -v 1.0.10 \
 && gem install fluent-plugin-aws-elasticsearch-service -v 2.4.1 \
 && gem install fluent-plugin-kube-events-timestamp -v 0.1.3 \
 && gem install fluent-plugin-grok-parser -v 2.6.2 \
  # main branch
 && gem specific_install -l https://github.com/acquia/fluent-plugin-syslog_rfc5424.git --ref 5199be67d1a385f529fa1d6b6023e95ba7fac27d \
  # master branch
 && gem specific_install -l https://github.com/banzaicloud/fluent-plugin-gcs.git --ref 9fbff74734ddce289f9a10734a338b3c7f0065b5 \
  # rwlock branch
 && gem specific_install -l https://github.com/kube-logging/fluent-plugin-label-router.git --ref 2ff43789f895735b7852e0a1a809280e22d9e8ef \
  # main branch
 && gem specific_install -l https://github.com/SumoLogic/sumologic-kubernetes-fluentd.git --ref 3a5a2e29e3776107a1bf0c9435de5cc2b14fc7f6 -d fluent-plugin-kubernetes-metadata-filter \
 && gem specific_install -l https://github.com/SumoLogic/sumologic-kubernetes-fluentd.git --ref 3a5a2e29e3776107a1bf0c9435de5cc2b14fc7f6 -d fluent-plugin-enhance-k8s-metadata \
 && gem specific_install -l https://github.com/SumoLogic/sumologic-kubernetes-fluentd.git --ref 3a5a2e29e3776107a1bf0c9435de5cc2b14fc7f6 -d fluent-plugin-kubernetes-sumologic \
 # https://github.com/SumoLogic/sumologic-kubernetes-fluentd/blob/main/Dockerfile#L35
 #&& gem specific_install https://github.com/abonas/kubeclient --ref 831e360772c717aab5ca086521c45c86ee51435e \
  # master branch
 && gem specific_install https://github.com/abonas/kubeclient --ref 6b40464cd71e2789ea6de94f618f793a8db9ac49 \
 && apk del .build-deps $buildDeps \
 && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem /usr/lib/ruby/gems/2.*/gems/fluentd-*/test

RUN addgroup -S fluent && adduser -S -G fluent fluent \
    # for log storage (maybe shared with host)
    && mkdir -p /fluentd/log \
    # configuration/plugins path (default: copied from .)
    && mkdir -p /fluentd/etc /fluentd/plugins \
    && chown -R fluent /fluentd && chgrp -R fluent /fluentd \
    && mkdir -p /buffers && chown -R fluent /buffers

COPY fluent.conf /fluentd/etc/
COPY entrypoint.sh /bin/
COPY healthy.sh /bin/

ENV FLUENTD_CONF="fluent.conf"
ENV RUBYLIB="/usr/lib/ruby/gems/3.1.0/gems/resolv-0.2.2/lib"

ENV LD_PRELOAD=""
EXPOSE 24224 5140

USER fluent
ENTRYPOINT ["tini",  "--", "/bin/entrypoint.sh"]
CMD ["fluentd"]
