package example

import (
	"github.com/obnahsgnaw/goutils/notifyutil"
	"github.com/obnahsgnaw/goutils/notifyutil/channels"
	"log"
)

type DemoNotification struct{}

func (d DemoNotification) Channels() []notifyutil.ChannelName {
	return []notifyutil.ChannelName{channels.LogChannel}
}

func (d DemoNotification) To(notifyutil.ChannelName) notifyutil.Target {
	return "test"
}

func (d DemoNotification) Data(notifyutil.ChannelName) notifyutil.Data {
	return notifyutil.NewStrData("Demo", "this is demo")
}

func (d DemoNotification) SuccessHandle(notifyutil.ChannelName) {
	log.Println("demo: handle success")
}

func (d DemoNotification) FailedHandle(notifyutil.ChannelName, error) {
	log.Println("demo: handle failed")
}

func (d DemoNotification) RetryMax() int {
	return 3
}
