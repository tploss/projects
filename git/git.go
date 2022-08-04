package git

import (
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	gssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"gitlab.com/tploss/projects/globals"
)

func ClonePullAll(project globals.Project) error {
	fmt.Println(git.Open)
	return nil
}

type GitCmd struct {
	PrivateKey string `name:"private-key" short:"k" type:"existingfile" env:"SSH_KEY" default:"~/.ssh/id_rsa" help:"Path to SSH private key that should be used"`
}

type PullCmd struct {
	Project string `short:"p" help:"Only pull/clone repos of given project"`
}

func (p PullCmd) Run(g *globals.G) error {
	auth, err := gssh.NewSSHAgentAuth("")
	fmt.Println(auth.User)
	if err != nil {
		return err
	}
	pulled := make([]globals.Repo, 0)
	for _, proj := range g.Conf.Projects {
		if p.Project != "" && proj.Name != p.Project {
			continue
		}
		for _, repo := range proj.Repos {
			rDir, err := g.Conf.Dir(repo.Name)
			if err != nil {
				return err
			}
			if info, err := os.Stat(filepath.Join(rDir, ".git")); err != nil {
				// clone needed
				globals.EnsureDir(rDir)
				_, err := git.PlainClone(rDir, false, &git.CloneOptions{URL: repo.Url, Auth: auth})
				if err != nil {
					return err
				}
				pulled = append(pulled, repo)
			} else if !info.IsDir() {
				return fmt.Errorf("repo directory %s contains .git but it is not a directory", rDir)
			} else {
				// pull needed
				re, err := git.PlainOpen(rDir)
				if err != nil {
					return err
				}
				work, err := re.Worktree()
				if err != nil {
					return err
				}
				if err := work.Pull(&git.PullOptions{Auth: auth}); err != nil {
					return err
				}
				pulled = append(pulled, repo)
			}
		}
	}
	if len(pulled) == 0 {
		return fmt.Errorf("no repos were found for %s", p.Project)
	}
	return nil
}
