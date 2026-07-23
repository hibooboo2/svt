package exec

import (
	"log/slog"
	"os/exec"
)

func Command(cmd string, args ...string) *exec.Cmd {
	slog.Info("Running command: ", "cmd", cmd, "args", args)
	return exec.Command(cmd, args...)
}
