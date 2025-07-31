package emailutil

type Email struct {
	subject     string
	from        string
	to          []string
	cc          string
	bcc         string
	attaches    []string
	contentType string
	content     string
}

func NewEmail(subject string, content string, to ...string) *Email {
	return &Email{
		subject: subject,
		from:    "global",
		to:      to,
		content: content,
	}
}

func (s *Email) From(from string) {
	s.from = from
}

func (s *Email) ContentType(contentType string) {
	s.contentType = contentType
}

func (s *Email) Content(content string) {
	s.content = content
}

func (s *Email) Cc(cc string) {
	s.cc = cc
}

func (s *Email) Bcc(bcc string) {
	s.bcc = bcc
}

func (s *Email) Attaches(attach ...string) {
	s.attaches = append(s.attaches, attach...)
}
