package edit

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"gitlab.com/tploss/projects/globals"
	"gopkg.in/yaml.v3"
)

type EditCmd struct {
	Project string `arg:"" optional:"" help:"Specify a specific project to edit instead of entire config file"`
}

func (e EditCmd) Run(g *globals.G) error {
	editor := getEditor()
	if editor == "" {
		return errors.New("$EDITOR was not set, use e.g. 'export EDITOR=\"$(which vim)\"'")
	}

	tempWrite, err := os.CreateTemp("", "projects-*.yaml")
	if err != nil {
		return err
	}
	defer os.Remove(tempWrite.Name())
	defer tempWrite.Close()

	openEditor := func() error {
		cmd := exec.Command(editor, tempWrite.Name())
		globals.AttachOwnWriters(cmd)
		return cmd.Run()
	}

	enc := yaml.NewEncoder(tempWrite)
	// define function for encoding so that closing is not forgotten
	encode := func(d interface{}) error {
		err := enc.Encode(d)
		enc.Close()
		return err
	}

	// Need a separate io.File since if we use the same as we did for writing we will hit an EOF error
	tempRead, err := os.Open(tempWrite.Name())
	if err != nil {
		return err
	}
	defer tempRead.Close()
	dec := yaml.NewDecoder(tempRead)

	if e.Project == "" {
		if err = encode(g.Conf); err != nil {
			return err
		}
		if err = openEditor(); err != nil {
			return err
		}

		var conf globals.Conf
		if err = dec.Decode(&conf); err != nil {
			return err
		}
		g.Conf = conf
	} else {
		p, _ := g.Conf.GetProject(e.Project)
		if err = encode(p); err != nil {
			return err
		}
		if err = openEditor(); err != nil {
			return err
		}

		var proj globals.Project
		if err = dec.Decode(&proj); err != nil {
			return err
		}
		g.Conf.SetProject(e.Project, proj)
	}

	if err = g.Conf.Validate(); err != nil {
		return fmt.Errorf("validation of cofiguration: %w", err)
	}
	return g.WriteConf()
}

func getEditor() string {
	if val := os.Getenv("EDITOR"); val != "" {
		return val
	}
	/* IDEA: try some well known editors
	for _, e := range []string{"vim", "nano"} {
		if editor, err := exec.LookPath(e); err == nil {
			return editor
		}
	}
	*/
	return ""
}
