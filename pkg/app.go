package notes

import (
	"fmt"
	"os"
	"path"

	"github.com/rwxrob/uniq"
)

const CONTENT_DIR = "content"
const DEFAULT_COLOR = "\033[1;32m"
const DEFAULT_PREFIX = ""

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

// INITIALIZATION

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

// CRUD

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

func (app *App) Get(id string) (*Note, error) {
	if id == "" || id == "latest" {
		return app.Notes.GetLatest()
	}
	return app.Notes.Get(id)
}

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

// TODO
func (app *App) Pin(note *Note) error {
	if note.Pinned {
		return fmt.Errorf("note [%s] is already pinned", note.ID)
	}
	// Repo move

	// Commit

	note.Pinned = true
	return nil
}

// TODO
func (app *App) Unpin(note *Note) error {
	if !note.Pinned {
		return fmt.Errorf("note [%s] is not pinned", note.ID)
	}
	// Repo move

	// Commit

	note.Pinned = false
	return nil
}

// TODO
func (app *App) Query(q string) (*Notes, error) {
	// Find on content
	return &Notes{}, nil
}

// TODO
func (app *App) FindByTag(tags []string) (*Notes, error) {
	return &Notes{}, nil
}

func (app *App) GetLatest() (*Note, error) {
	return app.Notes.GetLatest()
}

func (app *App) SetLatest(latest *Note) {
	app.Notes.SetLatest(latest)
}

func (app *App) SetColor(c string) error {
	if c == "" {
		return fmt.Errorf("color can't be empty")
	}
	app.Color = c
	return nil
}

func (app *App) SetPrefix(p string) error {
	if p == "" {
		return fmt.Errorf("prefix can't be empty")
	}
	app.Prefix = p
	return nil
}

func (app *App) GetTransitions() []string {
	var transitions []string
	for _, c := range *app.Categories {
		for _, t := range *c.Transitions {
			transitions = append(transitions, t.Name)
		}
	}
	return transitions
}
