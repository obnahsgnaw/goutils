package channels

import (
	"github.com/obnahsgnaw/goutils/notifyutil"
	"log"
)

const LogChannel notifyutil.ChannelName = "log"

type LogChannelHandler struct {
}

func NewLogChannel() *LogChannelHandler {
	return &LogChannelHandler{}
}

func (s *LogChannelHandler) Send(to notifyutil.Target, data notifyutil.Data) error {
	log.Println("Log channel: ", to, ">>>", data.Content())
	return nil
}
