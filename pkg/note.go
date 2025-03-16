package notes

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/fs/file"
)

// The filename used for note content.
const FILE_NAME = "README.md"

// Represents a single note.
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

// Initializes a note, reading its title and content from disk.
func (n *Note) Init() error {
	title, err := n.GetTitle()
	if err != nil {
		return err
	}
	n.Title = title

	content, err := n.GetContent()
	if err != nil {
		return err
	}
	n.Content = content

	n.GetTags()

	return nil
}

// Returns the full path to the note's file.
func (n *Note) GetPath() string {
	return path.Join(n.GetDir(), n.GetFileName())
}

// Returns the directory where the note is stored.
func (n *Note) GetDir() string {
	return path.Join(n.Category.Path, n.ID)
}

// Returns the filename for the note, including a "*" prefix if pinned.
func (n *Note) GetFileName() string {
	if n.Pinned {
		return "*" + FILE_NAME
	}
	return FILE_NAME
}

// Reads the note's content from its file.
func (n *Note) GetContent() (string, error) {
	bs, err := os.ReadFile(n.GetPath())
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// Extracts tags from the note's content (currently a stub).
func (n *Note) GetTags() []string {
	return []string{""} // Placeholder; needs actual tag extraction logic
}

// Extracts the note's title from the first line of its content file.
func (n *Note) GetTitle() (string, error) {
	head, err := file.Head(n.GetPath(), 1)
	if err != nil {
		return "", err
	}
	title := strings.TrimPrefix(head[0], "# ")
	return title, nil
}

// Generates a string representation of the note.
func (n *Note) Print() string {
	return fmt.Sprintf("[%s] %s", n.ID, n.Title)
}

// Opens the note's file in the user's default editor.
func (n *Note) Edit() error {
	return file.Edit(n.GetPath())
}

// Creates the note's directory and file, writing the initial title.
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

// Commits changes to the note's file to the repository.
func (n *Note) Commit(action string) error {
	title, err := n.GetTitle()
	if err != nil {
		return err
	}
	message := action + ": " + title
	err = n.Repository.Commit(message)
	return err
}

// Moves the note to a new category.
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

// Moves the note to the next category according to its transitions.
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

// Updates the note's content and commits the changes to the repository.
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

// Removes the note's directory and commits the removal to the repository.
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
