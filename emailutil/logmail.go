package emailutil

import (
	"strings"
)

type LogMail struct {
}

func newLogMail() *LogMail {
	return &LogMail{}
}

func (ge *LogMail) DialAndSend(m ...*Email) error {
	for _, mm := range m {
		toLogMail(mm)
	}
	return nil
}

func (ge *LogMail) Ssl(_ bool) {

}

func toLogMail(e *Email) {
	println("From  :", e.from)
	println("To    :", strings.Join(e.to, ","))
	if e.cc != "" {
		println("Cc    :", e.cc)
	}
	if e.bcc != "" {
		println("Bcc   :", e.bcc)
	}
	println("Subject:", e.subject)
	println("Content:(" + e.contentType + ")")
	println(e.content)
	for _, attach := range e.attaches {
		println("Attach:", attach)
	}
}
