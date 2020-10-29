# Build /go/bin/smartctl_exporter
FROM quay.io/prometheus/golang-builder:1.13-base AS builder
ADD . /go/src/github.com/Sheridan/smartctl_exporter/
RUN cd /go/src/github.com/Sheridan/smartctl_exporter && make

# Container image
FROM ubuntu:18.04
WORKDIR /
RUN apt-get update \
    && apt-get install smartmontools/bionic-backports -y --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/Sheridan/smartctl_exporter/bin/smartctl_exporter /bin/smartctl_exporter
COPY docker-entrypoint.sh /
COPY smartctl_exporter.yaml /
RUN chmod +x /docker-entrypoint.sh
EXPOSE 9633
ENTRYPOINT ["/docker-entrypoint.sh"]
