#!/usr/bin/env bash
set -o errexit -o errtrace -o nounset -o pipefail

# Ensure gitlint.
# Not included in dependencies JSON because it takes a long time to install
# in CI (since it depends on python@3.x) and we only use it locally.
ensure-dependency "gitlint" "brew install --quiet gitlint"

# Ensure $USER owns /usr/local/{bin,share}.
# Allows for running `make install` w/out sudo interruptions.
if [[ ! -w /usr/local/bin ]]; then
    echo "Changing owner of /usr/local/bin to $USER."
    sudo chown -R "${USER}" /usr/local/bin
fi
if [[ ! -w /usr/local/share ]]; then
    echo "Changing owner of /usr/local/share to $USER."
    sudo chown -R "${USER}" /usr/local/share
fi

# Ensure local repo.
if ! git rev-parse --is-inside-work-tree &>/dev/null; then
    if gum confirm "Create local git repo?"; then
        echo "Creating local repo."
        git init
        git add .
        git commit -m "feat: initial commit"
    fi
fi

# Ensure local repo hooks.
if [[ -d .git ]]; then
    echo "Updating .git/hooks."
    mkdir -p .git/hooks
    rm -f .git/hooks/*.sample
    cp -f bin/githooks/* .git/hooks/
    chmod +x .git/hooks/*
fi
