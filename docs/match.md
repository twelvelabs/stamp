# Match

Target a subset of the destination to update.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`default`](#default) |  | ➖ | ➖ | ➖ | <p>A default value to use if the JSON path expression is not found. |
| [`pattern`](#pattern) | string | ➖ | ➖ | `""` | <p>A regexp (content type: text) or JSON path expression (content type: json, yaml) |
| [`source`](#source) | string | ➖ | ✅ | `"line"` | <p>Determines how regexp patterns should be applied. |

### `default`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
|  | ➖ | ➖ | ➖ |

A default value to use if the JSON path expression is not found.

### `pattern`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | `""` |

A regexp (content type: text) or JSON path expression (content type: json, yaml). When empty, will match everything.

### `source`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"line"` |

Determines how regexp patterns should be applied.

Allowed Values:

- `"file"`: Match the entire file.
- `"line"`: Match each line.
