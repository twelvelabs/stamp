# MissingConfig

Determines what to do when updating an existing file and
the destination path is missing.

> [!IMPORTANT]
> Only used in [update] and [delete] tasks.

[update]: https://github.com/twelvelabs/stamp/tree/main/docs/update_task.md
[delete]: https://github.com/twelvelabs/stamp/tree/main/docs/delete_task.md

Allowed Values:

- `"ignore"`: Do nothing. The task becomes a noop.
- `"touch"`: Create an empty file.
- `"error"`: Raise an error.
