package sse

import "github.com/obnahsgnaw/goutils/strutil"

type Message struct {
	Event string
	Data  string
	pairs []string
}

func NewMessage(data string, ext ...string) *Message {
	if len(ext)%2 != 0 {
		ext = append(ext, "")
	}
	return &Message{
		Event: "message",
		Data:  data,
		pairs: ext,
	}
}

func NewEventMessage(event, data string, ext ...string) *Message {
	if len(ext)%2 != 0 {
		ext = append(ext, "")
	}
	return &Message{
		Event: event,
		Data:  data,
		pairs: ext,
	}
}

func Encode(m *Message) string {
	var msg string
	if m.Event != "" {
		msg = strutil.ToString("event: ", m.Event, "\n")
	}
	msg += strutil.ToString("data: ", m.Data)
	for i := 0; i < len(m.pairs); i += 2 {
		msg += m.pairs[i] + m.pairs[i+1]
	}
	msg += "\n\n"
	return msg
}
