package sse

type Message struct {
	Event string
	Data  string
}

func NewMessage(data string) *Message {
	return &Message{
		"message",
		data,
	}
}

func NewEventMessage(event, data string) *Message {
	return &Message{
		event,
		data,
	}
}

func (m *Message) Encode() string {
	var msg string
	if m.Event != "" {
		msg += "event: " + m.Event + "\n"
	}
	msg += "data: " + m.Data + "\n\n"
	return msg
}
