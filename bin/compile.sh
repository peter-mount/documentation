#!/bin/bash

function cleanup() {
  set +x
  echo "Fixing permissions"
  cd /work &&\
  chown -R $USERID .
}
trap cleanup EXIT

set -x
cd /work &&\
npm install &&\
doctool -6502 content/6502/ &&\
doctool -bbc content/bbc/ &&\
hugo
