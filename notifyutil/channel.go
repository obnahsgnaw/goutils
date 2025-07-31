package notifyutil

type ChannelName string

func (channel ChannelName) String() string {
	return string(channel)
}

type Target string

func (target Target) String() string {
	return string(target)
}

type Data interface {
	Subject() string
	Content() string
}

type StrData struct {
	subject string
	data    string
}

func NewStrData(subject, data string) *StrData {
	return &StrData{subject: subject, data: data}
}

func (s *StrData) Content() string {
	return s.data
}

func (s *StrData) Subject() string {
	return s.subject
}

type Channel interface {
	Send(to Target, data Data) error
}
