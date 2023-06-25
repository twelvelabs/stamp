package modify

//go:generate go-enum -f=$GOFILE -t ../enums.tmpl --marshal --names

// Action determines what type of modification to perform.
// ENUM(append, prepend, replace, delete).
type Action string

// MergeType determines slice merge behavior.
// ENUM(concat, upsert, replace).
type MergeType string
