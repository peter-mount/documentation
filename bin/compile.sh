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
$@
