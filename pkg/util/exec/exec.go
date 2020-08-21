package exec

import (
	"os"
	"os/exec"
)

// RunCommand is a simplified Command runner that includes command output only
// when the command produces an error.
func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	stout := NewPrefixWriter(os.Stdout, "\t")
	defer stout.Close()

	cmd.Stdout = stout

	sterr := NewPrefixWriter(os.Stderr, "\t")
	defer sterr.Close()

	cmd.Stderr = sterr

	return cmd.Run()
}

// RunCommand is a simplified Command runner that includes command output only
// when the command produces an error.
func RunCommandWorkingDir(command, wd string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = wd

	stout := NewPrefixWriter(os.Stdout, "\t")
	defer stout.Close()

	cmd.Stdout = stout

	sterr := NewPrefixWriter(os.Stderr, "\t")
	defer sterr.Close()

	cmd.Stderr = sterr

	return cmd.Run()
}
