package globals

import (
	"fmt"
	"regexp"
)

type Conf struct {
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

func (c Conf) GetProject(name string) (Project, bool) {
	for _, p := range c.Projects {
		if p.Name == name {
			return p, true
		}
	}
	return Project{Name: name, Repos: []Repo{{Name: "Placeholder", Url: "Placeholder"}}}, false
}

func (c Conf) Validate() error {
	for _, p := range c.Projects {
		if err := p.validate(); err != nil {
			return fmt.Errorf("config has an issue: %w", err)
		}
	}
	return nil
}

func (p Project) validate() error {
	if err := validate(p.Name); err != nil {
		return fmt.Errorf("Project name %s has an issue: %w", p.Name, err)
	}
	for _, r := range p.Repos {
		if err := r.validate(); err != nil {
			return fmt.Errorf("Project repos have an issue: %w", err)
		}
	}
	return nil
}

func (r Repo) validate() error {
	if err := validate(r.Name); err != nil {
		return fmt.Errorf("Repo name %s has an issue: %w", r.Name, err)
	}
	return nil
}

func validate(s string) error {
	pattern := `^[a-zA-Z0-9_-]+$`
	if !regexp.MustCompile(pattern).Match([]byte(s)) {
		return fmt.Errorf("%s does not match %s", s, pattern)
	}
	return nil
}
