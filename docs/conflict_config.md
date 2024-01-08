# ConflictConfig

Determines what to do when creating a new file and
the destination path already exists.

> [!IMPORTANT]
> Only used in [create] tasks.

[create]: https://github.com/twelvelabs/stamp/tree/main/docs/create_task.md

Allowed Values:

- `"keep"`: Keep the existing path. The task becomes a noop.
- `"replace"`: Replace the existing path.
- `"prompt"`: Prompt the user.
