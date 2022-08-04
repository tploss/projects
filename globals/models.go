package globals

import (
	"fmt"
	"os"
	"path/filepath"
)

type Conf struct {
	BaseDir  string    `yaml:"base_dir"`
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Name  string `yaml:"name"`
	Repos []Repo `yaml:"repos,omitempty"`
}

type Repo struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

// Can also be used to rename a project by giving the old name for the
// name parameter but using the new name within the proj struct
func (c *Conf) SetProject(name string, newP Project) {
	for i, p := range c.Projects {
		if p.Name == name {
			c.Projects[i] = newP
			return
		}
	}
	c.Projects = append(c.Projects, newP)
}

func (c *Conf) RemoveProject(name string) {
	for i, p := range c.Projects {
		if p.Name == name {
			c.Projects = append(c.Projects[:i], c.Projects[i+1:]...)
		}
	}
}

func (c Conf) GetProject(name string) (Project, bool) {
	for _, p := range c.Projects {
		if p.Name == name {
			return p, true
		}
	}
	return Project{Name: name, Repos: []Repo{{Name: "Placeholder", Url: "Placeholder"}}}, false
}

func (p Project) RemoveFromDisk(baseDir string) error {
	pPath := filepath.Join(baseDir, p.Name)
	return os.RemoveAll(pPath)
}

func (c Conf) Dir(repo string) (string, error) {
	for _, p := range c.Projects {
		if path := p.dir(repo); path != "" {
			return filepath.Join(c.BaseDir, path), nil
		}
	}
	return "", fmt.Errorf("repo %s was not found in config", repo)
}

func (p Project) dir(repo string) string {
	for _, r := range p.Repos {
		if r.Name == repo {
			return filepath.Join(p.Name, r.Name)
		}
	}
	return ""
}
