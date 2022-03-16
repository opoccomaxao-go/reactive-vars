package reactive

import (
	"sync/atomic"

	"github.com/opoccomaxao-go/event/v3"
)

type сommonVariable interface {
	Value() interface{}
	Name() string
	Prefix() string
	FullName() string
}

type Variable[T comparable] interface {
	сommonVariable
	OnChange(func(T)) event.Subscriber
	Get() T
	Set(T)
}

type variable[T comparable] struct {
	event  event.Event[T]
	name   string
	prefix string
	value  atomic.Value
	log    func(name string, value interface{})
}

func (v *variable[T]) Value() interface{} {
	return v.value.Load()
}

func (v *variable[T]) Name() string {
	return v.name
}

func (v *variable[T]) Prefix() string {
	return v.prefix
}

func (v *variable[T]) FullName() string {
	return v.prefix + ":" + v.name
}

func (v *variable[T]) OnChange(fn func(T)) event.Subscriber {
	return v.event.Subscribe(fn)
}

func (v *variable[T]) Get() T {
	return v.value.Load().(T)
}

func (v *variable[T]) Set(value T) {
	prev := v.value.Swap(value)
	val, _ := prev.(T)
	if val != value {
		v.log(v.FullName(), value)
		v.event.Publish(value)
	}
}
