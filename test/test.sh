#!/bin/sh

set -eu

curl -f -v 'http://drawio:8080/proxy?url=http://www.example.com/index.html'
