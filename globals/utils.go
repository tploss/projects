package globals

import (
	"os"
	"os/exec"
)

func AttachOwnWriters(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}
