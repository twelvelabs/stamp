# Task

A task to execute in the destination directory.

## Variants

- [CreateTask](create_task.md#createtask)
- [UpdateTask](update_task.md#updatetask)
- [DeleteTask](delete_task.md#deletetask)
- [GeneratorTask](generator_task.md#generatortask)

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`type`](#type) | string | ✅ | ✅ | ➖ | <p>The task type. |

### `type`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| string | ✅ | ✅ | ➖ |

The task type.

Allowed Values:

- `"create"`
- `"update"`
- `"delete"`
- `"generator"`
