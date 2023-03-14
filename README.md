# Stamp

[![build](https://github.com/twelvelabs/stamp/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/twelvelabs/stamp/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/twelvelabs/stamp/branch/main/graph/badge.svg?token=AVI3Z0Y6WE)](https://codecov.io/gh/twelvelabs/stamp)

Stamp is a CLI tool for scaffolding new projects.

Project templates are packaged as generators, and are easy to create and share
(they're just a directory with [a generator.yaml file](https://github.com/gostamp/generator-app/blob/main/generator.yaml)).

## Usage

```bash
# Show all installed generators
stamp list

# Add a generator from a local directory
stamp add ~/my/generator/dir

# Add a generator from a remote origin
# Origin can be anything supported by https://github.com/hashicorp/go-getter
stamp add github.com/gostamp/generator-app

# Run the `app` generator
stamp new app
```

## Status

Currently under active development and working towards an initial alpha release.
