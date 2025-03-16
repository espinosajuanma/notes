package notes

import (
	"fmt"
	"os"
	"path"

	"github.com/rwxrob/uniq"
)

// The name of the directory where notes are stored.
const CONTENT_DIR = "content"

// The default color for output, using ANSI escape codes for bright green.
const DEFAULT_COLOR = "\033[1;32m"

// The default prefix for output, which is an empty string.
const DEFAULT_PREFIX = ""

// Represents the main application structure, holding all core components and data.
type App struct {
	Name       string
	Config     *Config
	Categories *Categories
	Repository *Repository
	Notes      *Notes
	Path       string
	Prefix     string
	Color      string
}

// Creates a new App instance.
// It initializes the App struct with the provided config
func NewApp(name, path string, config *Config) *App {
	app := &App{
		Name:       name,
		Path:       path,
		Config:     config,
		Repository: &Repository{Path: path},
		Categories: &Categories{},
		Notes:      &Notes{},
		Color:      DEFAULT_COLOR,
		Prefix:     DEFAULT_PREFIX,
	}

	return app
}

// Initializes the application, including configuration, repository, categories, and notes.
// It returns an error if any initialization step fails.
func (app *App) Init() error {
	if app.Path == "" {
		return fmt.Errorf("path can't be empty")
	}

	// Initializes Configurations
	err := app.InitConfig()
	if err != nil {
		return err
	}

	// Initializes Repository
	err = app.Repository.Init()
	if err != nil {
		return err
	}

	// Initializes categories
	for _, c := range *app.Categories {
		c.Init()
	}

	// Initializes notes
	err = app.InitNotes()
	if err != nil {
		return err
	}

	return nil
}

// Initializes the application's configuration, including setting up categories and transitions.
// It returns an error if any configuration step fails.
func (app *App) InitConfig() error {
	err := app.Config.Init()
	if err != nil {
		return err
	}

	t := app.Config
	// Set categories
	for _, c := range t.Categories {
		cat := &Category{
			Name:        c.Name,
			Title:       c.Title,
			Hidden:      c.Hidden,
			Starter:     c.Starter,
			Weight:      c.Weight,
			Notes:       &Notes{},
			Transitions: &Transitions{},
			Path:        path.Join(app.Path, CONTENT_DIR, c.Name),
		}
		app.Categories.Set(c.Name, cat)
	}
	// Set transitions
	for _, c := range t.Categories {
		cat, err := app.Categories.Find(c.Name)
		if err != nil {
			return err
		}
		cat.Transitions.Init(c.Transitions, app.Categories)
	}
	return nil
}

// Initializes the notes by reading the content directory and creating Note instances.
// It returns an error if any note initialization step fails.
func (app *App) InitNotes() error {
	for _, c := range *app.Categories {
		dir, err := os.ReadDir(c.Path)
		if err != nil {
			return err
		}
		for _, n := range dir {
			if !n.IsDir() {
				continue
			}
			note := &Note{
				Category:   c,
				ID:         n.Name(),
				Repository: app.Repository,
			}
			err := note.Init()
			app.Notes.Set(note.ID, note)
			c.Notes.Set(note.ID, note)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Creates a new note with the given title.
// It returns a pointer to the created Note and an error if any.
func (app *App) Create(title string) (*Note, error) {
	c, err := app.Categories.GetStarter()
	if err != nil {
		return &Note{}, err
	}
	note := &Note{
		Title:      title,
		Category:   c,
		ID:         uniq.Isosec(),
		Repository: app.Repository,
	}
	note.Create()
	note.Update(true)
	app.SetLatest(note)
	return note, nil
}

// Retrieves a note by its ID or the latest note.
// If id is empty or "latest", it returns the latest note.
// Otherwise, it returns the note with the specified ID.
func (app *App) Get(id string) (*Note, error) {
	if id == "" || id == "latest" {
		return app.Notes.GetLatest()
	}
	return app.Notes.Get(id)
}

// Moves a note to a different category based on a transition.
// It returns a pointer to the new Category and an error if any.
func (app *App) Move(id, transition string) (*Category, error) {
	note, err := app.Get(id)
	if err != nil {
		return &Category{}, err
	}

	// Checks if is a valid move
	newCategory, err := note.Category.Transitions.Find(transition)
	if err != nil {
		return &Category{}, err
	}

	// Repo Moves
	source := note.GetDir()
	target := newCategory.Path // Should check if there is additional
	err = app.Repository.Move(source, target)
	if err != nil {
		return &Category{}, err
	}
	message := fmt.Sprintf("%s [%s]", transition, note.Title)
	if app.Repository.Commit(message) != nil {
		return &Category{}, err
	}

	note.Category = newCategory
	app.SetLatest(note)
	err = app.InitNotes()
	if err != nil {
		return &Category{}, err
	}
	return newCategory, nil
}

// Open the $EDITOR for the user edit a note.
// It returns an error if any.
func (app *App) Update(id string) error {
	note, err := app.Get(id)
	if err != nil {
		return err
	}
	err = note.Edit()
	if err != nil {
		return err
	}
	err = note.Update(false)
	if err != nil {
		return err
	}
	app.SetLatest(note)
	return nil
}

// Removes a note.
// It returns an error if any.
func (app *App) Remove(id string) error {
	note, err := app.Get(id)
	if err != nil {
		return err
	}
	// Removes file
	err = note.Remove()
	if err != nil {
		return err
	}
	// Removes from map
	note.Category.Notes.Unset(note)
	app.Notes.Unset(note)
	return nil
}

// Pins a note.
// It returns an error if any.
// TODO: Implement repository move and commit.
func (app *App) Pin(note *Note) error {
	return fmt.Errorf("Pin not implemented")
}

// Unpins a note.
// It returns an error if any.
// TODO: Implement repository move and commit.
func (app *App) Unpin(note *Note) error {
	return fmt.Errorf("Unpin not implemented")
}

// Queries notes based on content.
// It returns a pointer to the matching Notes and an error if any.
// TODO: Implement content-based search.
func (app *App) Query(q string) (*Notes, error) {
	return &Notes{}, fmt.Errorf("Query not implemented")
}

// Finds notes by tags.
// It returns a pointer to the matching Notes and an error if any.
// TODO: Implement tag-based search.
func (app *App) FindByTag(tags []string) (*Notes, error) {
	return &Notes{}, fmt.Errorf("FindByTag not implemented")
}

// Gets the latest note.
// It returns a pointer to the latest Note and an error if any.
func (app *App) GetLatest() (*Note, error) {
	return app.Notes.GetLatest()
}

// Sets the latest note.
func (app *App) SetLatest(latest *Note) {
	app.Notes.SetLatest(latest)
}

// Sets the output color.
// It returns an error if the color is empty.
func (app *App) SetColor(c string) error {
	if c == "" {
		return fmt.Errorf("color can't be empty")
	}
	app.Color = c
	return nil
}

// Sets the output prefix.
// It returns an error if the prefix is empty.
func (app *App) SetPrefix(p string) error {
	if p == "" {
		return fmt.Errorf("prefix can't be empty")
	}
	app.Prefix = p
	return nil
}

// Gets a list of all available transitions.
// It returns a slice of transition names.
func (app *App) GetTransitions() []string {
	var transitions []string
	for _, c := range *app.Categories {
		for _, t := range *c.Transitions {
			transitions = append(transitions, t.Name)
		}
	}
	return transitions
}
