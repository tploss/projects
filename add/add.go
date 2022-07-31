package add

import (
	"gitlab.com/tploss/projects/edit"
	"gitlab.com/tploss/projects/globals"
)

type AddCmd struct {
	Project string `arg:"" required:"" help:"Name for the new project"`
}

func (a AddCmd) Run(g *globals.G) error {
	return edit.EditCmd{Project: a.Project}.Run(g)
}
