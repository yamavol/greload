package internal

import (
	"os"
	"os/exec"
	"runtime"
)

// ExecuteCommand executes a shell command in an OS-independent manner.
// On Windows, it uses cmd /C, on Unix-like systems it uses sh -c.
func ExecuteCommand(command string) error {
	var cmd *exec.Cmd

	// Determine shell based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	// Inherit stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
