package exec

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func ExecAndPrint(cmd *exec.Cmd) (string, error) {
	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer
	cmd.Run()
	cmd.Wait()
	errBufferStr := errBuffer.String()
	if len(errBufferStr) > 0 {
		return "", errors.New(errBufferStr)
	}
	return strings.TrimSpace(outBuffer.String()), nil
}
