package notes

import (
	"fmt"
	"path/filepath"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fs/dir"
)

// Represents a Git repository for managing note data.
type Repository struct {
	Path   string
	Remote string
}

// Initializes the Git repository.  Creates the directory if needed and runs `git init`.
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

// Executes a Git command within the repository.
func (r *Repository) Command(args ...string) error {
	a := []string{"git", "-C", r.Path}
	a = append(a, args...)
	return Z.Exec(a...)
}

// Sets the remote URL for the repository.
func (r *Repository) SetRemote(url string) error {
	if url == "" {
		return fmt.Errorf("remote URL can't be empty")
	}
	r.Remote = url
	return nil
}

// Adds a file to the Git staging area.
func (r *Repository) Add(path string) error {
	if path == "" {
		return fmt.Errorf("path can't be empty")
	}
	path = r.CleanPath(path)
	return r.Command("add", path)
}

// Moves a file or directory within the repository.
func (r *Repository) Move(source, target string) error {
	if source == "" || target == "" {
		return fmt.Errorf("source or target paths can't be empty")
	}
	source = r.CleanPath(source)
	target = r.CleanPath(target)
	return r.Command("mv", source, target)
}

// Removes a file or directory from the repository.
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

// Commits changes to the repository.
func (r *Repository) Commit(message string) error {
	if message == "" {
		return fmt.Errorf("commit message can't be empty")
	}
	return r.Command("commit", "-q", "-m", message)
}

// Pushes changes to the remote repository.
func (r *Repository) Push() error {
	if r.Remote == "" {
		return nil // No remote set, so nothing to push
	}
	return r.Command("push")
}

// Shows the commit history of the repository.
func (r *Repository) Log() error {
	return r.Command("log")
}

// Shows the status of the repository.
func (r *Repository) Status() error {
	return r.Command("status")
}

// Cleans a given path, making it relative to the repository root.
func (r *Repository) CleanPath(p string) string {
	path, err := filepath.Rel(r.Path, p)
	if err != nil {
		return ""
	}
	return path
}
