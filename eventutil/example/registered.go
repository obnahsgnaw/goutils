package example

import (
	"github.com/obnahsgnaw/goutils/eventutil"
	"log"
)

// 1. define
var _RegisteredEvent = newRegisteredEvent()

func init() {
	// 2. register
	_RegisteredEvent.RegisterTo(eventutil.DefaultManger())
}

// FireRegisteredEvent 3. fires with data
func FireRegisteredEvent(uid int) {
	_RegisteredEvent.Fire(uid)
}

type RegisteredEvent struct {
	eventutil.BaseEventBuilder
}

func newRegisteredEvent() *RegisteredEvent {
	e := &RegisteredEvent{}
	e.Initialize(e)
	return e
}

func (e *RegisteredEvent) Topic() eventutil.Topic {
	return "demo"
}

func (e *RegisteredEvent) Handle(evt *eventutil.Event) {
	// 4. parse data and handle
	uid := evt.Data.Get(0).(int)
	log.Println("User:", uid, " register successfully")
}
