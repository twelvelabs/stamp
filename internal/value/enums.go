package value

//go:generate go-enum -f=$GOFILE --marshal --names

// DataType is the data type of a value.
// ENUM(bool, int, intSlice, string, stringSlice)
type DataType string

// PromptConfig determines when a value should prompt.
// ENUM(always, never, on-empty, on-unset)
type PromptConfig string

// InputMode determines whether the value is a flag or positional argument.
// ENUM(arg, flag)
type InputMode string
