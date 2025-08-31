SHELL := /bin/bash

# Env var inputs for container image builds
# REGISTRY: The registry to which the build image should be pushed to.
# IMAGE: The name of the image to build and publish in the afore mentioned registry.
# PLATFORM: The platform for which the image should be built
REGISTRY ?= docker.io
IMAGE ?= gotify/cli
PLATFORM ?= linux/amd64,linux/arm64,linux/386,linux/arm/v7,linux/riscv64

# Env var inputs for all builds
# VERSION: The version for which the container image or the binary is being built.
#          When it is not provided, no version will be specified in the built package.
# COMMIT: The commit of this project for which the cli is being built, for reference in the tool's "version" command.
# LD_FLAGS: Build flags, for the tool's "version" command.
COMMIT ?= $(shell git rev-parse --verify HEAD)
LD_FLAGS ?= $(if $(VERSION),-X main.Version=${VERSION}) \
	-X main.BuildDate=$(shell date "+%F-%T") \
	-X main.Commit=${COMMIT}

ifdef GOTOOLCHAIN
	GO_VERSION=$(GOTOOLCHAIN)
else
	GO_VERSION=$(shell go mod edit -json | jq -r .Toolchain | sed -e 's/go//')
endif

build-docker-multiarch:
	docker buildx build \
		$(if $(DOCKER_BUILD_PUSH),--push) \
		-t ${REGISTRY}/${IMAGE}:master \
		$(if $(VERSION),-t ${REGISTRY}/${IMAGE}:latest) \
		$(if $(VERSION),-t ${REGISTRY}/${IMAGE}:${VERSION}) \
		$(if $(VERSION),-t ${REGISTRY}/${IMAGE}:$(shell echo $(VERSION) | cut -d '.' -f -2)) \
		$(if $(VERSION),-t ${REGISTRY}/${IMAGE}:$(shell echo $(VERSION) | cut -d '.' -f -1)) \
		--build-arg GO_VERSION=$(GO_VERSION) \
		--build-arg LD_FLAGS="$(LD_FLAGS)" \
		--platform $(PLATFORM) \
		-f docker/Dockerfile .

clean:
	rm -rf build

build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64       go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-windows-amd64.exe cli.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386         go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-windows-386.exe   cli.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64       go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-linux-amd64       cli.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=386         go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-linux-386         cli.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64       go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-linux-arm64       cli.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm GOARM=7 go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-linux-arm-7       cli.go
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64       go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-darwin-amd64      cli.go
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64       go build -ldflags="${LD_FLAGS}" -o build/gotify-cli-darwin-arm64      cli.go
