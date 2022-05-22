package reactive

import (
	"github.com/opoccomaxao-go/event/v3"
)

type SyncVariable[T comparable] interface {
	typedVariable[T]

	// GetSyncPtr returns pointer to synchronous value. Do not change value, read-only use.
	GetSyncPtr() *T
	// OnChange sends to listener pointer to synchronous value. Do not change listeners value, read-only use.
	OnChange(listener func(*T)) event.Subscriber
	Repeat()
}

type syncVariable[T comparable] struct {
	typedVariableImpl[T]

	syncValue T
	event     event.Event[*T]
}

var _ SyncVariable[struct{}] = (*syncVariable[struct{}])(nil)

func (v *syncVariable[T]) GetSyncPtr() *T {
	return &v.syncValue
}

func (v *syncVariable[T]) OnChange(listener func(*T)) event.Subscriber {
	return v.event.Subscribe(listener)
}

func (v *syncVariable[T]) Repeat() {
	v.event.Publish(&v.syncValue)
}

func (v *syncVariable[T]) Set(value T) {
	prev := v.value.Swap(value)
	v.syncValue, _ = prev.(T)
	if v.syncValue != value {
		v.event.Publish(&v.syncValue)
	}
}

func (v *syncVariable[T]) init() {
	v.Set(v.syncValue)
}

func NewSyncVariable[T comparable](name string) SyncVariable[T] {
	res := &syncVariable[T]{
		typedVariableImpl: typedVariableImpl[T]{
			commonVariableImpl: commonVariableImpl{
				name:   name,
				prefix: getPrefix[T](),
			},
		},
		event: event.NewEvent[*T](),
	}

	res.init()

	return res
}
