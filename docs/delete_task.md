# DeleteTask

Deletes a path in the destination directory.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`dst`](#dst) | [Destination](destination.md#destination) | ✅ | ➖ | ➖ | <p>The destination path. |
| [`each`](#each) | string | ➖ | ➖ | ➖ | <p>Set to a comma separated value and the task will be executued once per-item. |
| [`if`](#if) | string | ➖ | ➖ | `"true"` | <p>Determines whether the task should be executed. |
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

Set to a comma separated value and the task will be executued once per-item. On each iteration, the `_Item` and `_Index` values will be set accordingly.

Examples:

```yaml
each: foo, bar, baz
```

```yaml
each: '{{ .SomeList | join "," }}'
```

### `if`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | `"true"` |

Determines whether the task should be executed. The value must be [coercible](https://pkg.go.dev/strconv#ParseBool) to a boolean.

Examples:

```yaml
if: "true"
```

```yaml
if: '{{ .SomeBool }}'
```

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ✅ | `"delete"` |

Deletes a path in the destination directory.

Allowed Values:

- `"delete"`
