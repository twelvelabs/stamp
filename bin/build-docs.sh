#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

usage() {
    echo "Usage: $(basename "$0") <project>"
}
export project="${1?$(usage)}"

rm -rf build/docs
mkdir -p build/docs

# Generate markdown docs from the generated JSON schema files.
schemadoc gen --in build/schemas --out build/docs
# Ensure the markdown files are clean.
# Running twice because (sigh) some fixes appear to short circuit others
# and end up generating results that require subsequent fixes :shrug:.
markdownlint --fix ./build/docs/*.md || markdownlint --fix ./build/docs/*.md
