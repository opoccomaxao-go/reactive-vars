package reactive

import "sync/atomic"

type commonVariable interface {
	Value() interface{}
	Name() string
	Prefix() string
	FullName() string
}

type commonVariableImpl struct {
	name   string
	prefix string
	value  atomic.Value
}

var _ commonVariable = (*commonVariableImpl)(nil)

func (v *commonVariableImpl) Value() interface{} {
	return v.value.Load()
}

func (v *commonVariableImpl) Name() string {
	return v.name
}

func (v *commonVariableImpl) Prefix() string {
	return v.prefix
}

func (v *commonVariableImpl) FullName() string {
	return v.prefix + ":" + v.name
}
