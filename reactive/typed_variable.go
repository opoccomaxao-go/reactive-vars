package reactive

type typedVariable[T comparable] interface {
	CommonVariable

	Get() T
	Set(T)
}

type typedVariableImpl[T comparable] struct {
	commonVariableImpl
}

var _ typedVariable[struct{}] = (*typedVariableImpl[struct{}])(nil)

func (v *typedVariableImpl[T]) Get() T {
	return v.value.Load().(T)
}

func (v *typedVariableImpl[T]) Set(value T) {
	v.value.Store(value)
}
