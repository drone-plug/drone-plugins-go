#!/bin/sh

set -e
go build .

set +e

set -x

DRONE=true ./drone-downstream

DRONE=true \
     PLUGIN_VERSION=asd \
     ./drone-downstream

DRONE=true \
     PLUGIN_PLUGIN_DEBUG=1 \
     PLUGIN_VERSION=1.0 \
     ./drone-downstream

DRONE=true \
     PLUGIN_TOKEN=tokenvalue \
     PLUGIN_SERVER=servervalue \
     ./drone-downstream
