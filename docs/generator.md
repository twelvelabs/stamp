# Generator

Stamp generator metadata.

Example:

```yaml
name: "greet"
description: "Generates a text file with a greeting message."

# When run, the generator will prompt the user for these values.
# Alternately, the user can pass them in via flags.
values:
  # Values are prompted in the order they are defined.
  - key: "Name"
    default: "Some Name"

  # Subsequent values can reference those defined prior.
  # This allows for sensible, derived defaults.
  - key: "Greeting"
    default: "Hello, {{ .Name }}."

# Next, the generator executes a series of tasks.
# Tasks have access to the values defined above.
tasks:
  # Render the inline content as a template string.
  # Write it to <./some_name.txt> in the destination directory.
  - type: create
    src:
      content: "{{ .Greeting }}"
    dst:
      path: "{{ .Name | underscore }}.txt"
```

```shell
# Save the above to <./greeting/generator.yaml>.
# The following will prompt for values, then write <./some_name.txt>.
stamp new ./greeting

# Pass an alternate destination dir as the second argument.
# The following creates </some/other/dir/some_name.txt>.
stamp new ./greeting /some/other/dir

# Install the generator so you can refer to it by name
# rather than filesystem path.
stamp add ./greeting
stamp new greet

# You can also publish it to a git repo or upload it as an archive
# and share it with others:
stamp add git@github.com:username/my-generator.git
stamp add github.com/username/my-generator
stamp add https://example.com/my-generator.tar.gz
```

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`description`](#description) | string | ➖ | ➖ | ➖ | <p>The generator description. |
| [`name`](#name) | string | ✅ | ➖ | ➖ | <p>The generator name. |
| [`tasks`](#tasks) | [Task](task.md#task)[] &#124; null | ➖ | ➖ | ➖ | <p>A list of generator tasks. |
| [`values`](#values) | [Value](value.md#value)[] &#124; null | ➖ | ➖ | ➖ | <p>A list of generator input values. |
| [`visibility`](#visibility) | string | ➖ | ✅ | `"public"` | <p>How the generator may be viewed or invoked. |

### `description`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ➖ | ➖ |

The generator description. The first line is shown when listing all generators. The full description is used when viewing generator help/usage text.

### `name`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ➖ | ➖ |

The generator name.

### `tasks`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Task](task.md#task)[] &#124; null | ➖ | ➖ | ➖ |

A list of generator [tasks](https://github.com/twelvelabs/stamp/tree/main/docs/task.md).

### `values`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| [Value](value.md#value)[] &#124; null | ➖ | ➖ | ➖ |

A list of generator input [values](https://github.com/twelvelabs/stamp/tree/main/docs/value.md).

### `visibility`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ➖ | ✅ | `"public"` |

How the generator may be viewed or invoked.

Allowed Values:

- `"public"`: Callable anywhere.
- `"hidden"`: Public, but hidden in the generator list.
- `"private"`: Only callable as a sub-generator. Never displayed.
