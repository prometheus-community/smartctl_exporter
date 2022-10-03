ARG ARCH="amd64"
ARG OS="linux"
FROM ${ARCH}/alpine:3
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

RUN apk add smartmontools

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/smartctl_exporter /bin/smartctl_exporter

EXPOSE      9633
USER        nobody
ENTRYPOINT  [ "/bin/smartctl_exporter" ]
