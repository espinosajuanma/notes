package notes

import "fmt"

type Notes map[string]*Note

func (ns Notes) Set(ID string, n *Note) {
	ns[ID] = n
}

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

func (ns Notes) Unset(note *Note) {
	delete(ns, note.ID)
}

func (ns Notes) GetLatest() (*Note, error) {
	for _, n := range ns {
		if n.Latest {
			return n, nil
		}
	}
	return &Note{}, fmt.Errorf("latest task not found")
}

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

func (ns Notes) Print() string {
	p := ""
	for _, n := range ns {
		t, _ := n.GetTitle()
		p += fmt.Sprintf("  [%s] %s\n", n.ID, t)
	}
	return p
}
