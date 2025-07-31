package example

import (
	"github.com/obnahsgnaw/goutils/emailutil"
	"testing"
)

func TestVerifyCodeEmail(t *testing.T) {
	s := emailutil.NewDevManager()

	e := VerifyCodeEmail()
	e.RegisterTo(s)
	err := e.Send("123456", "xxx@xxxx.com")
	if err != nil {
		t.Error(err)
	}
}
