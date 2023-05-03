package modify

//go:generate go-enum -f=$GOFILE --marshal --names

// Action determines what type of modification to perform.
// ENUM(append, prepend, replace, delete).
type Action string

// SliceMerge determines slice merge behavior.
// ENUM(concat, upsert, replace).
type SliceMerge string
