package globals

import (
	"fmt"
	"regexp"
)

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
