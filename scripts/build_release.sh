#!/bin/bash

set -e

NAME="albio"

ROOT=$(dirname $0)/..

rm -fR snapshot/
goxc -tasks='xc archive' -bc 'linux,amd64,!arm darwin,amd64,!arm' -d .
cp -p "$ROOT"/snapshot/linux_amd64/"$NAME" "$ROOT"/snapshot/"$NAME"_linux_amd64
cp -p "$ROOT"/snapshot/darwin_amd64/"$NAME" "$ROOT"/snapshot/"$NAME"_darwin_amd64
