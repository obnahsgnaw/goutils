package singletonutil

import "testing"

type DemoA struct {
}

func (s *DemoA) Name() string {
	return "this is demo a"
}

type DemoB struct {
}

func (s *DemoB) Name() string {
	return "this is demo b"
}

func TestNewBuilder(t *testing.T) {
	b1 := NewBuilder(func() interface{} {
		return &DemoA{}
	})
	b2 := NewBuilder(func() interface{} {
		return &DemoB{}
	})
	da, ok := b1.Get().(*DemoA)
	if !ok {
		t.Errorf("builder 1 need  reture demo a but got %v", da)
		return
	}
	if da.Name() != "this is demo a" {
		t.Errorf("builder 1  reture not got")
		return
	}
	db, ok := b2.Get().(*DemoB)
	if !ok {
		t.Errorf("builder 2 need  reture demo b but got %v", db)
		return
	}

	if db.Name() != "this is demo b" {
		t.Errorf("builder 2  reture not got")
		return
	}
}
