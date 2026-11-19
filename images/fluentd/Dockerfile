FROM ruby:4.0.7-slim-trixie AS ruby-slim

RUN echo 'gem: --no-document' >> /etc/gemrc

### Build stages: compile gems with native extensions
FROM ruby-slim AS build-base

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        build-essential autoconf automake libtool pkg-config \
        git zlib1g-dev libgeoip-dev libsnappy-dev libssl-dev \
 && rm -rf /var/lib/apt/lists/* \
 && gem install fluentd -v 1.19.3 \
 && gem install bigdecimal -v 3.3.1 \
 && gem install ostruct -v 0.6.3 \
 && gem install oj -v 3.17.3 \
 && fluent-gem install specific_install -v 0.3.8 \
 && fluent-gem install fluent-plugin-label-router -v 0.5.0 \
 && fluent-gem install snappy -v 0.0.15 \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

FROM build-base AS build-filters

COPY filters/Gemfile /Gemfile.filters
COPY filters/Gemfile.lock /Gemfile.filters.lock

ENV BUNDLE_WITHOUT="output,test,development"

RUN fluent-gem install --file /Gemfile.filters \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

FROM build-filters AS build-full

COPY outputs/Gemfile /Gemfile.outputs
COPY outputs/Gemfile.lock /Gemfile.outputs.lock

RUN fluent-gem specific_install -l https://github.com/kube-logging/fluent-plugin-syslog_rfc5424.git --ref 4ab9f7df3757b0e31e4bc209acab05a518efdce3 \
 && fluent-gem specific_install -l https://github.com/kube-logging/fluent-plugin-elasticsearch.git --ref f711c4322e37852a83caab0385ba4567031eb9ab \
 && fluent-gem install --file /Gemfile.outputs \
 && rm -rf /usr/local/bundle/cache/* && find /usr/local/bundle -name '*.o' -delete

### Runtime stages: no compiler, only the gems from the matching build stage
FROM ruby-slim AS runtime

RUN groupadd --system --gid 101 fluent \
 && useradd --system --gid fluent --uid 100 --create-home --shell /usr/sbin/nologin fluent \
 && mkdir -p /fluentd/log /fluentd/etc /fluentd/plugins /buffers \
 && chown -R fluent:fluent /fluentd /buffers \
 && apt-get update \
 && apt-get install -y --no-install-recommends \
        ca-certificates tini libjemalloc2 \
        libgeoip1t64 libmaxminddb0 libsnappy1v5 libcurl4t64 \
 && rm -rf /var/lib/apt/lists/*

COPY fluent.conf /fluentd/etc/
COPY entrypoint.sh healthy.sh /bin/

ENV FLUENTD_CONF="fluent.conf"
ENV LD_PRELOAD="libjemalloc.so.2"
EXPOSE 24224 5140

USER fluent

ENTRYPOINT ["tini", "--", "/bin/entrypoint.sh"]
CMD ["fluentd"]

FROM runtime AS base

COPY --from=build-base --chown=fluent:fluent /usr/local/bundle /usr/local/bundle

### Image with all the filters
FROM runtime AS filters

ENV BUNDLE_WITHOUT="output,test,development"

COPY --from=build-filters --chown=fluent:fluent /usr/local/bundle /usr/local/bundle

### Image with all the filters and outputs
FROM runtime AS full

ENV BUNDLE_WITHOUT="output,test,development"

COPY --from=build-full --chown=fluent:fluent /usr/local/bundle /usr/local/bundle
