package cmdutil

import "testing"

func TestLogCmd(t *testing.T) {
	err := RunLogCmd(func(msg string) {
		println("Log:" + msg)
	}, "ping", "www.baidu.com")
	if err != nil {
		t.Error(err)
		return
	}
}
