package eventutil

type Topic string

func (t Topic) String() string {
	return string(t)
}

type Handler func(*Event)

type Data []interface{}

func (d Data) Get(i int) interface{} {
	if i < 0 || i >= len(d) {
		return nil
	}
	return d[i]
}

// Event the event target
type Event struct {
	Topic    Topic
	Data     Data
	provider func() *Manger
}

func newEvent(topic Topic, data Data, manager func() *Manger) *Event {
	return &Event{
		Topic:    topic,
		Data:     data,
		provider: manager,
	}
}

// Fire to fire event
func (e *Event) Fire() {
	if e.provider != nil {
		e.provider().Fire(e)
	}
}
