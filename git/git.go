package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	"gitlab.com/tploss/projects/globals"
)

func ClonePullAll(project globals.Project) error {
	fmt.Println(git.Open)
	return nil
}
