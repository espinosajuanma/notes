package notes

import "fmt"

// Represents a collection of notes.
type Notes map[string]*Note

// Adds a note to the collection.
func (ns Notes) Set(ID string, n *Note) {
	ns[ID] = n
}

// Retrieves a note from the collection by its ID.  Returns an error if not found.
func (ns Notes) Get(id string) (*Note, error) {
	if id == "" || id == "latest" {
		return ns.GetLatest()
	}

	n, ok := ns[id]
	if !ok {
		return &Note{}, fmt.Errorf("task [%s] not found", id)
	}
	return n, nil
}

// Removes a note from the collection.
func (ns Notes) Unset(note *Note) {
	delete(ns, note.ID)
}

// Retrieves the latest note from the collection.
func (ns Notes) GetLatest() (*Note, error) {
	for _, n := range ns {
		if n.Latest {
			return n, nil
		}
	}
	return &Note{}, fmt.Errorf("latest task not found")
}

// Sets the specified note as the latest in the collection.
func (ns Notes) SetLatest(latest *Note) {
	for _, n := range ns {
		if n.Latest && n != latest {
			n.Latest = false
		}
		if n == latest {
			n.Latest = true
		}
	}
}

// Retrieves a sub-collection containing only pinned notes.
func (ns Notes) GetPinned() (*Notes, error) {
	list := Notes{}
	for _, n := range ns {
		if n.Pinned {
			list[n.ID] = n
		}
	}
	if len(list) > 0 {
		return &list, nil
	}
	return &list, fmt.Errorf("no pinned notes")
}

// Retrieves a sub-collection containing only notes belonging to a specific category.
func (ns Notes) GetByCategory(c *Category) (*Notes, error) {
	list := Notes{}
	for _, n := range ns {
		if n.Category == c {
			list[n.ID] = n
		}
	}
	if len(list) > 0 {
		return &list, nil
	}
	return &list, fmt.Errorf("category doesn't has notes")
}

// Generates a formatted string representation of all notes in the collection.
func (ns Notes) Print() string {
	p := ""
	for _, n := range ns {
		t, _ := n.GetTitle()
		p += fmt.Sprintf("  [%s] %s\n", n.ID, t)
	}
	return p
}
