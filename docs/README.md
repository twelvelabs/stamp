# Stamp documentation

This is... a bit of a work in progress 😬

## Adding a generator

```bash
# Install a local generator
stamp add ./path/to/generator
# Install a remote generator from a repo
stamp add github.com/some-user/some-repo
# Install a remote generator from an archive
stamp add https://example.com/generator.tar.gz
```

## Running a generator

```bash
stamp new some:generator:name
```

## Creating a generator

TODO

## Values

TODO

## Task types

- [create](./tasks/create.md) - creates a file
- [update](./tasks/update.md) - updates a file
- [delete](./tasks/delete.md) - delete a file
- [generator](./tasks/generator.md) - delegate to another generator
