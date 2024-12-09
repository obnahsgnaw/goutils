package cmdutil

import (
	"bytes"
	"os/exec"
)

func RunLogCmd(log func(msg string), cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err = cmd.Start(); err != nil {
		return err
	}

	for {
		temp := make([]byte, 1024)
		_, err = stdout.Read(temp)
		if err != nil {
			break
		}
		lines := bytes.Split(bytes.Trim(temp, "\x00"), []byte{'\n'})
		for _, line := range lines {
			lineStr := string(line)
			if lineStr != "" && log != nil {
				log(lineStr)
			}
		}
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}
