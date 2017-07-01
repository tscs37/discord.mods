#!/bin/sh

set -eux

build() {
    cd ./installer
    ./build.sh
}

test() {
    echo "No tests defined so far"
}

test
build