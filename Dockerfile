# Dockerfile for building the toolchain
#
# This Dockerfile accepts multiple build-args to support local
# repositories rather than using the public ones - mainly to speed
# up builds and limit
#
# prefix    prefix added to standard docker images, used when having
#           a local repository rather than use the public docker hub
#
# aptrepo   Replacement path to a local debian repository
#
# Hugo version to install
#
# ===================================================================
# Container with current debian, chromium & dependencies
ARG prefix
FROM ${prefix}debian:11-slim AS base
ARG aptrepo

# We need ca-certificates installed first before our local apt proxy defined in the aptrepo build-arg
RUN if [ ! -z "${aptrepo}" ]; then \
        apt-get update &&\
        apt-get install -y ca-certificates &&\
        sed -i \
          -e "s|http://deb.debian.org/|${aptrepo}|" \
          -e "s|http://security.debian.org/|${aptrepo}|" \
          /etc/apt/sources.list \
    ;fi &&\
    apt-get update &&\
    apt-get install -y ca-certificates chromium &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# ===================================================================
# Retrieve hugo (precompiled) & build our tools
ARG prefix
ARG arch=amd64
ARG goos=linux
FROM ${prefix}golang:alpine AS build
ARG hugoVersion=0.87.0

# Required commands for the build
RUN apk add --no-cache git tzdata
#RUN apk add --no-cache curl git tzdata zip

# Get Hugo extended build with version in the hugoVersion build-arg
RUN mkdir -pv /dest/bin/ &&\
    cd /tmp &&\
    wget -q -O /tmp/hugo.tgz https://github.com/gohugoio/hugo/releases/download/v${hugoVersion}/hugo_extended_${hugoVersion}_Linux-64bit.tar.gz &&\
    tar xpf hugo.tgz &&\
    cp -p hugo /dest/bin/

WORKDIR /work

# Seems go 1.16+ have changed things so this is required otherwise modules are not handled correctly with go.sum breaking
# via https://github.com/golang/go/issues/44129#issuecomment-860060061
RUN go env -w GOFLAGS=-mod=mod

# Download go module dependencies
COPY go.mod .
RUN go mod download

# Build our tools
COPY tools/ tools/
RUN CGO_ENABLED=0 go build -o /dest/bin/doctool tools/bin/main.go

# ===================================================================
# Final stage add all resources to base image
FROM base AS final

# Copy the go tool build
COPY --from=build /dest/bin/* /usr/local/bin/

WORKDIR /work
