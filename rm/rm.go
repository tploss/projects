package rm

import (
	"fmt"

	"gitlab.com/tploss/projects/globals"
)

type RmCmd struct {
	Projects []string `arg:"" required:"" help:"Specify the project(s) to remove from the file system"`
	All      bool     `short:"a" help:"Also remove from project configuration"`
}

func (r RmCmd) Run(g *globals.G) error {
	for _, name := range r.Projects {
		proj, exists := g.Conf.GetProject(name)
		if !exists {
			return fmt.Errorf("project %s does not exist in configuration %s", proj, g.Config)
		}

		if err := proj.RemoveFromDisk(g.Conf.BaseDir); err != nil {
			return err
		}

		if r.All {
			g.Conf.RemoveProject(name)
		}
	}
	return nil
}
