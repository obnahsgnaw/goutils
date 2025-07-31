package example

import (
	"github.com/obnahsgnaw/goutils/notifyutil"
	"github.com/obnahsgnaw/goutils/notifyutil/channels"
	"log"
)

type RegisteredNotification struct {
	notifyutil.BaseNotificationBuilder
	user  string
	email string
}

func NewRegisteredNotification(user, email string) *RegisteredNotification {
	n := &RegisteredNotification{user: user, email: email}
	n.Initialize(n)
	return n
}

func (d RegisteredNotification) Channels() []notifyutil.ChannelName {
	return []notifyutil.ChannelName{channels.LogChannel, channels.EmailChannel}
}

func (d RegisteredNotification) To(notifyutil.ChannelName) notifyutil.Target {
	return notifyutil.Target(d.email)
}

func (d RegisteredNotification) Data(notifyutil.ChannelName) notifyutil.Data {
	return notifyutil.NewStrData("Register Successfully", d.user+", register successfully")
}

func (d RegisteredNotification) SuccessHandle(notifyutil.ChannelName) {
	log.Println("RegisteredNotification: handle success")
}

func (d RegisteredNotification) FailedHandle(notifyutil.ChannelName, error) {
	log.Println("RegisteredNotification: handle failed")
}

func (d RegisteredNotification) RetryMax() int {
	return 3
}
