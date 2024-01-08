# Destination

The destination path.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`conflict`](#conflict) | string | ➖ | ✅ | `"prompt"` | <p>Determines what to do when creating a new file and the destination path already exists. |
| [`content_type`](#content_type) | string | ➖ | ✅ | ➖ | <p>Specifies the content type of the file. |
| [`missing`](#missing) | string | ➖ | ✅ | `"ignore"` | <p>Determines what to do when updating an existing file and the destination path is missing. |
| [`mode`](#mode) | string | ➖ | ➖ | `"0666"` | <p>An optional POSIX mode to set on the file path. |
| [`path`](#path) | string | ✅ | ➖ | ➖ | <p>The file path relative to the destination directory. |

### `conflict`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"prompt"` |

Determines what to do when creating a new file and
the destination path already exists.

> [!IMPORTANT]
> Only used in [create] tasks.

[create]: https://github.com/twelvelabs/stamp/tree/main/docs/create_task.md

Allowed Values:

- `"keep"`: Keep the existing path. The task becomes a noop.
- `"replace"`: Replace the existing path.
- `"prompt"`: Prompt the user.

### `content_type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | ➖ |

Specifies the content type of the file.
Inferred from the file extension by default.

When the content type is JSON or YAML, the file will be
parsed into a data structure before use.
When updating files, the content type determines
the behavior of the [match.pattern] attribute.

[match.pattern]: https://github.com/twelvelabs/stamp/tree/main/docs/match.md#pattern

Allowed Values:

- `"json"`
- `"yaml"`
- `"text"`

### `missing`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"ignore"` |

Determines what to do when updating an existing file and
the destination path is missing.

> [!IMPORTANT]
> Only used in [update] and [delete] tasks.

[update]: https://github.com/twelvelabs/stamp/tree/main/docs/update_task.md
[delete]: https://github.com/twelvelabs/stamp/tree/main/docs/delete_task.md

Allowed Values:

- `"ignore"`: Do nothing. The task becomes a noop.
- `"touch"`: Create an empty file.
- `"error"`: Raise an error.

### `mode`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | `"0666"` |

An optional [POSIX mode](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation) to set on the file path.

Examples:

```yaml
mode: "0755"
```

```yaml
mode: '{{ .ModeValue }}'
```

### `path`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ➖ | ➖ |

The file path relative to the destination directory. Attempts to traverse outside the destination directory will raise a runtime error

When creating new files, the [conflict](#conflict) attribute will be used if the path already exists. When updating or deleting files, the [missing](#missing) attribute will be used if the path does not exist.
