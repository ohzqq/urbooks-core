package urbooks

import (
	"bytes"
	"fmt"
	"os/exec"
)

func shellout(args ...string) (error, string, string) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)
	cmdArgs := []string{"-c"}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("bash", cmdArgs...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	fmt.Println(cmd.String())
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}
