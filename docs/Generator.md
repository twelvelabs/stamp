# Generator

Stamp generator metadata.

## Generator Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `description` | string |  |  | The generator description. The first line is shown when listing all generators. The full description is used when viewing generator help/usage text. |
| `name` | string |  |  | The generator name. |
| `tasks` | [Task](#task)[] &#124; null |  |  |  |
| `values` | [Value](#value)[] &#124; null |  |  |  |

## CreateTask

Creates a new file in the destination directory.

### CreateTask Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `dst` | [Destination](#destination) |  |  | The destination path. |
| `each` | string |  |  | Set to a comma separated value and the task will be executued once per-item. On each iteration, the `_Item` and `_Index` values will be set accordingly. |
| `if` | string |  | true | Determines whether the task should be executed. The value must be coercible to a boolean. |
| `src` | [Source](#source) |  |  | The source path or inline content. |
| `type` | string |  |  | Creates a new file in the destination directory.<br><br>Allowed value:<br>• create |

## DeleteTask

Deletes a file in the destination directory.

### DeleteTask Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `dst` | [Destination](#destination) |  |  | The destination path. |
| `each` | string |  |  | Set to a comma separated value and the task will be executued once per-item. On each iteration, the `_Item` and `_Index` values will be set accordingly. |
| `if` | string |  | true | Determines whether the task should be executed. The value must be coercible to a boolean. |
| `type` | string |  |  | Deletes a file in the destination directory.<br><br>Allowed value:<br>• delete |

## Destination

The destination path.

### Destination Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `conflict` | string |  | prompt | ConflictConfig determines what to do when destination paths already exist.<br><br>Allowed values:<br>• keep<br>• replace<br>• prompt |
| `content_type` | string |  |  | An explicit content type. Inferred from the file extension by default.<br><br>Allowed values:<br>• json<br>• yaml<br>• text |
| `missing` | string |  | ignore | MissingConfig determines what to do when destination paths are missing.<br><br>Allowed values:<br>• ignore<br>• touch<br>• error |
| `mode` | string |  | 0666 | An optional POSIX file mode to set on the file path. |
| `path` | string |  |  | The file path relative to the destination directory. Attempts to traverse outside the destination directory will raise a runtime error |

## GeneratorTask

Executes another generator.

### GeneratorTask Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `each` | string |  |  | Set to a comma separated value and the task will be executued once per-item. On each iteration, the `_Item` and `_Index` values will be set accordingly. |
| `if` | string |  | true | Determines whether the task should be executed. The value must be coercible to a boolean. |
| `name` | string |  |  |  |
| `type` | string |  |  | Executes another generator.<br><br>Allowed value:<br>• generator |
| `values` | object &#124; null |  | map[] |  |

## Source

The source path or inline content.

### Source Variants

- [Source Content](#source-content)
- [Source Path](#source-path)

## Source Content

The source content.

### Source Content Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `content` |  | ✅ |  | Inline content. Can be any type. String keys and/or values will be rendered as templates. |
| `content_type` | string |  |  | FileType specifies the content type of the destination path.<br><br>Allowed values:<br>• json<br>• yaml<br>• text |

## Source Path

The source path.

### Source Path Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `content_type` | string |  |  | FileType specifies the content type of the destination path.<br><br>Allowed values:<br>• json<br>• yaml<br>• text |
| `path` | string | ✅ |  | The file path relative to the source directory. Attempts to traverse outside the source directory will raise a runtime error. |

## Task

A generator task.

### Task Variants

- [CreateTask](#createtask)
- [DeleteTask](#deletetask)
- [GeneratorTask](#generatortask)
- [UpdateTask](#updatetask)

## Action

The action to perform on the destination.

### Action Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `merge` | string |  | concat | MergeType determines slice merge behavior.<br><br>Allowed values:<br>• concat<br>• upsert<br>• replace |
| `type` | string |  | replace | Action determines what type of modification to perform.<br><br>Allowed values:<br>• append<br>• prepend<br>• replace<br>• delete |

## Match

Target a subset of the destination to update.

### Match Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `default` |  |  |  | A default value to use if the JSON path expression is not found. |
| `pattern` | string |  |  | A regexp or JSON path expression. |
| `source` | string |  | line | MatchSource determines whether match patterns should be applied per-line or to the entire file.<br><br>Allowed values:<br>• file<br>• line |

## UpdateTask

Updates a file in the destination directory.

### UpdateTask Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `action` | [Action](#action) |  |  | The action to perform on the destination. |
| `description` | string |  |  | An optional description of what is being updated. |
| `dst` | [Destination](#destination) |  |  | The destination path. |
| `each` | string |  |  | Set to a comma separated value and the task will be executued once per-item. On each iteration, the `_Item` and `_Index` values will be set accordingly. |
| `if` | string |  | true | Determines whether the task should be executed. The value must be coercible to a boolean. |
| `match` | [Match](#match) |  |  | Target a subset of the destination to update. |
| `src` | [Source](#source) |  |  | The source path or inline content. |
| `type` | string |  |  | Updates a file in the destination directory.<br><br>Allowed value:<br>• update |

## Value

A generator input value.

### Value Properties

| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
| `default` |  |  |  |  |
| `flag` | string |  |  |  |
| `help` | string |  |  |  |
| `if` | string |  | true |  |
| `key` | string |  |  |  |
| `mode` | string |  | flag | InputMode determines whether the value is a flag or positional argument.<br><br>Allowed values:<br>• arg<br>• flag<br>• hidden |
| `name` | string |  |  |  |
| `options` | array &#124; null |  | [] |  |
| `prompt` | string |  | on-unset | PromptConfig determines when a value should prompt.<br><br>Allowed values:<br>• always<br>• never<br>• on-empty<br>• on-unset |
| `transform` | string |  |  |  |
| `type` | string |  | string | DataType is the data type of a value.<br><br>Allowed values:<br>• bool<br>• int<br>• intSlice<br>• string<br>• stringSlice |
| `validate` | string |  |  |  |
