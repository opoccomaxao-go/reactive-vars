package reactive

import (
	"sync/atomic"

	"github.com/opoccomaxao-go/event"
)

type Variable interface {
	OnChange(func())
	Name() string
	Value() interface{}
	Prefix() string
}

type variable struct {
	ee    *event.Emitter
	name  string
	id    string
	value atomic.Value
	log   func(name string, value interface{})
}

func (v *variable) set(value interface{}) {
	prev := v.value.Swap(value)
	if prev != value {
		v.log(v.name, value)
		v.ee.Emit(v.id)
	}
}

func (v *variable) OnChange(fn func()) {
	v.ee.On(v.id, func(...interface{}) { fn() })
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
