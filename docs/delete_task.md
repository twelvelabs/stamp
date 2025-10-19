# DeleteTask

Deletes a path in the destination directory.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`dst`](#dst) | [Destination](destination.md#destination) | ✅ | ➖ | ➖ | <p>The destination path. |
| [`each`](#each) | string | ➖ | ➖ | ➖ |  |
| [`if`](#if) | string | ➖ | ➖ | `"true"` |  |
| [`type`](#type) | string | ✅ | ✅ | `"delete"` | <p>Deletes a path in the destination directory. |

### `dst`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Destination](destination.md#destination) | ✅ | ➖ | ➖ |

The destination path.

### `each`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

### `if`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | `"true"` |

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ✅ | `"delete"` |

Deletes a path in the destination directory.

Allowed Values:

- `"delete"`
