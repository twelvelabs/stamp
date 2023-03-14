package stamp

import "testing"

// dummy test that exercises all the generated enum methods
// so that they don't mess up our overall coverage numbers :/.
func TestConflict(t *testing.T) {
	_ = ConflictConfigNames()
	c := ConflictConfig("foo")
	_ = c.IsValid()
	_ = c.String()
	_, _ = c.MarshalText()
	_ = (&c).UnmarshalText([]byte{})
}
