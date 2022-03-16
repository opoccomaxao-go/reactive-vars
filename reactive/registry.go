package reactive

import (
	"log"
	"os"
	"sync"

	"github.com/opoccomaxao-go/event/v3"
)

type Registry interface {
	Float(name string) Variable[float64]
	Bool(name string) Variable[bool]
	Dump() map[string]interface{}
	Debug()
}

type registry struct {
	mu     sync.Mutex
	config Config
	vars   []сommonVariable
	log    func(string, interface{})
}

type Config struct {
	Name string
}

func discard(string, interface{}) {}

func New(config Config) Registry {
	return &registry{
		config: config,
		vars:   make([]сommonVariable, 0, 100),
		log:    discard,
	}
}

func (r *registry) Find(prefix, name string) сommonVariable {
	for _, v := range r.vars {
		if v.Prefix() == prefix && v.Name() == name {
			return v
		}
	}

	return nil
}

func (r *registry) Float(name string) Variable[float64] {
	return getVariable[float64](r, name)
}

func (r *registry) Bool(name string) Variable[bool] {
	return getVariable[bool](r, name)
}

func (r *registry) Debug() {
	logger := log.New(os.Stdout, "", log.Flags())
	r.log = func(name string, value interface{}) {
		logger.Printf("%s:%s = %v\n", r.config.Name, name, value)
	}
}

func (r *registry) Log(name string, value interface{}) {
	r.log(name, value)
}

func (r *registry) Dump() map[string]interface{} {
	res := map[string]interface{}{}

	for _, v := range r.vars {
		res[v.FullName()] = v.Value()
	}

	return res
}

func getVariable[T comparable](r *registry, name string) Variable[T] {
	prefix := getPrefix[T]()

	r.mu.Lock()
	defer r.mu.Unlock()

	res, ok := r.Find(prefix, name).(Variable[T])
	if ok {
		return res
	}

	res = &variable[T]{
		event:  event.NewEvent[T](),
		prefix: prefix,
		name:   name,
		log:    r.Log,
	}

	res.Set(getZeroValue[T]())

	r.vars = append(r.vars, res)

	return res
}
