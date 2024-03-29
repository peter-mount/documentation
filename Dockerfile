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
ARG npmrepo
ARG nodejs=https://nodejs.org/dist/v18.12.1/node-v18.12.1-linux-x64.tar.xz
# Needed post node15 otherwise npm will install at / & thats no longer permitted
WORKDIR /root

# We need ca-certificates installed first before our local apt proxy defined in the aptrepo build-arg
RUN if [ ! -z "${aptrepo}" ]; then \
        apt-get update &&\
        apt-get install -y ca-certificates &&\
        sed -i \
          -e "s|http://deb.debian.org/|${aptrepo}|" \
          -e "s|http://security.debian.org/|${aptrepo}|" \
          /etc/apt/sources.list \
    ;fi

RUN if [ ! -z "${npmrepo}" ]; then echo registry=${npmrepo} >~/.npmrc ;fi &&\
    apt-get update &&\
    apt-get install -y ca-certificates chromium git curl xz-utils texlive texlive-latex-extra

# Install latest stable nodejs
RUN curl ${nodejs} -o /tmp/node.txz && (cd /usr/local;tar xJpf /tmp/node.txz) &&\
    (cd /usr/local/bin;for i in ../node-*/bin/*;do ln -s ${i};done) &&\
    npm install npm@latest -g &&\
    npm install -g postcss postcss-cli &&\
    npm install autoprefixer

# cleanup
RUN chmod -R +rx /root &&\
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* ~/.npmrc

# Required to ensure that nodejs can resolve the home directory. Relies on uid being 1000
RUN echo "user:x:1000:1000::/home:/bash" >>/etc/passwd && mkdir -p /home && chown 1000:1000 /home

# ===================================================================
# Retrieve hugo (precompiled) & build our tools
ARG prefix
ARG arch=amd64
ARG goos=linux
FROM ${prefix}golang:alpine AS build
#ARG hugoVersion=0.89.2
ARG hugoVersion=0.108.0

# Required commands for the build
RUN apk add --no-cache tzdata

# Get Hugo extended build with version in the hugoVersion build-arg
RUN cd /tmp &&\
    wget -q -O /tmp/hugo.tgz https://github.com/gohugoio/hugo/releases/download/v${hugoVersion}/hugo_extended_${hugoVersion}_Linux-64bit.tar.gz &&\
    tar xpf hugo.tgz &&\
    mkdir -pv /dest/usr/local/bin/ &&\
    cp -p hugo /dest/usr/local/bin/

WORKDIR /work

# Seems go 1.16+ have changed things so this is required otherwise modules are not handled correctly with go.sum breaking
# via https://github.com/golang/go/issues/44129#issuecomment-860060061
RUN go env -w GOFLAGS=-mod=mod

# Download go module dependencies
COPY go.mod .
RUN go mod download

# intermediates that don't get detected by go.mod
# Review these every so often
RUN go get github.com/mailru/easyjson \
           github.com/chromedp/sysutil \
           golang.org/x/net \
           github.com/aymerick/douceur \
           github.com/gobwas/ws \
           github.com/xuri/efp \
           golang.org/x/crypto \
           golang.org/x/text \
           github.com/richardlehane/mscfb \
           github.com/mohae/deepcopy \
           github.com/gorilla/css \
           github.com/josharian/intern \
           github.com/gobwas/pool \
           github.com/gobwas/httphead \
           github.com/richardlehane/msoleps

# Build our tools
COPY tools/ tools/
RUN CGO_ENABLED=0 go build -o /dest/usr/local/bin/doctool tools/bin/main.go

# ===================================================================
# Final stage add all resources to base image
FROM base AS final

WORKDIR /work

# Copy the go tool build
COPY --from=build /dest/ /
