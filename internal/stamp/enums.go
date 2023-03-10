package stamp

//go:generate go-enum -f=$GOFILE --marshal --names

// Conflict determines what to do when destination paths already exist.
// ENUM(keep, replace, prompt).
type Conflict string
