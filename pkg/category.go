package notes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rwxrob/fs/dir"
)

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

func (c *Category) Init() error {
	exist := dir.Exists(c.Path)
	if !exist {
		return dir.Create(c.Path)
	}
	return nil
}

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

type Categories map[string]*Category

func (cs Categories) Find(name string) (*Category, error) {
	c, ok := cs[name]
	if !ok {
		return &Category{}, fmt.Errorf("no category with name [%s]", name)
	}
	return c, nil
}

func (cs Categories) Set(name string, c *Category) {
	cs[name] = c
}

func (cs Categories) GetVisible() *Categories {
	list := &Categories{}
	for _, c := range cs {
		if !c.Hidden && len(*c.Notes) > 0 {
			list.Set(c.Name, c)
		}
	}
	return list
}

func (cs Categories) GetStarter() (*Category, error) {
	for _, c := range cs {
		if c.Starter {
			return c, nil
		}
	}
	return &Category{}, fmt.Errorf("no starter category found")
}

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
