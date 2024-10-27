# Action

Determines what type of modification to perform.

The append/prepend behavior differs slightly depending on
the destination content type. Strings are concatenated,
numbers are added, and objects are recursively merged.
Arrays are concatenated by default, but that behavior can
be customized via the 'merge' enum.

Replace and delete behave consistently across all types.

Allowed Values:

- `"append"`: Append to the destination content.
- `"prepend"`: Prepend to the destination content.
- `"replace"`: Replace the destination.
- `"delete"`: Delete the destination content.
