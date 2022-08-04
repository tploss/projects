package list

import (
	"fmt"
	"io"
	"os"

	"gitlab.com/tploss/projects/globals"
	"gopkg.in/yaml.v3"
)

var printLocation io.Writer = os.Stdout

type ListCmd struct {
	Projects bool   `short:"p" xor:"type" help:"Only list projects"`
	Repos    string `short:"r" xor:"type" help:"Only list repos of given project"`

	Format string `short:"f" enum:"plain,yaml" default:"yaml" help:"Format for output, choose from: ${enum}"`
}

func (l ListCmd) Run(g *globals.G) error {
	yF, pF := l.Format == "yaml", l.Format == "plain"
	var data interface{} = nil
	if l.Projects {
		ps := make([]globals.Project, len(g.Conf.Projects))
		for i, p := range g.Conf.Projects {
			ps[i] = globals.Project{Name: p.Name, Repos: nil}
		}
		data = ps
	} else if l.Repos != "" {
		for _, p := range g.Conf.Projects {
			if p.Name == l.Repos {
				data = p.Repos
				break
			}
		}
	} else {
		// either all was selected or we default to all
		data = g.Conf.Projects
	}

	if data != nil {
		if yF {
			enc := yaml.NewEncoder(printLocation)
			enc.Encode(data)
			_ = enc.Close()
		} else if pF {
			fmt.Fprintln(printLocation, data)
		}
	}

	return nil
}
