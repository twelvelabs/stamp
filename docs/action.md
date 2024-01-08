# Action

The action to perform on the destination.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`merge`](#merge) | string | ➖ | ✅ | `"concat"` | <p>Determines merge behavior for arrays - either when modifying them directly or when recursively merging objects containing arrays. |
| [`type`](#type) | string | ➖ | ✅ | `"replace"` | <p>Determines what type of modification to perform. |

### `merge`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"concat"` |

Determines merge behavior for arrays - either when modifying them directly
or when recursively merging objects containing arrays.

Allowed Values:

- `"concat"`: Concatenate source and destination arrays.
- `"upsert"`: Add source array items if not present in the destination.
- `"replace"`: Replace the destination with the source.

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"replace"` |

Determines what type of modification to perform.

The append/prepend behavior differs slightly depending on
the destination content type. Strings are concatenated,
numbers are added, and objects are recursively merged.
Arrays are concatenated by default, but that behavior can
be customized via the 'merge' enum.

Replace and delete behave consistently across all types.

Allowed Values:

- `"append"`: Append to the destination content.
- `"prepend"`: Prepend to the destination content.
- `"replace"`: Replace the destination.
- `"delete"`: Delete the destination content.
