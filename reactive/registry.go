package reactive

import (
	"log"
	"os"
	"sync"
)

type Registry interface {
	Float(name string) SyncVariable[float64]
	Bool(name string) SyncVariable[bool]
	Dump() map[string]interface{}
	GetAll() []CommonVariable
}

type registry struct {
	mu     sync.Mutex
	config Config
	vars   []CommonVariable
	logger *log.Logger
}

type Config struct {
	Name  string
	Debug bool
}

func New(config Config) Registry {
	return &registry{
		config: config,
		vars:   make([]CommonVariable, 0, 100),
		logger: log.New(os.Stdout, "", log.Flags()),
	}
}

func (r *registry) Find(prefix, name string) CommonVariable {
	for _, v := range r.vars {
		if v.Prefix() == prefix && v.Name() == name {
			return v
		}
	}

	return nil
}

func (r *registry) Float(name string) SyncVariable[float64] {
	return getVariable[float64](r, name)
}

func (r *registry) Bool(name string) SyncVariable[bool] {
	return getVariable[bool](r, name)
}

func (r *registry) log(variable CommonVariable) {
	r.logger.Printf("%s = %v\n", variable.FullName(), variable.Value())
}

func (r *registry) Dump() map[string]interface{} {
	res := map[string]interface{}{}

	for _, v := range r.vars {
		res[v.FullName()] = v.Value()
	}

	return res
}

func (r *registry) GetAll() []CommonVariable {
	return append([]CommonVariable{}, r.vars...)
}

func getLoggerListener[T comparable](r *registry, variable Variable[T]) func(T) {
	return func(T) {
		r.log(variable)
	}
}

func getVariable[T comparable](r *registry, name string) SyncVariable[T] {
	prefix := getPrefix[T]() + ":" + r.config.Name

	r.mu.Lock()
	defer r.mu.Unlock()

	{
		resVariable, ok := r.Find(prefix, name).(SyncVariable[T])
		if ok {
			return resVariable
		}
	}

	tempVariable := newSyncVariable[T](name, prefix)
	r.vars = append(r.vars, &tempVariable)

	tempVariable.init()
	if r.config.Debug {
		tempVariable.OnChange(getLoggerListener[T](r, &tempVariable)).Async()
	}

	return &tempVariable
}
