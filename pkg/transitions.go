package notes

import (
	"fmt"
)

type Transitions map[string]*Category

func (ts Transitions) Set(name string, category *Category) {
	ts[name] = category
}

func (ts Transitions) Find(name string) (*Category, error) {
	c, ok := ts[name]
	if !ok {
		return &Category{}, fmt.Errorf("invalid transition")
	}
	return c, nil
}

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
