# PromptConfig

Determines when a [value] should prompt for input.

[value]: https://github.com/twelvelabs/stamp/tree/main/docs/value.md

Allowed Values:

- `"always"`: Always prompt.
- `"never"`: Never prompt.
- `"on-empty"`: Only when input OR default is blank/zero.
- `"on-unset"`: Only when not explicitly set via CLI.
