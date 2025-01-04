package cmdutil

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func RunLogCmd(logFn func(msg string), cmdName string, args ...string) (err error) {
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

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			lines := bytes.Split(bytes.Trim(scanner.Bytes(), "\x00"), []byte{'\n'})
			for _, line := range lines {
				lineStr := string(line)
				if lineStr != "" && logFn != nil {
					logFn(lineStr)
				}
			}
		}
		if err1 := scanner.Err(); err1 != nil {
			log.Fatal(err1)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			lines := bytes.Split(bytes.Trim(scanner.Bytes(), "\x00"), []byte{'\n'})
			for _, line := range lines {
				lineStr := string(line)
				if lineStr != "" && logFn != nil {
					logFn(lineStr)
				}
			}
		}
		if err1 := scanner.Err(); err1 != nil {
			log.Fatal(err1)
		}
	}()

	return cmd.Wait()
}

func RunCmd(cmdName string, args ...string) (err error) {
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return
	}

	return cmd.Wait()
}
