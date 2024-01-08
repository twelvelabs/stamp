# SourcePath

The source path.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`content_type`](#content_type) | string | ➖ | ✅ | ➖ | <p>Specifies the content type of the file. |
| [`path`](#path) | string | ✅ | ➖ | ➖ | <p>The file path relative to the generator source directory (\src) |

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

### `path`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ➖ | ➖ |

The file path relative to the generator source directory (\_src). Attempts to traverse outside the source directory will raise a runtime error.

The file will be rendered as a Go [text/template](https://pkg.go.dev/text/template) and have access to all the [values](value.md) defined by the generator.

The file _may_ be parsed depending on it's [content type](#content_type). Note that this happens post-render. This allows for dynamic data sources when creating/updating structured files.
