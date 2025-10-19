# CreateTask

Creates a new path in the destination directory.

When using source templates, the [src.path](source_path.md#path)
attribute may be a file or a directory path. When the latter,
the source directory will be copied to the destination path recursively.

Examples:

```yaml
tasks:
  - type: create
    # Render <./_src/README.tpl> (using the values defined in the generator)
    # and write it to <./README.md> in the destination directory.
    # If the README file already exists in the destination dir,
    # keep the existing file and do not bother prompting the user.
    src:
      path: "README.tpl"
    dst:
      path: "README.md"
      conflict: keep
```

```yaml
values:
  - key: "FirstName"
    default: "Some Name"

tasks:
  - type: create
    # Render the inline content as a template and write it to
    # <./some_name/greeting.txt> in the destination directory.
    src:
      content: "Hello, {{ .FirstName }}!"
    dst:
      path: "{{ .FirstName | underscore }}/greeting.txt"
```

```yaml
tasks:
  - type: create
    # Render all the files in <./_src/scripts/> (using the values defined in the generator),
    # copy them to <./scripts/> in the destination directory, then make them executable.
    src:
      path: "scripts/"
    dst:
      path: "scripts/"
      mode: "0755"
```

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`dst`](#dst) | [Destination](destination.md#destination) | ✅ | ➖ | ➖ | <p>The destination path. |
| [`each`](#each) | string | ➖ | ➖ | ➖ |  |
| [`if`](#if) | string | ➖ | ➖ | `"true"` |  |
| [`src`](#src) | [Source](source.md#source) | ✅ | ➖ | ➖ | <p>The source path or inline content. |
| [`type`](#type) | string | ✅ | ✅ | `"create"` | <p>Creates a new path in the destination directory. |

### `dst`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Destination](destination.md#destination) | ✅ | ➖ | ➖ |

The destination path.

Examples:

```yaml
dst:
    path: README.md
```

```yaml
dst:
    mode: "0755"
    path: bin/build.sh
```

### `each`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

### `if`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | `"true"` |

### `src`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Source](source.md#source) | ✅ | ➖ | ➖ |

The source path or inline content.

Examples:

```yaml
src:
    path: README.tpl
```

```yaml
src:
    content: Hello, {{ .FirstName }}!
```

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ✅ | `"create"` |

Creates a new path in the destination directory.

Allowed Values:

- `"create"`

Examples:

```yaml
type: create
```
