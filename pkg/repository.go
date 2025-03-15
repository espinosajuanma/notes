package notes

import (
	"fmt"
	"path/filepath"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fs/dir"
)

type Repository struct {
	Path   string
	Remote string
}

func (r *Repository) Init() error {
	if r.Path == "" {
		return fmt.Errorf("repository path can't be empty")
	}
	// Create directory if doesn't exist
	exists := dir.Exists(r.Path)
	if !exists {
		err := dir.Create(r.Path)
		if err != nil {
			return err
		}
	}

	// Init git repository
	return r.Command("init", "-q")
}

func (r *Repository) Command(args ...string) error {
	a := []string{"git", "-C", r.Path}
	a = append(a, args...)
	return Z.Exec(a...)
}

func (r *Repository) SetRemote(url string) error {
	if url == "" {
		return fmt.Errorf("path can't be empty")
	}
	r.Remote = url
	return nil
}

func (r *Repository) Add(path string) error {
	if path == "" {
		return fmt.Errorf("path can't be empty")
	}
	path = r.CleanPath(path)
	return r.Command("add", path)
}

func (r *Repository) Move(source, target string) error {
	if source == "" || target == "" {
		return fmt.Errorf("source [%s] or target [%s] can't be empty", source, target)
	}
	source = r.CleanPath(source)
	target = r.CleanPath(target)
	return r.Command("mv", source, target)
}

func (r *Repository) Remove(path string, recursive bool) error {
	if path == "" {
		return fmt.Errorf("path can't be empty")
	}
	path = r.CleanPath(path)
	if recursive {
		return r.Command("rm", "-rf", path)
	}
	return r.Command("rm", path)
}

func (r *Repository) Commit(message string) error {
	if message == "" {
		return fmt.Errorf("commit message can't be empty")
	}
	return r.Command("commit", "-q", "-m", message)
}

func (r *Repository) Push() error {
	if r.Remote == "" {
		return nil
	}
	return r.Command("push")
}

func (r *Repository) Log() error {
	return r.Command("log")
}

func (r *Repository) Status() error {
	return r.Command("status")
}

func (r *Repository) CleanPath(p string) string {
	path, err := filepath.Rel(r.Path, p)
	if err != nil {
		return ""
	}
	return path
}
