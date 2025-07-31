package emailutil

import "errors"

type Manager struct {
	dialer     Dialer
	ssl        bool
	globalFrom map[string]string
	emails     []*Email
}

func newManager(d Dialer, o ...Option) *Manager {
	s := &Manager{
		dialer:     d,
		globalFrom: make(map[string]string),
	}
	for _, opt := range o {
		opt(s)
	}
	s.dialer.Ssl(s.ssl)
	return s
}

func (s *Manager) Build(subject string, content string, to ...string) *Email {
	e := NewEmail(subject, content, to...)
	s.emails = append(s.emails, e)
	return e
}

func (s *Manager) Send(m ...*Email) error {
	m = append(m, s.emails...)
	s.emails = nil
	for _, mm := range m {
		if mm.from == "" {
			mm.from = "global"
		}
		if v, ok := s.globalFrom[mm.from]; ok {
			mm.from = v
		}
		if mm.from == "global" {
			mm.from = "None"
		}
		if len(mm.to) == 0 {
			return errors.New("you must specify the mail[" + mm.subject + "] sent to")
		}
	}
	return s.dialer.DialAndSend(m...)
}

type Option func(*Manager)

func Ssl(ssl bool) Option {
	return func(s *Manager) {
		s.ssl = ssl
	}
}

func From(from string) Option {
	return func(s *Manager) {
		s.globalFrom["global"] = from
	}
}

func NamedFrom(key, from string) Option {
	return func(s *Manager) {
		s.globalFrom[key] = from
	}
}

func NewManager(host string, port int, username string, password string, o ...Option) *Manager {
	return newManager(newGoMail(host, port, username, password), o...)
}

func NewDevManager() *Manager {
	return newManager(newLogMail())
}
