# Update Task

Updates an existing file in the destination directory.

## Schema

| field                                      | type       | default    |
| ------------------------------------------ | ---------- | ---------- |
| [`action`](#action)                        | any        |            |
| [`action.type`](#actiontype)               | enum       | replace    |
| [`action.merge`](#actionmerge)             | enum       | concat     |
| [`description`](#description)              | string     |            |
| [`dst`](#dst)                              | string     |            |
| [`file_type`](#file_type)                  | string     | _inferred_ |
| [`match`](#match)                          | any        |            |
| [`match.default`](#matchdefault)           | any        | nil        |
| [`match.pattern`](#matchpattern)           | string     | _all_      |
| [`match.source`](#matchsource)             | enum       | line       |
| [`missing`](#missing)                      | enum       | ignore     |
| [`mode`](#mode)                            | string     |            |
| [`src`](#src)                              | any        |            |

### `action`

The type of update to perform is configured via the `action` field. It has two forms:

- A "long form" version:

  ```yaml
  action:
    type: "append"
    merge: "upsert"
  ```

- And a "short form" version that just specifies the [action type](#actiontype):

  ```yaml
  action: "append"
  ```

### `action.type`

The action to perform. Can be one of:

- `prepend`: prepend source content before the destination target.
- `append`: append source content after the destination target.
- `replace` replace the destination target with the source content (default).
- `delete` delete the destination target.

The **prepend/append** behavior varies slightly depending on the destination
target [type](#file_type):

- String are concatenated.
- Numbers are added.
- Boolean values are logically conjoined (i.e. `AND` or `&&`).
- Objects and arrays are recursively merged (append) or reverse-merged (prepend).

Both the destination target and the source content must be (or be coercible to)
the same type.

### `action.merge`

Determines the logic to use when merging arrays. Can be one of:

- `concat`: Concatenate both arrays (default).
- `upsert`: Concatenate any elements of the second array not already in the first.
- `replace`: Replace the first array with the second. Useful when merging
  objects containing array values.

The following table illustrates each option:

|        | concat (default)       | upsert            | replace      |
| ------ | ---------------------- | ----------------- | ------------ |
| first  | `[foo, bar]`           | `[foo, bar]`      | `[foo, bar]` |
| second | `[bar, baz]`           | `[bar, baz]`      | `[bar, baz]` |
| result | `[foo, bar, bar, baz]` | `[foo, bar, baz]` | `[bar, baz]` |

### `description`

Optional text describing the update. Will be shown to users as part of the
`stamp new` output.

Example:

```yaml
- type: update
  dst: config.json
  description: "adding {{ .Language }} settings"
```

Assuming `.Language` is "Python", the user will see:

```text
[update] path/to/dst/config.json (adding Python settings)
```

### `dst`

The destination path to be updated. Should be specified relative to the
destination directory.

### `file_type`

The file type of the destination path. Determines the behavior of the
[match pattern](#matchpattern) and whether to encode/decode content.

Can be one of: `json`, `text`, or `yaml`. If not specified, will be inferred
from the file extension.

### `match`

An optional match pattern to target. It has two forms:

- A "long form" version:

  ```yaml
  match:
    # Match everything between the words "foo" and "bar"
    # even if it crosses multiple lines
    # (the ?s flag causes `.` to match newlines).
    pattern: "(?s)foo(.*)bar"
    # Apply the pattern to the entire file, not each line.
    source: file
  ```

- And a "short form" version that just specifies a [match pattern](#matchpattern):

  ```yaml
  # Matches single lines containing nothing but "foo"
  match: "^foo$"
  ```

### `match.default`

> Only used when the [file type](#file_type) is `json` or `yaml`.

Configures a value to set (prior to update) if the
[JSON path expression](#matchpattern) is not found.
Defaults to `nil`.

For example:

```yaml
# Would normally be a noop if config.yaml contained only `{}`,
# Because nothing matched.
- type: update
  dst: config.yaml
  match: $.tags
  action: append
  src:
    - foo

# While this would result in `{ tags: [foo] }`
- type: update
  dst: config.yaml
  match:
    pattern: $.tags
    default: []
  action: append
  src:
    - foo
```

### `match.pattern`

An optional match pattern to target in the destination path. The default
behavior matches everything.

The pattern format depends on the [file type](#file_type):

- If `dst` is a text file, then `match` should be a [Go regular expression](https://pkg.go.dev/regexp/syntax).
- If `dst` is a JSON or YAML file, then `match` should be a [JSON path expression](https://goessner.net/articles/JsonPath/).

### `match.source`

> Only used when the [file type](#file_type) is `text`.

Determines the source material the [regular expression](#matchpattern)
will search. Can be one of:

- `file`: Search the entire file at once.
- `line`: Search each line of the file individually (default).

The default makes it easier to mimic the behavior of tools like `grep` and
`sed` (where ^ and $ match the line), without having to remember to add regexp
flags to the beginning of the expression. The downside is that multi-line
operations are not possible.

To perform multi-line operations, set the source to `file` and remember to
use the `(?m)` and `(?s)` [regexp flags](https://pkg.go.dev/regexp/syntax#hdr-Syntax)
when appropriate.

### `missing`

What to do if the [destination path](#dst) is missing. Can be one of:

- `ignore`: Do nothing (default).
- `error`: Return an error.

### `mode`

The POSIX file mode (in octal notation) to set on the [destination path](#dst).
Should be supplied as a YAML string literal (not octal number).

### `src`

The new content that will be added to the [destination path](#dst). Can be
either a path relative to the generator's `_src` directory, or inline content.

Examples:

```yaml
# Template path.
- type: update
  dst: example.txt
  src: "template-for-{{ .Name }}.txt"

# Inline string
- type: update
  dst: example.txt
  src: "Content for {{ .Name }}"

# Inline data
- type: update
  dst: config.yaml
  src:
    foo: true
    bar:
      - "Content for {{ .Name }} 1"
      - "Content for {{ .Name }} 2"
```

## Examples

Append templated content to `script.sh` and ensure executable:

```yaml
- type: update
  dst: script.sh
  action: append
  src: |
    echo "Hello, {{ .Name }}"
  mode: "0755"
```

Replace all instances of "hello" or "goodbye" with a personalized greeting.
Note that you can mix and match regexp capture groups and template strings:

```yaml
- type: update
  dst: file.txt
  match: "(hello|goodbye)"
  action: replace # could be omitted since it is the default
  src: "${1}, {{ .Name }}"
```

Delete three specific lines that appear together:

```yaml
- type: update
  dst: file.txt
  match:
    pattern: "foo\nbar\nbaz\n"
    source: file
  action: delete
```

Append a new item to the `tags` array in `config.yaml`:

```yaml
- type: update
  dst: config.yaml
  match: $.tags
  action: append
  src: "{{ .TagName }}"
```

Merge the object in `overrides.yaml` into `config.yaml`. Note that:

- `overrides.yaml` is rendered as a template prior to parsing.
- Any arrays in `overrides.yaml` will replace those with the same name
  in `config.yaml` (rather than the default behavior of concatenation).

```yaml
- type: update
  dst: config.yaml
  action:
    type: append
    merge: replace
  src: overrides.yaml
```
