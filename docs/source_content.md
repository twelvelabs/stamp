# SourceContent

The source content.

## Properties

| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
| [`content`](#content) |  | ✅ | ➖ | ➖ | <p>Inline content. |

### `content`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
|  | ✅ | ➖ | ➖ |

Inline content. Can be any type. String keys and/or values will be rendered as templates.

Examples:

```yaml
content: '{{ .ValueOne }}'
```

```yaml
content:
    - '{{ .ValueOne }}'
    - '{{ .ValueTwo }}'
```

```yaml
content:
    foo: '{{ .ValueOne }}'
```
