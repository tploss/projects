package globals

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v3"
)

type G struct {
	Config  string           `short:"c" type:"path" help:"Path to config/projects file" default:"~/.config/projects.yaml"`
	Version kong.VersionFlag `short:"v"`
	Conf    Conf             `kong:"-"`
}

func (g *G) AfterApply() error {
	dir, err := ExpandTilde(g.Config)
	if err != nil {
		return err
	}
	g.Config = dir

	if _, err := os.Stat(g.Config); err != nil {
		return fmt.Errorf("Config/Projects file could not be opened: %w", err)
	}

	f, err := os.Open(g.Config)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = yaml.NewDecoder(f).Decode(&g.Conf); err != nil {
		return err
	}
	if err = g.Conf.Validate(); err != nil {
		return err
	}

	return nil
}

func (g *G) WriteConf() error {
	f, err := os.Create(g.Config)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(f)
	defer enc.Close()
	return enc.Encode(g.Conf)
}
