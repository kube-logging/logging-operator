FROM alpine:3.9
LABEL maintainer "Fluentd developers <fluentd@googlegroups.com>"
LABEL Description="Fluentd docker image" Vendor="Fluent Organization" Version="1.9.2"

# Do not split this into multiple RUN!
# Docker creates a layer for every RUN-Statement
# therefore an 'apk delete' has no effect
RUN apk update \
 && apk add --no-cache \
        ca-certificates \
        ruby ruby-irb ruby-etc ruby-webrick \
        tini \
 && apk add --no-cache --virtual .build-deps \
        build-base linux-headers \
        ruby-dev gnupg \
 && echo 'gem: --no-document' >> /etc/gemrc \
 && gem install oj -v 3.8.1 \
 && gem install json -v 2.3.0 \
 && gem install async-http -v 0.50.0 \
 && gem install ext_monitor -v 0.1.2 \
 && gem install fluentd -v 1.9.2 \
 && gem install bigdecimal -v 1.4.4 \
 && apk del .build-deps \
 && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem

RUN addgroup -S fluent && adduser -S -g fluent fluent \
    # for log storage (maybe shared with host)
    && mkdir -p /fluentd/log \
    # configuration/plugins path (default: copied from .)
    && mkdir -p /fluentd/etc /fluentd/plugins \
    && chown -R fluent /fluentd && chgrp -R fluent /fluentd


COPY fluent.conf /fluentd/etc/
COPY entrypoint.sh /bin/


ENV FLUENTD_CONF="fluent.conf"

ENV LD_PRELOAD=""
EXPOSE 24224 5140

USER fluent
ENTRYPOINT ["tini",  "--", "/bin/entrypoint.sh"]
CMD ["fluentd"]
