package globals

import (
	"os"
	"os/exec"
	"path/filepath"
)

func AttachOwnWriters(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func EnsureDir(path string) {
	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		os.Remove(path)
	}
	os.MkdirAll(path, os.ModePerm)
}

func prependHome(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path), nil
}

func ExpandTilde(path string) (string, error) {
	if []rune(path)[0] == '~' {
		return prependHome(string([]rune(path)[1:]))
	}
	return path, nil
}
