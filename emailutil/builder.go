package emailutil

import (
	"bytes"
	"errors"
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/obnahsgnaw/goutils/structutil"
	"text/template"
)

type EmailBuilder interface {
	From() string
	Subject() string
	ContentType() string
	Template() string
	Cc() string
	Bcc() string
	Attachments() []string
}

type BaseEmailBuilder struct {
	structutil.NamedStruct
	errutil.ErrBuilder
	impl    EmailBuilder
	manager *Manager
}

func (s *BaseEmailBuilder) Initialize(impl EmailBuilder) {
	s.impl = impl
	s.ParseName(s.impl)
	s.ErrPrefix = s.GetName()
}

func (s *BaseEmailBuilder) From() string {
	return "" // use global
}

func (s *BaseEmailBuilder) Cc() string {
	return ""
}

func (s *BaseEmailBuilder) Bcc() string {
	return ""
}

func (s *BaseEmailBuilder) Attachments() []string {
	return nil
}

func (s *BaseEmailBuilder) ContentType() string {
	return "text/html"
}

func (s *BaseEmailBuilder) RegisterTo(m *Manager) {
	s.manager = m
}

func (s *BaseEmailBuilder) Send(data any, to ...string) error {
	if s.impl == nil {
		panic("BaseEmailBuilder: Initialize() first when email Send()")
	}
	if s.manager == nil {
		panic("BaseEmailBuilder: RegisterTo() manager first when email Send()")
	}
	m := NewEmail(s.impl.Subject(), "", to...)
	m.From(s.impl.From())
	m.Cc(s.impl.Cc())
	m.Bcc(s.impl.Bcc())
	t := s.impl.Template()
	if t == "" {
		return errors.New("email content template empty")
	}
	tmpl, err := template.New("email").Parse(t)
	if err != nil {
		return errors.New("email content template invalid")
	}
	var body bytes.Buffer
	if tErr := tmpl.Execute(&body, data); tErr != nil {
		return errors.New("email content template invalid: " + tErr.Error())
	}
	m.ContentType(s.impl.ContentType())
	m.Content(body.String())
	m.Attaches(s.impl.Attachments()...)

	return s.manager.Send(m)
}
