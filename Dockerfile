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
# npmrepo   Local NPM repository
#
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

# ===================================================================
# Container with JDK11, node & basic tools
ARG prefix
FROM ${prefix}openjdk:11 AS jdk
ARG aptrepo
ARG npmrepo

# --build-arg aptrepo= domain/path to replace url in apt for using a local nexus3 apt proxy repository
RUN if [ ! -z "${aptrepo}" ]; then \
    sed -i \
      -e "s|http://deb.debian.org/|${aptrepo}|" \
      -e "s|http://security.debian.org/|${aptrepo}|" \
      /etc/apt/sources.list \
    ;fi

# --build-arg npmrepo= path to a local npm repository
RUN if [ ! -z "${npmrepo}" ]; then echo registry=${npmrepo} >~/.npmrc ;fi

RUN apt-get update &&\
    apt-get install -y nodejs npm &&\
    echo;echo "Ignore npm errors whilst we upgrade it to the newest version" &&\
    npm install npm@latest -g &&\
    echo;echo "Purging APT caches" &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

FROM jdk AS npm

WORKDIR /work
COPY package.json .
RUN npm install
