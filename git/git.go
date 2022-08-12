package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/tploss/projects/globals"
)

type GitCmd struct {
	Project string   `short:"p" optional:"" help:"Only pull/clone repos of given project"`
	GitArgs []string `arg:"" passthrough:""`
}

func (gc GitCmd) Run(g *globals.G) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	var projects []globals.Project
	if gc.Project != "" {
		proj, exists := g.Conf.GetProject(gc.Project)
		if !exists {
			return fmt.Errorf("project %s is not found in the config", gc.Project)
		}
		projects = []globals.Project{proj}
	} else {
		projects = g.Conf.Projects
	}

	for _, p := range projects {
		for _, r := range p.Repos {
			fmt.Printf("Handling %s\n", filepath.Join(p.Name, r.Name))
			path := g.Conf.Dir(p, r)

			if err = ensureGit(path, r); err != nil {
				return err
			}
			err = os.Chdir(path)
			if err != nil {
				return err
			}
			cmd := exec.Command("git", gc.GitArgs...)
			globals.AttachOwnWriters(cmd)
			if err = cmd.Run(); err != nil {
				return err
			}
			fmt.Println("========")
		}
	}
	return os.Chdir(wd)
}

func ensureGit(path string, repo globals.Repo) error {
	globals.EnsureDir(path)
	gPath := filepath.Join(path, ".git")

	if stat, err := os.Stat(gPath); err != nil {
		if files, _ := ioutil.ReadDir(path); len(files) != 0 {
			return fmt.Errorf("path %s does not contain a git repository but has content", path)
		}

		// have to clone since git repo does not exist
		fmt.Println("Cloning repository since it does not exist yet")
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		if err = os.Chdir(path); err != nil {
			return err
		}
		cmd := exec.Command("git", "clone", repo.Url, ".")
		globals.AttachOwnWriters(cmd)
		if err = cmd.Run(); err != nil {
			return err
		}
		return os.Chdir(wd)
	} else {
		if !stat.IsDir() {
			return fmt.Errorf("path %s exists but is not a directory", gPath)
		}
	}
	return nil
}
