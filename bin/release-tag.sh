#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"
source "${SCRIPT_DIR}/release-status.sh"

if [[ "$CURRENT_VERSION" == "$NEXT_VERSION" ]]; then
    echo "Nothing to tag."
    exit 0
fi

if [[ "${CI:-}" != "true" ]]; then
    # When running locally, prompt before creating the tag (Safety First™).
    if ! gum confirm --default=false "Create and push tag $NEXT_VERSION"; then
        echo "Aborting."
        exit 0
    fi
fi

# Copy build artifacts over to docs and commit.
mkdir -p \
    build/docs \
    build/schemas \
    docs
cp build/docs/*.md docs/
cp build/schemas/*.json docs/
git add docs/
if [[ $(git status --porcelain 2>/dev/null) != "" ]]; then
    git commit --gpg-sign --message "chore(release): $NEXT_VERSION [skip ci]"
    git push origin main
fi

# Tag and push.
git tag \
    --sign "$NEXT_VERSION" \
    --message "$NEXT_VERSION"
git push origin --tags
