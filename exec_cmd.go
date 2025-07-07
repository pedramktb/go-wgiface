package wgiface

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execCmd(cmd *exec.Cmd) (string, error) {
	stderr, stdout := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stderr, cmd.Stdout = stderr, stdout
	if err := cmd.Run(); err != nil {
		return stdout.String(), fmt.Errorf("%w: %s", err, stderr.String())
	}
	return stdout.String(), nil
}
