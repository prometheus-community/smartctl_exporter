ARG ARCH="amd64"
ARG OS="linux"
FROM --platform=${OS}/${ARCH} alpine:3
ARG ARCH="amd64"
ARG OS="linux"
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

RUN apk add smartmontools

COPY .build/${OS}-${ARCH}/smartctl_exporter /bin/smartctl_exporter

EXPOSE      9633
USER        nobody
ENTRYPOINT  [ "/bin/smartctl_exporter" ]
