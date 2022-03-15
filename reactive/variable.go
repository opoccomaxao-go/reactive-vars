package reactive

import (
	"sync/atomic"

	"github.com/opoccomaxao-go/event/v2"
)

type Variable interface {
	OnChange(func(interface{}))
	Name() string
	Value() interface{}
	Prefix() string
}

type variable struct {
	event event.Event
	name  string
	value atomic.Value
	log   func(name string, value interface{})
}

func (v *variable) set(value interface{}) {
	prev := v.value.Swap(value)
	if prev != value {
		v.log(v.name, value)
		v.event.Publish(value)
	}
}

func (v *variable) OnChange(fn func(interface{})) {
	v.event.Subscribe(fn)
}

func (v *variable) Name() string {
	return v.name
}

func (v *variable) Value() interface{} {
	return v.value.Load()
}

func (variable) Prefix() string {
	return ""
}
