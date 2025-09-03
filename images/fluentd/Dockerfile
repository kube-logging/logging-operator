FROM ruby:3.4.4-alpine3.20 AS base

ARG BUILD_DEPS=" \
      make gcc g++ libc-dev \
      wget bzip2 zlib-dev git linux-headers \
      automake autoconf libtool build-base \
      ruby-dev libc6-compat geoip-dev \
      snappy-dev gnupg bash openssl-dev \
      ruby-bigdecimal \
      "

RUN addgroup -S -g 101 fluent && adduser -S -G fluent -u 100 fluent \
 # for log storage (maybe shared with host)
 && mkdir -p /fluentd/log \
 # configuration/plugins path (default: copied from .)
 && mkdir -p /fluentd/etc /fluentd/plugins \
 && chown -R fluent /fluentd && chgrp -R fluent /fluentd \
 && mkdir -p /buffers && chown -R fluent /buffers \
 && apk add --no-cache \
        ca-certificates \
        ruby ruby-irb ruby-etc \
        tini libmaxminddb geoip \
        snappy libssl3 libcrypto3 \
 && apk add --no-cache $BUILD_DEPS \
 && echo 'gem: --no-document' >> /etc/gemrc \
 && gem install fluentd -v 1.18.0 \
 && fluent-gem install specific_install -v 0.3.8 \
 && fluent-gem install fluent-plugin-label-router -v 0.5.0 \
 && fluent-gem install snappy -v 0.0.15 \
 && find /usr/local/bundle/gems/ -newer /etc/gemrc -exec chown fluent:fluent {} \; \
 && apk del $BUILD_DEPS \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

COPY fluent.conf /fluentd/etc/
COPY entrypoint.sh /bin/
COPY healthy.sh /bin/

ENV FLUENTD_CONF="fluent.conf"
ENV LD_PRELOAD=""
EXPOSE 24224 5140

USER fluent

ENTRYPOINT ["tini",  "--", "/bin/entrypoint.sh"]
CMD ["fluentd"]

### Image with all the filters
FROM base AS filters

ADD ./filters/Gemfile /Gemfile.filters
ADD ./filters/Gemfile.lock /Gemfile.filters.lock

USER root

ENV BUNDLE_WITHOUT="output,test,development"

RUN apk add --no-cache $BUILD_DEPS \
 && touch /etc/gemrc \
 && fluent-gem install --file Gemfile.filters \
 && find /usr/local/bundle/gems/ -newer /etc/gemrc -exec chown fluent:fluent {} \; \
 && apk del $BUILD_DEPS \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

USER fluent

### Image with all the filters and outputs
FROM filters AS full

ADD ./outputs/Gemfile /Gemfile.outputs
ADD ./outputs/Gemfile.lock /Gemfile.outputs.lock

USER root

RUN apk add --no-cache $BUILD_DEPS \
 && rm -rf /usr/include/fortify \
 && touch /etc/gemrc \
 && fluent-gem specific_install -l https://github.com/kube-logging/fluent-plugin-syslog_rfc5424.git --ref 4ab9f7df3757b0e31e4bc209acab05a518efdce3 \
 && fluent-gem install --file /Gemfile.outputs \
 && find /usr/local/bundle/gems/ -newer /etc/gemrc -exec chown fluent:fluent {} \; \
 && apk del $BUILD_DEPS \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

USER fluent
