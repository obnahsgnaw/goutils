package emailutil

import "gopkg.in/gomail.v2"

type GoMail struct {
	dialer   *gomail.Dialer
	host     string
	port     int
	username string
	password string
}

func newGoMail(host string, port int, username string, password string) *GoMail {
	return &GoMail{
		dialer:   gomail.NewDialer(host, port, username, password),
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (ge *GoMail) DialAndSend(m ...*Email) error {
	var emails []*gomail.Message
	for _, mm := range m {
		emails = append(emails, toGoMail(mm))
	}
	return ge.dialer.DialAndSend(emails...)
}

func (ge *GoMail) Ssl(ssl bool) {
	ge.dialer.SSL = ssl
}

func toGoMail(e *Email) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", e.to...)
	m.SetHeader("Subject", e.subject)
	m.SetBody(e.contentType, e.content)
	if e.cc != "" {
		m.SetHeader("Cc", e.cc)
	}
	if e.bcc != "" {
		m.SetHeader("Bcc", e.bcc)
	}
	for _, attach := range e.attaches {
		m.Attach(attach)
	}
	return m
}
