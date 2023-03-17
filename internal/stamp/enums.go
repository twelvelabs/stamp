package stamp

//go:generate go-enum -f=$GOFILE --marshal --names

// ConflictConfig determines what to do when destination paths already exist.
// ENUM(keep, replace, prompt).
type ConflictConfig string

// MissingConfig determines what to do when destination paths are missing.
// ENUM(ignore, error).
type MissingConfig string
