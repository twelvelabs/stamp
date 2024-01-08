# FileType

Specifies the content type of the file.
Inferred from the file extension by default.

When the content type is JSON or YAML, the file will be
parsed into a data structure before use.
When updating files, the content type determines
the behavior of the [match.pattern] attribute.

[match.pattern]: https://github.com/twelvelabs/stamp/tree/main/docs/match.md#pattern

Allowed Values:

- `"json"`
- `"yaml"`
- `"text"`
