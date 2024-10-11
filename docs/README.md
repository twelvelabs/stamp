# Getting Started

## Creating a generator

Create a directory containing two files:

- `generator.yaml`:

  ```yaml
  name: my-generator-name
  description: A generator that creates a text file.

  values:
    - key: Greeting
      type: string
      default: Hello

  tasks:
    - type: create
      src:
        path: message.tpl
      dst:
        path: message.txt
  ```

- `_src/message.tpl`

  ```text
  {{ .Greeting }}, World!
  ```

When installed and executed, this generator will:

- Prompt the user for a greeting (defaulting to "Hello")
- Render the template to a file named `message.txt`

See the `generator.yaml` [documentation](./generator.md) for all generator options.

## Adding a generator

```bash
# Install a local generator
stamp add ./path/to/generator
# Install a remote generator from a repo
stamp add github.com/some-user/some-repo
# Install a remote generator from an archive
stamp add https://example.com/generator.tar.gz
```

See [twelvelabs/generator-app](https://github.com/twelvelabs/generator-app) for an example generator repo.

## Running a generator

```bash
stamp new my-generator-name
```
