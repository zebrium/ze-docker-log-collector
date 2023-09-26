FROM alpine:latest AS certificates
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates


FROM fluent/fluentd:v1.16

USER root

COPY --from=certificates /etc/ssl/certs/ /usr/lib/ssl/certs/

RUN gem install fluent-plugin-zebrium_output\
  docker-api\
  docker\
  fluent-plugin-rewrite-tag-filter\
  fluent-plugin-multi-format-parser\
  fluent-plugin-s3\
  fluent-plugin-record-reformer\
  fluent-plugin-concat

COPY config/fluent.conf /fluentd/etc/

USER fluent

EXPOSE 24224

ENV FLUSH_INTERVAL "60s"
ENV BUFFER_CHUNK_LIMIT_SIZE "8MB"
ENV BUFFER_CHUNK_LIMIT_RECORDS "40000"
ENV BUFFER_TOTAL_LIMIT_SIZE "1GB"
ENV BUFFER_RETRY_TIMEOUT "1h"
ENV BUFFER_RETRY_MAX_TIMES "360"
ENV BUFFER_RETRY_WAIT "10s"
ENV VERIFY_SSL "true"
ENV ZE_DEPLOYMENT_NAME "default"
ENV ZE_LOG_LEVEL "info"
ENV ZE_LOG_COLLECTOR_TYPE "docker"
ENV LOG_FORWARDER_MODE "true"