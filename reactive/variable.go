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
	Repeat()
}

type Variable[T comparable] interface {
	сommonVariable
	OnChange(func(T)) event.Subscriber
	Get() T
	Set(T)
}

var _ Variable[struct{}] = (*variable[struct{}])(nil)

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

func (v *variable[T]) Repeat() {
	v.event.Publish(v.value.Load().(T))
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
		v.event.Publish(value)
		go v.log(v.FullName(), value)
	}
}

func NewVariable[T comparable](name string) Variable[T] {
	res := &variable[T]{
		event:  event.NewEvent[T](),
		name:   name,
		prefix: getPrefix[T](),
		log:    discard,
	}

	res.Set(getZeroValue[T]())

	return res
}
