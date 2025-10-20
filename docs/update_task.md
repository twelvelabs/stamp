# UpdateTask

Updates a file in the destination directory.

The default behavior is to replace the entire file with the
source content, but you can optionally specify alternate
[actions](#action) (prepend, append, or delete) or [target](#match)
a subsection of the destination file.
If the destination file is structured (JSON, YAML), then you
may target a JSON path pattern, otherwise it will be treated
as plain text and you can target via regular expression.

Examples:

```yaml
tasks:
  - type: update
    # Render <./_src/COPYRIGHT.tpl> and append it
    # to the end of the README.
    # If the README does not exist in the destination dir,
    # then do nothing.
    src:
      path: "COPYRIGHT.tpl"
    action:
      type: "append"
    dst:
      path: "README.md"
```

```yaml
tasks:
  - type: update
    # Update <./package.json> in the destination dir.
    # If the file is missing, create it.
    dst:
      path: "package.json"
      missing: "touch"
    # Don't update the entire file - just the dependencies section.
    # If the dependencies section is missing, initialize it to an empty object.
    match:
      pattern: "$.dependencies"
      default: {}
    # Append (i.e. merge) the source content to the dependencies section.
    # The default behavior is to fully replace the matched pattern
    # with the source content.
    action:
      type: "append"
    # Use this inline object as the source content.
    # We could alternately reference a source file
    # containing a JSON object.
    src:
      content:
        lodash: "4.17.21"
```

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`action`](#action) | [UpdateAction](update_action.md#updateaction) | ➖ | ➖ | ➖ | <p>The action to perform on the destination. |
| [`description`](#description) | string | ➖ | ➖ | ➖ | <p>An optional description of what is being updated. |
| [`dst`](#dst) | [Destination](destination.md#destination) | ✅ | ➖ | ➖ | <p>The destination path. |
| [`each`](#each) | string | ➖ | ➖ | ➖ | <p>Set to a comma separated value and the task will be executued once per-item. |
| [`if`](#if) | string | ➖ | ➖ | `"true"` | <p>Determines whether the task should be executed. |
| [`match`](#match) | [UpdateMatch](update_match.md#updatematch) | ➖ | ➖ | ➖ | <p>Target a subset of the destination to update. |
| [`src`](#src) | [Source](source.md#source) | ✅ | ➖ | ➖ | <p>The source path or inline content. |
| [`type`](#type) | string | ✅ | ✅ | `"update"` | <p>Updates a file in the destination directory. |

### `action`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [UpdateAction](update_action.md#updateaction) | ➖ | ➖ | ➖ |

The action to perform on the destination.

### `description`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

An optional description of what is being updated.

### `dst`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Destination](destination.md#destination) | ✅ | ➖ | ➖ |

The destination path.

### `each`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

Set to a comma separated value and the task will be executued once per-item. On each iteration, the _Item and_Index values will be set accordingly.

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

### `match`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [UpdateMatch](update_match.md#updatematch) | ➖ | ➖ | ➖ |

Target a subset of the destination to update.

### `src`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Source](source.md#source) | ✅ | ➖ | ➖ |

The source path or inline content.

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ✅ | `"update"` |

Updates a file in the destination directory.

Allowed Values:

- `"update"`
