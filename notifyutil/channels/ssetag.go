package channels

import (
	"github.com/obnahsgnaw/goutils/httputil/sse"
	"github.com/obnahsgnaw/goutils/notifyutil"
)

const SseTagChannel notifyutil.ChannelName = "sse-tag"

type SseTagChannelHandler struct {
	m *sse.Manager
}

func NewSseTagChannel(m *sse.Manager) *SseTagChannelHandler {
	return &SseTagChannelHandler{m}
}

func (s *SseTagChannelHandler) Send(to notifyutil.Target, data notifyutil.Data) error {
	s.m.BroadcastTag(to.String(), sse.NewMessage(data.Content()))
	return nil
}
