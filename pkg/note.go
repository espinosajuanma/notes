package notes

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/fs/file"
)

const FILE_NAME = "README.md"

type Note struct {
	ID         string
	Category   *Category
	Repository *Repository
	App        *App
	Title      string
	Content    string
	Tags       []string
	Latest     bool
	Pinned     bool
}

func (n *Note) Init() error {
	title, err := n.GetTitle()
	if err != nil {
		return err
	}
	n.Title = title

	// Check if pin

	// Read content
	content, err := n.GetContent()
	if err != nil {
		return err
	}
	n.Content = content

	// Parse tags
	n.GetTags()

	return nil
}

// PATH functions

func (n *Note) GetFileName() string {
	if n.Pinned {
		return "*" + FILE_NAME
	}
	return FILE_NAME
}

func (n *Note) GetDir() string {
	return path.Join(n.Category.Path, n.ID)
}

func (n *Note) GetPath() string {
	return path.Join(n.GetDir(), n.GetFileName())
}

func (n *Note) GetContent() (string, error) {
	bs, err := os.ReadFile(n.GetPath())
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (n *Note) GetTags() []string {
	return []string{""}
}

func (n *Note) GetTitle() (string, error) {
	file.Head(n.GetPath(), 0)
	dir := n.GetPath()
	head, err := file.Head(dir, 1)
	if err != nil {
		return "", err
	}
	title := strings.TrimPrefix(head[0], "# ")
	return title, nil
}

func (n *Note) Print() string {
	return fmt.Sprintf("[%s] %s", n.ID, n.Title)
}

// File Actions

func (n *Note) Edit() error {
	return file.Edit(n.GetPath())
}

func (n *Note) Create() (*Note, error) {
	if dir.Exists(n.GetDir()) {
		return n, fmt.Errorf("%s already exists", n.GetDir())
	}
	err := dir.Create(n.GetDir())
	if err != nil {
		return n, err
	}
	err = file.Overwrite(n.GetPath(), "# "+n.Title+"\n")
	if err != nil {
		return n, err
	}
	return n, nil
}

// Git Actions

func (n *Note) Commit(action string) error {
	title, err := n.GetTitle()
	if err != nil {
		return err
	}
	message := action + ": " + title
	err = n.Repository.Commit(message)
	return err
}

func (n *Note) Move(c *Category) error {
	current := n.GetDir()
	n.Category = c
	new := n.GetDir()
	err := n.Repository.Move(current, new)
	if err != nil {
		return err
	}
	err = n.Commit("To " + c.Title)
	if err != nil {
		return err
	}
	return nil
}

func (n *Note) MoveNext() error {
	c, err := n.Category.Transitions.Find("next")
	if err != nil {
		return err
	}
	err = n.Move(c)
	if err != nil {
		return err
	}
	return nil
}

func (n *Note) Update(isNew bool) error {
	err := n.Repository.Add(n.GetPath())
	if err != nil {
		return err
	}
	message := "Updated"
	if isNew {
		message = "Created"
	}
	err = n.Commit(message)
	if err != nil {
		return err
	}
	return nil
}

func (n *Note) Remove() error {
	err := n.Repository.Remove(n.GetDir(), true)
	if err != nil {
		return err
	}
	err = n.Commit("Remove")
	if err != nil {
		return err
	}
	return nil
}
