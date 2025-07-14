package sse

type Message struct {
	Event string
	Data  string
}

func (m *Message) Encode() string {
	var msg string
	if m.Event != "" {
		msg += "event: " + m.Event + "\n"
	}
	msg += "data: " + m.Data + "\n\n"
	return msg
}
