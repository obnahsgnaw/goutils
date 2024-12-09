package redis

import (
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func TestNewBuilder(t *testing.T) {
	rds := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "20240101123456",
		DB:       0,
	})
	b := New(rds)

	l1 := b.Compete("test1", 25*time.Second)

	if l1.Error() != nil {
		t.Fatal(l1.Error())
	}
	if l1.Exist() {
		t.Fatal("l1 need not exist but not")
	}

	l2 := b.Compete("test1", 25*time.Second)
	if l2.Error() != nil {
		t.Fatal(l2.Error())
	}
	if !l2.Exist() {
		t.Fatal("l2 need exist but not")
	}

	l1.Unlock()

	l3 := b.Compete("test1", 5*time.Second)
	if l3.Error() != nil {
		t.Fatal(l3.Error())
	}
	if l3.Exist() {
		t.Error("l3 need not exist but not")
		return
	}

	time.Sleep(5 * time.Second)

	l4 := b.Compete("test1", 5*time.Second)
	if l4.Error() != nil {
		t.Fatal(l4.Error())
	}
	if l4.Exist() {
		t.Error("l4 need ok after l3 expired, but not")
		return
	}
}
