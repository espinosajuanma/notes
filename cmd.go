package notes

import (
	"fmt"
	"strconv"
	"strings"

	n "github.com/espinosajuanma/notes/pkg"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/vars"
)

var app *n.App

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()

	// Init Configs
	config := &n.Config{
		Template: Z.Vars.Get("template"),
	}
	if config.Template == "" {
		names, _ := config.GetTemplates()
		term.Print("Which template would you like to use:")
		for i, name := range names {
			term.Printf("  %d. %s", i+1, name)
		}
		name := term.Prompt("Choose: ")
		n, _ := strconv.Atoi(name)
		config.Template = names[n-1]
	}
	if config.Template == "" {
		panic("you must set a template or user config file")
	}

	path := "./tasks" // fix
	app = n.NewApp(Z.ExeName, path, config)

	err := app.Init()
	if err != nil {
		panic(err)
	}
	err = Z.Vars.Set("template", config.Template)
	if err != nil {
		panic(err)
	}

	// Init git
	remote := Z.Vars.Get(".remote")
	if remote != "" {
		app.Repository.Remote = remote
	}

	// Inite Notes
	latest := Z.Vars.Get(".latest")
	if latest != "" {
		n, err := app.Get(latest)
		if err == nil {
			app.SetLatest(n)
		}
	}

	app.SetColor(Z.Vars.Get(".color"))
	app.SetPrefix(Z.Vars.Get(".prefix"))
}

var Cmd = &Z.Cmd{
	Name: `notes`,
	Commands: []*Z.Cmd{
		MoveCmd,
		ListCmd, AddCmd, EditCmd, RemoveCmd,
		LogCmd, StatusCmd,
		help.Cmd, vars.Cmd, conf.Cmd,
	},
	Shortcuts: Z.ArgMap{
		"color":  {"var", "set", "color"},
		"prefix": {"var", "set", "prefix"},
	},
	Version:     `v0.0.1`,
	Source:      `https://github.com/espinosajuanma/notes`,
	Issues:      `https://github.com/espinosajuanma/notes/issues`,
	Summary:     `Notes app`,
	Description: ``,
}

// CRUD

var ListCmd = &Z.Cmd{
	Name:        "list",
	Commands:    []*Z.Cmd{help.Cmd},
	Aliases:     []string{"ls"},
	Summary:     "",
	Description: "",
	Call: func(caller *Z.Cmd, args ...string) error {
		if len(args) >= 1 {
			c, err := app.Categories.Find(args[0])
			if err != nil {
				return err
			}
			s := c.Print()
			if s != "" {
				term.Print(s)
			}
			return nil
		}

		s := app.Categories.Print()
		if s != "" {
			term.Print(s)
			return nil
		}

		root := caller.Root()
		cmds := root.Commands
		for _, cmd := range cmds {
			if cmd.Name == "help" {
				cmd.Call(root, "help")
				return nil
			}
		}
		return nil
	},
}

var AddCmd = &Z.Cmd{
	Name:        "add",
	Aliases:     []string{"create", "new"},
	Commands:    []*Z.Cmd{},
	Summary:     "",
	Description: "",
	Usage:       "[title]",
	Call: func(caller *Z.Cmd, args ...string) error {
		var title string
		if len(args) >= 1 {
			title = args[0]
		}
		if title == "" {
			title = term.Prompt("Title: ")
		}
		if title == "" {
			return fmt.Errorf("title can'n be empty")
		}
		note, err := app.Create(title)
		if err != nil {
			return err
		}

		app.Notes.Set(note.ID, note)
		app.SetLatest(note)

		y := term.Prompt("Would you like to edit the note? (y/N) ")
		y = strings.ToLower(y)
		term.Print(y)
		if y == "y" || y == "yes" {
			err := app.Update(note.ID)
			if err != nil {
				return err
			}
		}

		y = term.Prompt("Would you like to move the note? (y/N) ")
		y = strings.ToLower(y)
		if y == "y" || y == "yes" {
			if note.MoveNext() != nil {
				return err
			}
		}

		return Z.Vars.Set(".latest", note.ID)
	},
}

var EditCmd = &Z.Cmd{
	Name:        "edit",
	Aliases:     []string{"ed", "e"},
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     "",
	Description: "",
	Usage:       "[id]",
	Call: func(caller *Z.Cmd, args ...string) error {
		id := "latest"
		if len(args) > 0 {
			id = args[0]
		}
		err := app.Update(id)
		if err != nil {
			return err
		}
		note, err := app.Notes.GetLatest()
		if err != nil {
			return err
		}
		return Z.Vars.Set(".latest", note.ID)
	},
}

var RemoveCmd = &Z.Cmd{
	Name:        "remove",
	Aliases:     []string{"rm"},
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     "",
	Description: "",
	Usage:       "[id]",
	Call: func(caller *Z.Cmd, args ...string) error {
		id := "latest"
		if len(args) > 0 {
			id = args[0]
		}
		return app.Remove(id)
	},
}

// Transition

var MoveCmd = &Z.Cmd{
	Name:        "move",
	Aliases:     []string{"mv"},
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     "",
	Description: "",
	Usage:       "[transition] [id]",
	Call: func(caller *Z.Cmd, args ...string) error {
		id := "latest"
		var transition string
		if len(args) == 0 {
			return ListCmd.Call(caller)
		}
		if len(args) > 1 {
			transition = args[0]
			id = args[1]
		} else {
			transition = args[0]
		}
		category, err := app.Move(id, transition)
		if err != nil {
			return err
		}
		s := category.Print()
		term.Print(s)
		return nil
	},
}

// Git

var LogCmd = &Z.Cmd{
	Name:        "log",
	Aliases:     []string{},
	Commands:    []*Z.Cmd{},
	Summary:     "",
	Description: "",
	Usage:       "",
	Call: func(caller *Z.Cmd, args ...string) error {
		return app.Repository.Log()
	},
}

var StatusCmd = &Z.Cmd{
	Name:        "status",
	Aliases:     []string{},
	Commands:    []*Z.Cmd{},
	Summary:     "",
	Description: "",
	Usage:       "",
	Call: func(caller *Z.Cmd, args ...string) error {
		return app.Repository.Status()
	},
}

/*
var RemoteCmd = &Z.Cmd{
	Name:        "remote",
	Commands:    []*Z.Cmd{},
	Summary:     "",
	Description: "",
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) > 0 {
			err := x.Root().Set("remote", args[0])
			if err != nil {
				return err
			}
		}

		remote, _ := x.Root().Get("remote")
		term.Print(remote)
		return nil
	},
}
*/
