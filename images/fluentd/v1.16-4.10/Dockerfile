# TODO every time the base image is changed please make sure BASEIMAGEDEPS is updated (further down in this file)
# https://gallery.ecr.aws/sumologic/kubernetes-fluentd
FROM public.ecr.aws/sumologic/kubernetes-fluentd:1.16.5-sumo-0-alpine as sumo

FROM ruby:3.2.5-alpine3.20 as base

ARG BUILD_DEPS=" \
      make gcc g++ libc-dev \
      wget bzip2 zlib-dev git linux-headers \
      automake autoconf libtool build-base \
      ruby-dev libc6-compat geoip-dev \
      snappy-dev gnupg bash openssl-dev \
      "

RUN addgroup -S -g 101 fluent && adduser -S -G fluent -u 100 fluent \
    # for log storage (maybe shared with host)
    && mkdir -p /fluentd/log \
    # configuration/plugins path (default: copied from .)
    && mkdir -p /fluentd/etc /fluentd/plugins \
    && chown -R fluent /fluentd && chgrp -R fluent /fluentd \
    && mkdir -p /buffers && chown -R fluent /buffers

COPY --from=sumo --chown=fluent:fluent /usr/local/bundle /usr/local/bundle

RUN apk add --no-cache \
        ca-certificates \
        ruby ruby-irb ruby-etc \
        tini libmaxminddb geoip \
        snappy \
 && apk add --no-cache $BUILD_DEPS \
 && apk add --no-cache "libssl3>=3.1.4-r1" \
 && apk add --no-cache "libcrypto3>=3.1.4-r1" \
 && echo 'gem: --no-document' >> /etc/gemrc \
 # BASEIMAGEDEPS
 # need to rebuild these gems to have the extensions installed
 && fluent-gem install \
         bigdecimal:1.4.4 \
         oj:3.16.1 \
         google-protobuf:3.25.5 \
         nio4r:2.7.3 \
         snappy:0.3.0 \
         webrick:1.8.2 \
 # The only required fluentd dependency is the label router
 # kubeclient install is upgrading the one from the base image
 && fluent-gem install specific_install -v 0.3.8 \
 && fluent-gem specific_install -l https://github.com/ManageIQ/kubeclient.git --ref 054bff2c5e31a555004be2b3c4d32fb9dc5e6a0f \
 && fluent-gem specific_install -l https://github.com/kube-logging/fluent-plugin-label-router.git --ref namespace-label \
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
FROM filters as full

ADD ./outputs/Gemfile /Gemfile.outputs
ADD ./outputs/Gemfile.lock /Gemfile.outputs.lock

USER root

RUN apk add --no-cache $BUILD_DEPS \
 && touch /etc/gemrc \
 # TODO get rid of this by submitting the patch upstream
 && fluent-gem specific_install -l https://github.com/kube-logging/fluent-plugin-gcs.git --ref ad6a7cb0f166e2c1648954a928e1e9b8cdcbf93d \
 && fluent-gem specific_install -l https://github.com/kube-logging/fluent-plugin-syslog_rfc5424.git --ref 4ab9f7df3757b0e31e4bc209acab05a518efdce3 \
 && fluent-gem install --file /Gemfile.outputs \
 && find /usr/local/bundle/gems/ -newer /etc/gemrc -exec chown fluent:fluent {} \; \
 && apk del $BUILD_DEPS \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

USER fluent
