package notes

import (
	"fmt"
)

// Represents a collection of named transitions, each pointing to a Category.
type Transitions map[string]*Category

// Associates a named transition with a specific Category.
func (ts Transitions) Set(name string, category *Category) {
	ts[name] = category
}

// Retrieves the Category associated with a given transition name.
// Returns an error if the transition name is not found.
func (ts Transitions) Find(name string) (*Category, error) {
	c, ok := ts[name]
	if !ok {
		return &Category{}, fmt.Errorf("invalid transition")
	}
	return c, nil
}

// Populates the Transitions map from a raw map of transition names to category names.
// It looks up each category name in the provided Categories and associates it with the transition.
// Returns an error if any category name is not found.
func (ts Transitions) Init(raw map[string]string, categories *Categories) error {
	for transition, target := range raw {
		category, err := categories.Find(target)
		if err != nil {
			return err
		}
		ts.Set(transition, category)
	}
	return nil
}
