# Build stage
FROM docker.io/golang:1.26-alpine3.23 AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=0 for static binary, -ldflags for smaller binary
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-w -s" -o smartctl_exporter .

# Runtime stage
FROM alpine:3
LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"
LABEL ARCH="amd64"
LABEL OS="linux"
RUN apk add --no-cache smartmontools

COPY --from=builder /build/smartctl_exporter /bin/smartctl_exporter

EXPOSE      9633
USER        nobody
ENTRYPOINT  [ "/bin/smartctl_exporter" ]
