package notes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rwxrob/fs/dir"
)

// Defines a category of notes.
type Category struct {
	Title       string
	Name        string
	Weight      int
	Path        string
	Transitions *Transitions
	Notes       *Notes
	Hidden      bool
	Starter     bool
}

// Creates the category directory if it doesn't exist.
func (c *Category) Init() error {
	exist := dir.Exists(c.Path)
	if !exist {
		return dir.Create(c.Path)
	}
	return nil
}

// Generates a formatted string representation of the category and its notes.
func (c *Category) Print() string {
	if len(*c.Notes) == 0 {
		return ""
	}
	out := fmt.Sprintf("%s%s%s\n\n", "\033[1;32m", strings.ToUpper(c.Title), "\033[0m") // Add prefix + color
	for _, n := range *c.Notes {
		out += fmt.Sprintf("  %s\n", n.Print())
	}
	return out
}

// A collection of categories, implemented as a map.
type Categories map[string]*Category

// Retrieves a category by its name.
func (cs Categories) Find(name string) (*Category, error) {
	c, ok := cs[name]
	if !ok {
		return &Category{}, fmt.Errorf("no category with name [%s]", name)
	}
	return c, nil
}

// Adds a category to the collection.
func (cs Categories) Set(name string, c *Category) {
	cs[name] = c
}

// Returns a new Categories collection containing only visible categories with notes.
func (cs Categories) GetVisible() *Categories {
	list := &Categories{}
	for _, c := range cs {
		if !c.Hidden && len(*c.Notes) > 0 {
			list.Set(c.Name, c)
		}
	}
	return list
}

// Retrieves the starter category from the collection.
func (cs Categories) GetStarter() (*Category, error) {
	for _, c := range cs {
		if c.Starter {
			return c, nil
		}
	}
	return &Category{}, fmt.Errorf("no starter category found")
}

// Generates a formatted string representation of all visible categories and their notes, sorted by weight.
func (cs Categories) Print() string {
	list := cs.GetVisible()

	// Sort
	weights := map[int]string{}
	for _, c := range *list {
		weights[c.Weight] = c.Name
	}
	var keys []int
	for k := range weights {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Print
	out := ""
	i := 0
	for _, k := range keys {
		c, _ := list.Find(weights[k])
		i++
		out += c.Print()
		if i != len(*list) {
			out += "\n"
		}
	}
	return out
}
