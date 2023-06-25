#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

usage() {
    echo "Usage: $(basename "$0") <project>"
}
export project="${1?$(usage)}"

mkdir -p build
go build -o "build/$project" .

rm -rf "build/${project}.schema.json"
mkdir -p build

# Generate JSON schema.
"build/$project" schema >"build/${project}.schema.json"
