package reactive

import (
	"sync"

	"golang.org/x/exp/slices"
)

type SyncVariable[T comparable] interface {
	Variable[T]

	// Sync keep value sync with internal value.
	Sync(*T)
	Unsync(*T)
}

type syncVariable[T comparable] struct {
	variable[T]

	syncs []*T
	mu    sync.Mutex
}

var _ SyncVariable[struct{}] = (*syncVariable[struct{}])(nil)

func (v *syncVariable[T]) Sync(ptr *T) {
	if ptr == nil {
		return
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	// find duplicates
	idx := slices.Index(v.syncs, ptr)
	if idx >= 0 {
		return
	}

	v.syncs = append(v.syncs, ptr)
}

func (v *syncVariable[T]) Unsync(ptr *T) {
	if ptr == nil {
		return
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	idx := slices.Index(v.syncs, ptr)
	if idx == -1 {
		return
	}

	// this is used to atomic replace of v.syncs to remove any locks from syncListener.
	res := make([]*T, idx, len(v.syncs))
	copy(res, v.syncs[:idx])

	v.syncs = append(res, v.syncs[idx+1:]...)
}

func (v *syncVariable[T]) syncListener(value T) {
	syncs := v.syncs

	for _, ptr := range syncs {
		*ptr = value
	}
}

func (v *syncVariable[T]) init() {
	v.variable.init()
	v.OnChange(v.syncListener)
}

func NewSyncVariable[T comparable](name string) SyncVariable[T] {
	res := &syncVariable[T]{
		variable: newVariable[T](name, getPrefix[T]()),
	}

	res.init()

	return res
}
