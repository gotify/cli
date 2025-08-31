ARG BUILDKIT_SBOM_SCAN_CONTEXT=true
ARG GO_VERSION=PLEASE_PROVIDE_GO_VERSION
ARG BUILDPLATFORM

FROM --platform=${BUILDPLATFORM} docker.io/golang:${GO_VERSION}-alpine AS builder

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH
ARG LD_FLAGS

WORKDIR /src
COPY . /src

RUN GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    GOARM=$(echo $TARGETPLATFORM | sed -En 's|^linux/arm/v(.*)|\1|p') \
    go build -mod=readonly -a -installsuffix cgo -ldflags "${LD_FLAGS}" -o /target/app/gotify-cli

FROM docker.io/library/alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /target/app/gotify-cli /gotify-cli

ENTRYPOINT ["/gotify-cli"]
