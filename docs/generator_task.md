# GeneratorTask

Executes another generator.

All values defined by the included generator are prepended
to the including generator's values list. This allows the
including generator to redefine values if needed.

You can also optionally define values to set in the included generator.
This can be useful to prevent the user from being prompted for
values that do not make sense in your use case (see example).

Example:

```yaml
name: "python-api"

tasks:
  - type: "generator"
    # Executes the gitignore generator, pre-setting
    # the "Language" value so the user isn't prompted.
    name: "gitignore"
    values:
      Language: "python"

  # ... other tasks ...
```

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`each`](#each) | string | ➖ | ➖ | ➖ | <p>Set to a comma separated value and the task will be executued once per-item. |
| [`if`](#if) | string | ➖ | ➖ | `"true"` | <p>Determines whether the task should be executed. |
| [`name`](#name) | string | ✅ | ➖ | ➖ | <p>The name of the generator to execute. |
| [`type`](#type) | string | ✅ | ✅ | `"generator"` | <p>Executes another generator. |
| [`values`](#values) | [Values](values.md#values) &#124; null | ➖ | ➖ | `{}` | <p>Optional key/value pairs to pass to the generator. |

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

### `name`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ➖ | ➖ |

The name of the generator to execute.

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ✅ | `"generator"` |

Executes another generator.

Allowed Values:

- `"generator"`

### `values`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Values](values.md#values) &#124; null | ➖ | ➖ | `{}` |

Optional key/value pairs to pass to the generator.

Examples:

```yaml
values:
    ValueKeyOne: foo
    ValueKeyTwo: bar
```
