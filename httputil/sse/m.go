package sse

import (
	"github.com/obnahsgnaw/goutils/strutil"
	"strconv"
)

type Message struct {
	Event string
	Data  string
	Id    string
	Retry int
}

func NewMessage(data string) *Message {
	return &Message{
		Event: "message",
		Data:  data,
	}
}

func NewEventMessage(event, data string) *Message {
	return &Message{
		Event: event,
		Data:  data,
	}
}

func Encode(m *Message) string {
	var msg string
	if m.Event != "" {
		msg = strutil.ToString("event: ", m.Event, "\n")
	}
	msg += strutil.ToString("data: ", m.Data)
	if m.Id != "" {
		msg += strutil.ToString("id: ", m.Id, "\n")
	}
	if m.Retry != 0 {
		msg += strutil.ToString("retry: ", strconv.Itoa(m.Retry), "\n")
	}
	msg += "\n\n"
	return msg
}
