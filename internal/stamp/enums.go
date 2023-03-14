package stamp

//go:generate go-enum -f=$GOFILE --marshal --names

// ConflictConfig determines what to do when destination paths already exist.
// ENUM(keep, replace, prompt).
type ConflictConfig string
