package channels

import (
	"github.com/obnahsgnaw/goutils/httputil/sse"
	"github.com/obnahsgnaw/goutils/notifyutil"
)

const SseGroupChannel notifyutil.ChannelName = "sse-group"

type SseGroupChannelHandler struct {
	m *sse.Manager
}

func NewSseGroupChannel(m *sse.Manager) *SseGroupChannelHandler {
	return &SseGroupChannelHandler{m}
}

func (s *SseGroupChannelHandler) Send(to notifyutil.Target, data notifyutil.Data) error {
	s.m.BroadcastGroup(to.String(), sse.NewMessage(data.Content()))
	return nil
}
