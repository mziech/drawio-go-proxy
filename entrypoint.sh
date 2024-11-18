#!/bin/sh

set -eu

mkdir -p /webroot
if [ ! -f /webroot/index.html ]; then
  echo "Downloading and installing draw.io ..."
  mkdir -p /tmp/unpack
  wget "https://github.com/jgraph/drawio/archive/$DRAWIO_BRANCH.zip"
  unzip "$DRAWIO_BRANCH.zip"
  cd "./drawio-$DRAWIO_BRANCH/src/main/webapp"
  mv * /webroot/
  rm -r /webroot/WEB-INF /webroot/META-INF
  cd /webroot
  rm -r /tmp/unpack
else
  echo "draw.io already installed."
fi

exec /app/drawio-go "$@"
