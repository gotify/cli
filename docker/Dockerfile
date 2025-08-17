ARG BUILDKIT_SBOM_SCAN_CONTEXT=true
# Suppress warning about invalid variable expansion
ARG GO_VERSION=PLEASE_PROVIDE_GO_VERSION
ARG DEBIAN=sid-slim
ARG BUILDPLATFORM
ARG TARGETPLATFORM

# Hack to normalize platform to match the chosed build image
# Get the gotify/build image tag
ARG __TARGETPLATFORM_DASHES=${TARGETPLATFORM/\//-}
ARG __TARGETPLATFORM_GO_NOTATION=${__TARGETPLATFORM_DASHES/arm\/v7/arm-7}

# --- Go Builder ---
FROM --platform=${BUILDPLATFORM} docker.io/golang:${GO_VERSION}-alpine AS builder

ARG LD_FLAGS=""

RUN apk add --no-cache ca-certificates git make bash

COPY . /src/gotify

RUN cd /src/gotify && \
    LD_FLAGS=${LD_FLAGS} make OUTPUT=/target/app/gotify-cli _build_within_docker

FROM docker.io/library/alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /target/app/gotify-cli /gotify-cli

ENTRYPOINT ["/gotify-cli"]
