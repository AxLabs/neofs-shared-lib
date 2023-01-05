#!/usr/bin/env bash

# This script extracts the content of the top-level section that corresponds to
# the requested release version specified by `$VERSION`.

sed '/^# '"${VERSION}"'$/,/^# /!d;//d;' CHANGELOG.md
