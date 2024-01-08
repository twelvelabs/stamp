package value

//go:generate go-enum -f=$GOFILE -t ../enums.tmpl --marshal --names

// Specifies the data type of a [value].
//
// [value]: https://github.com/twelvelabs/stamp/tree/main/docs/value.md
/*
	ENUM(
		bool        // Boolean.
		int         // Integer.
		intSlice    // Integer array/slice.
		string      // String.
		stringSlice // String array/slice.
	).
*/
type DataType string

// Determines when a [value] should prompt for input.
//
// [value]: https://github.com/twelvelabs/stamp/tree/main/docs/value.md
/*
	ENUM(
		always   // Always prompt.
		never    // Never prompt.
		on-empty // Only when input OR default is blank/zero.
		on-unset // Only when not explicitly set via CLI.
	).
*/
type PromptConfig string

// Determines how the [value] can be set.
//
// [value]: https://github.com/twelvelabs/stamp/tree/main/docs/value.md
/*
	ENUM(
		arg    // Can be set via positional argument OR prompt.
		flag   // Can be set via flag OR prompt.
		hidden // Can only be set via user config.
	).
*/
type InputMode string
