package channels

import (
	"github.com/obnahsgnaw/goutils/emailutil"
	"github.com/obnahsgnaw/goutils/notifyutil"
)

const EmailChannel notifyutil.ChannelName = "email"

type EmailChannelHandler struct {
	manager *emailutil.Manager
}

func NewEmailChannel(m *emailutil.Manager) *EmailChannelHandler {
	return &EmailChannelHandler{m}
}

func (s *EmailChannelHandler) Send(to notifyutil.Target, data notifyutil.Data) error {
	s.manager.Build(data.Subject(), data.Content(), to.String())
	return s.manager.Send()
}
