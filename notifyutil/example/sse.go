package example

import (
	"github.com/obnahsgnaw/goutils/notifyutil"
	"github.com/obnahsgnaw/goutils/notifyutil/channels"
	"log"
)

type SseNotification struct {
	notifyutil.BaseNotificationBuilder
}

func NewSseNotification() *SseNotification {
	n := &SseNotification{}
	n.Initialize(n)
	n.RegisterTo(notifyutil.DefaultManager())
	return n
}

func (d SseNotification) Channels() []notifyutil.ChannelName {
	return []notifyutil.ChannelName{channels.SseGroupChannel, channels.SseTagChannel}
}

func (d SseNotification) To(cn notifyutil.ChannelName) notifyutil.Target {
	switch cn {
	case channels.SseTagChannel:
		return "tag1"
	case channels.SseGroupChannel:
		return "group1"
	default:
		return ""
	}
}

func (d SseNotification) Data(notifyutil.ChannelName) notifyutil.Data {
	return notifyutil.NewStrData("Demo", "this is sse message")
}

func (d SseNotification) SuccessHandle(notifyutil.ChannelName) {
	log.Println("sse notification: handle success")
}

func (d SseNotification) FailedHandle(notifyutil.ChannelName, error) {
	log.Println("sse notification: handle failed")
}

func (d SseNotification) RetryMax() int {
	return 3
}
