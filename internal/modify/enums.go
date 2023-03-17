package modify

//go:generate go-enum -f=$GOFILE --marshal --names

// Action determines what type of modification to perform.
// ENUM(append, prepend, replace, delete).
type Action string
