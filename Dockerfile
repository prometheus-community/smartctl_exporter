# syntax=docker/dockerfile:1.6-labs

ARG GOLANG_VERSION="1.21"
ARG BUILD_IMAGE="golang:${GOLANG_VERSION}-alpine"
ARG GOLANGCI_LINT_IMAGE="golangci/golangci-lint:latest"
ARG ALPINE_VERSION="3.19"
ARG TARGET_OS="alpine"

# =============================================================================
FROM ${BUILD_IMAGE} as base

SHELL ["/bin/sh", "-e", "-u", "-o", "pipefail", "-o", "errexit", "-o", "nounset", "-c"]

WORKDIR /src/smartctl-exporter

ARG GO111MODULE="on"
ARG CGO_ENABLED="0"
ARG GOARCH="amd64"
ARG GOOS="linux"
ENV GO111MODULE="${GO111MODULE}" \
    CGO_ENABLED="${CGO_ENABLED}"  \
    GOARCH="${GOARCH}" \
    GOOS="${GOOS}"

RUN --mount=type=bind,source=./smartctl-exporter/go.mod,target=./go.mod \
    --mount=type=bind,source=./smartctl-exporter/go.sum,target=./go.sum \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download

# =============================================================================
FROM ${GOLANGCI_LINT_IMAGE} AS lint-base

# =============================================================================
FROM base AS lint

RUN --mount=type=bind,source=./smartctl-exporter/,target=./ \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    \
    golangci-lint run \
        --color never \
        --timeout 10m0s ./... | tee /linter_result.txt

# =============================================================================
FROM base AS test

RUN --mount=type=bind,source=./smartctl-exporter/,target=./ \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test -v -coverprofile=/cover.out ./...

# =============================================================================
FROM base as build

ARG APP_VERSION="docker"

RUN --mount=type=bind,source=./smartctl-exporter/,target=./ \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    \
    go build \
        -tags musl \
        -ldflags "-X 'main.AppVersion=${APP_VERSION}'" \
        -o /app && /app --version

# wait until other stages are done
COPY --from=lint /linter_result.txt /linter_result.txt
COPY --from=test /cover.out /cover.out

# =============================================================================
FROM alpine:${ALPINE_VERSION} as alpine-release

SHELL ["/bin/ash", "-e", "-u", "-o", "pipefail", "-o", "errexit", "-o", "nounset", "-c"]
RUN <<'EOF'
apk add --no-cache smartmontools
rm -rf /var/cache/apk/*
EOF

COPY --link --from=build /app /bin/smartctl-exporter

USER nobody
ENTRYPOINT [ "/bin/smartctl-exporter" ]

# =============================================================================
FROM ${TARGET_OS}-release as release
