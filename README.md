# Notes

Notes is a terminal app to manage a library of notes. It has the ability of dynamic categories and git control.

Other uses:

- Task Manager
- Zettelkasten
- Blog _*_

> _* Might need a new module to build the HTML from markdown_

## Install

You can just grab the latest binary [release](https://github.com/espinosajuanma/notes/releases).

This command can be installed as a standalone program or composed into a Bonzai command tree.

Standalone

```bash
go install github.com/espinosajuanma/notes/cmd/notes@latest
```

Composed

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/espinosajuanma/notes"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, notes.Cmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C notes notes
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Usage

- `notes` - shows list of note, same as `notes list` or `notes ls`
- `notes help` - shows help like manual pages
- `notes list <category>` - shows just category notes
- `notes add [title]` - create a new notes
- `notes remove [id]` - removes a notes, same as `notes rm` 
- `notes edit [id]` - opens the editor
- `notes [transition] [note]` - moves a note to another category

### Git

All changes get registered to on a git repository.

- `notes log` - shows latest git changes
- `notes remove [url]` - configure the remote git repository

### Note ID

The note ID is an unique identifier which is the current date.

Also you can use the keyword `last` which is the latest manipulated note.

## Multiple 

## Modified times

Use git to get latest modified times.

```bash
git log -1 --pretty="format:%ci" ./file.md
```

Remove `-1` to check all the modified dates.

