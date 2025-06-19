package structutil

import (
	"github.com/obnahsgnaw/goutils/cacheutil/static"
	"testing"
	"time"
)

type User struct {
	SyncStruct
	Id       int    `json:"id"`
	Username string `json:"username"`
}

var c = static.New()

func NewUser() *User {
	s := &User{}
	s.Init(c, "user", time.Second*10, s)
	return s
}

func TestSyncStruct(t *testing.T) {
	user1 := NewUser()
	user1.Id = 1
	user1.Username = "user1"
	err := user1.Save("1")
	if err != nil {
		t.Fatal(err)
		return
	}

	user2 := NewUser()
	hit, err := user2.Load("1")
	if err != nil {
		t.Fatal(err)
		return
	}
	if !hit {
		t.Fatal("not hit")
		return
	}
	println(user2.Id)
	println(user2.Username)
}
