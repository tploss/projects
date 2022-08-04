package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gitlab.com/tploss/projects/globals"
	"gopkg.in/yaml.v3"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getBaseDir() (string, error) {
	return os.UserHomeDir()
}

func ensureDir(dir string) error {
	if s, err := os.Stat(dir); err != nil {
		return os.Mkdir(dir, os.ModePerm)
	} else if !s.IsDir() {
		return fmt.Errorf("%s exists but is not a directory", dir)
	}
	return nil
}

func stashed() {
	baseDir, err := getBaseDir()
	panicOnErr(err)

	f, err := os.Open("example/projects.yaml")
	panicOnErr(err)
	var ps globals.Conf
	panicOnErr(yaml.NewDecoder(f).Decode(&ps))

	for _, p := range ps.Projects {
		pDir := filepath.Join(baseDir, p.Name)
		panicOnErr(ensureDir(pDir))
		for _, r := range p.Repos {
			rDir := filepath.Join(pDir, r.Name)
			panicOnErr(ensureDir(rDir))

			if fileList, err := os.ReadDir(rDir); err != nil {
				panic(err)
			} else if len(fileList) > 0 {
				fmt.Printf("Since %s is not empty, git clone is skipped\n", rDir)
				continue
			}

			cmd := exec.Command("git", "clone", r.Url, rDir)
			globals.AttachOwnWriters(cmd)
			fmt.Println(cmd)
			panicOnErr(cmd.Run())
		}
	}

}
