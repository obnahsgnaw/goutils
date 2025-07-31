package example

import (
	"github.com/obnahsgnaw/goutils/emailutil"
	"log"
)

type VerifyCode struct {
	emailutil.BaseEmailBuilder
}

var _verifyCodeEmail = newVerifyCodeEmail()

func newVerifyCodeEmail() *VerifyCode {
	e := &VerifyCode{}
	e.Initialize(e)
	return e
}

func VerifyCodeEmail() *VerifyCode {
	return _verifyCodeEmail
}

// Send Format the input anc call parent sends
func (e *VerifyCode) Send(code string, to string) error {
	return e.BaseEmailBuilder.Send(code, to)
}

func (e *VerifyCode) Subject() string {
	return "welcome"
}

func (e *VerifyCode) Template() string {
	return `Hello User,

Thank you for signing up! Your verification code is: {{.}}

Best regards,
Your Company`
}

func (e *VerifyCode) Failed(err error) {
	log.Println("handle verify code failed", err)
}
