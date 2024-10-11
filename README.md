# Stamp

[![build](https://github.com/twelvelabs/stamp/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/twelvelabs/stamp/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/twelvelabs/stamp/branch/main/graph/badge.svg?token=AVI3Z0Y6WE)](https://codecov.io/gh/twelvelabs/stamp)

Stamp is a CLI tool for scaffolding new projects.

Project templates are packaged as [generators](./docs/README.md), and are easy to create and share with others. 

## Installation

Choose one of the following:

- Download and manually install the latest [release](https://github.com/twelvelabs/stamp/releases/latest)
- Install with [Homebrew](https://brew.sh/) 🍺

  ```bash
  brew install twelvelabs/tap/stamp
  ```

- Install from source

  ```bash
  go install github.com/twelvelabs/stamp@latest
  ```

## Documentation

- [Getting Started](./docs/README.md)
- [Generator YAML Syntax](./docs/generator.md)

## Development

```shell
git clone git@github.com:twelvelabs/stamp.git
cd stamp

# Ensures all required dependencies are installed
# and bootstraps the project for local development.
make setup

make test
make build
make install

# Show help.
make
```
