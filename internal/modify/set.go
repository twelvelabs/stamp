package modify

import "fmt"

func NewSet(items ...any) *Set {
	s := &Set{
		keys:  []string{},
		items: map[string]any{},
	}
	s.Add(items...)
	return s
}

type Set struct {
	keys  []string
	items map[string]any
}

// Add adds an item to the set.
func (s *Set) Add(items ...any) {
	for _, item := range items {
		key := fmt.Sprintf("%#v", item)
		if _, found := s.items[key]; !found {
			s.keys = append(s.keys, key)
		}
		s.items[key] = item
	}
}

func (s *Set) Contains(item any) bool {
	key := fmt.Sprintf("%#v", item)
	if _, found := s.items[key]; found {
		return true
	}
	return false
}

func (s *Set) Len() int {
	return len(s.keys)
}

// ToSlice returns a slice containing all items in the set
// (in the order they were added).
func (s *Set) ToSlice() []any {
	items := []any{}
	for _, key := range s.keys {
		items = append(items, s.items[key])
	}
	return items
}
