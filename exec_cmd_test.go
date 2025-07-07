//go:build linux

package wgiface

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_execCmd_NoErr(t *testing.T) {
	cmd := exec.Command("echo", "secret word")
	out, err := execCmd(cmd)
	assert.NoError(t, err)
	assert.Contains(t, out, "secret word")
}

func Test_execCmd_Err(t *testing.T) {
	cmd := exec.Command("ls", "/nonexistent/directory")
	out, err := execCmd(cmd)
	expErr := &exec.ExitError{}
	assert.ErrorContains(t, err, "ls: cannot access '/nonexistent/directory': No such file or directory")
	assert.ErrorAs(t, err, &expErr)
	assert.Equal(t, 2, expErr.ExitCode())
	assert.Equal(t, out, "")
}
