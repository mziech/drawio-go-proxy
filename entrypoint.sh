#!/bin/sh

set -eu

if [ "$(id -u)" = "0" ]; then
  user=$(stat -c '%u' /webroot)
  if [ "$user" = "0" ]; then
    echo "Running as root because /webroot is owned by root. This is NOT recommended!"
  else
    echo "Changing to UID $user"
    exec su "$user" -c "$0" "$@"
  fi
else
  echo "Running with UID $(id -u)"
fi

if [ -z "${DRAWIO_VERSION:-}" ]; then
  echo "Determining latest version of draw.io ..."
  DRAWIO_VERSION=$(wget -q -O - https://raw.githubusercontent.com/jgraph/drawio/refs/heads/dev/VERSION)
fi

current_version=
if [ -f /webroot/VERSION ]; then
  current_version=$(cat /webroot/VERSION)
  echo "draw.io $current_version is installed."
else
  echo "No draw.io version installed."
fi

if [ "$current_version" != "$DRAWIO_VERSION" ]; then
  echo "Downloading and installing draw.io $DRAWIO_VERSION ..."
  rm -rf /tmp/unpack
  mkdir -p /tmp/unpack
  cd /tmp/unpack
  wget "https://github.com/jgraph/drawio/archive/v$DRAWIO_VERSION.zip"
  unzip "v$DRAWIO_VERSION.zip"
  cd "./drawio-$DRAWIO_VERSION"
  cp -fr src/main/webapp/* /webroot/ || true
  rm -r /webroot/WEB-INF /webroot/META-INF || true
  cd /webroot
  echo "$DRAWIO_VERSION" > VERSION
  rm -r /tmp/unpack
else
  echo "No update needed."
fi

echo "Running app"
exec /app/drawio-go "$@"
