# MergeType

Determines merge behavior for arrays - either when modifying them directly
or when recursively merging objects containing arrays.

Allowed Values:

- `"concat"`: Concatenate source and destination arrays.
- `"upsert"`: Add source array items if not present in the destination.
- `"replace"`: Replace the destination with the source.
