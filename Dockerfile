# Dockerfile for building the toolchain
# ===================================================================
# GO for the various tools
ARG prefix
ARG arch=amd64
ARG goos=linux
FROM ${prefix}golang:alpine AS gobuild
RUN apk add --no-cache \
    curl \
    git \
    tzdata \
    zip

WORKDIR /work

# Seems go 1.16+ have changed things so this is required otherwise modules are not handled correctly with go.sum breaking
# via https://github.com/golang/go/issues/44129#issuecomment-860060061
RUN go env -w GOFLAGS=-mod=mod

# Install dependencies
COPY go.mod .
RUN go mod download

# Build our tools
COPY tools/ tools/
RUN CGO_ENABLED=0 go build \
    -o /dest/bin/doctool \
    tools/bin/main.go
