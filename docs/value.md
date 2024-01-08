# Value

A generator input value.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`default`](#default) |  | ➖ | ➖ | ➖ | <p>The value default. |
| [`flag`](#flag) | string | ➖ | ➖ | ➖ | <p>The flag name for the value. |
| [`help`](#help) | string | ➖ | ➖ | ➖ | <p>Help text describing the value. |
| [`if`](#if) | string | ➖ | ➖ | `"true"` | <p>Determines whether the value is enabled. |
| [`key`](#key) | string | ➖ | ➖ | ➖ | <p>The variable name for the value. |
| [`mode`](#mode) | string | ➖ | ✅ | `"flag"` | <p>Determines how the [value] can be set. |
| [`name`](#name) | string | ➖ | ➖ | ➖ | <p>The display name shown when prompting for the value. |
| [`options`](#options) | array | ➖ | ➖ | ➖ | <p>A fixed set of valid options for the value. |
| [`prompt`](#prompt) | string | ➖ | ✅ | `"on-unset"` | <p>Determines when a [value] should prompt for input. |
| [`transform`](#transform) | string | ➖ | ➖ | ➖ | <p>Optional, comma-separated list of transform rules. |
| [`type`](#type) | string | ➖ | ✅ | `"string"` | <p>Specifies the data type of a [value] |
| [`validate`](#validate) | string | ➖ | ➖ | ➖ | <p>Optional, comma-separated list of validation rules. |

### `default`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
|  | ➖ | ➖ | ➖ |

The value default. Can refer to other values defined earlier in the list.

Examples:

```yaml
default: '{{ .OtherValue | underscore }}.txt'
```

### `flag`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

The flag name for the value. Will default to a [dash separated](https://pkg.go.dev/github.com/gobuffalo/flect#Dasherize) form of the [key](https://github.com/twelvelabs/stamp/tree/main/docs/value.md#key).

Examples:

```yaml
flag: custom-flag-name
```

### `help`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

Help text describing the value. Shown when prompting and when using the `--help` flag.

Examples:

```yaml
help: You should enter a random value.
```

### `if`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | `"true"` |

Determines whether the value is enabled. Can refer to other values defined earlier in the list (allows for dynamic prompts).

Examples:

```yaml
if: '{{ .UseDatabase }}'
```

```yaml
if: '{{ eq .Language "python" }}'
```

### `key`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

The variable name for the value. This is how you will refer to the value in template files.

Examples:

```yaml
key: MyValue
```

### `mode`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"flag"` |

Determines how the [value] can be set.

[value]: https://github.com/twelvelabs/stamp/tree/main/docs/value.md

Allowed Values:

- `"arg"`: Can be set via positional argument OR prompt.
- `"flag"`: Can be set via flag OR prompt.
- `"hidden"`: Can only be set via user config.

### `name`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

The display name shown when prompting for the value. Will default to a [humanized](https://pkg.go.dev/github.com/gobuffalo/flect#Humanize) form of the [key](https://github.com/twelvelabs/stamp/tree/main/docs/value.md#key).

Examples:

```yaml
name: Custom display name
```

### `options`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| array | ➖ | ➖ | ➖ |

A fixed set of valid options for the value. Will cause the value to be rendered as a single or multi-select when prompted (depending on data type). Attempts to assign a value not in this list will raise a validation error.

Examples:

```yaml
options:
    - foo
    - bar
```

### `prompt`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"on-unset"` |

Determines when a [value] should prompt for input.

Allowed Values:

- `"always"`: Always prompt.
- `"never"`: Never prompt.
- `"on-empty"`: Only when input OR default is blank/zero.
- `"on-unset"`: Only when not explicitly set via CLI.

### `transform`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

Optional, comma-separated list of [transform](https://github.com/twelvelabs/stamp/tree/main/docs/transform.md) rules.

Examples:

```yaml
transform: trim,uppercase
```

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"string"` |

Specifies the data type of a [value].

Allowed Values:

- `"bool"`: Boolean.
- `"int"`: Integer.
- `"intSlice"`: Integer array/slice.
- `"string"`: String.
- `"stringSlice"`: String array/slice.

### `validate`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

Optional, comma-separated list of [validation](https://github.com/go-playground/validator#baked-in-validations) rules.

Examples:

```yaml
validate: required,email
```
