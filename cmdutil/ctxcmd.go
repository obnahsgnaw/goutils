package cmdutil

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"
)

func RunCtxLogCmd(ctx context.Context, logFn func(msg string), cmdName string, args ...string) (err error) {
	ctx1, cancel := context.WithCancel(ctx)
	defer cancel()

	var stdout, stderr io.ReadCloser
	cmd := exec.Command(cmdName, args...)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		return
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	go func(in io.ReadCloser, fn func(msg string)) {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			lines := bytes.Split(bytes.Trim(scanner.Bytes(), "\x00"), []byte{'\n'})
			for _, line := range lines {
				lineStr := string(line)
				if lineStr != "" && fn != nil {
					fn(lineStr)
				}
			}
		}
		if err1 := scanner.Err(); err1 != nil {
			log.Fatal(err1)
		}
	}(stdout, logFn)

	go func(in io.ReadCloser, fn func(msg string)) {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			lines := bytes.Split(bytes.Trim(scanner.Bytes(), "\x00"), []byte{'\n'})
			for _, line := range lines {
				lineStr := string(line)
				if lineStr != "" && fn != nil {
					fn(lineStr)
				}
			}
		}
		if err1 := scanner.Err(); err1 != nil {
			log.Fatal(err1)
		}
	}(stderr, logFn)

	go func(ctx2 context.Context, cmd2 *exec.Cmd) {
		for {
			select {
			case <-ctx2.Done():
				_ = cmd2.Process.Kill()
				return
			}
		}
	}(ctx1, cmd)

	return cmd.Wait()
}
