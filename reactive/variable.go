package reactive

import (
	"github.com/opoccomaxao-go/event/v3"
)

type Variable[T comparable] interface {
	typedVariable[T]

	OnChange(func(T)) event.Subscriber
	Repeat()
}

var _ Variable[struct{}] = (*variable[struct{}])(nil)

type variable[T comparable] struct {
	typedVariableImpl[T]

	event event.Event[T]
}

func (v *variable[T]) OnChange(listener func(T)) event.Subscriber {
	return v.event.Subscribe(listener)
}

func (v *variable[T]) Repeat() {
	v.event.Publish(v.value.Load().(T))
}

func (v *variable[T]) Set(value T) {
	prev := v.value.Swap(value)
	val, _ := prev.(T)
	if val != value {
		v.event.Publish(value)
	}
}

func (v *variable[T]) init() {
	v.Set(getZeroValue[T]())
}

func newVariable[T comparable](name string, prefix string) variable[T] {
	return variable[T]{
		typedVariableImpl: typedVariableImpl[T]{
			commonVariableImpl: commonVariableImpl{
				name:   name,
				prefix: prefix,
			},
		},
		event: event.NewEvent[T](),
	}
}

func NewVariable[T comparable](name string) Variable[T] {
	res := newVariable[T](name, getPrefix[T]())

	res.init()

	return &res
}
