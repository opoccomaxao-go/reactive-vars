package reactive

import (
	"log"
	"os"
	"sync"

	"github.com/opoccomaxao-go/event"
)

type Registry interface {
	Float(name string) Float
	Bool(name string) Bool
	Dump() map[string]interface{}
	Debug()
}

type registry struct {
	mu     sync.Mutex
	config Config
	ee     *event.Emitter
	vars   []Variable
	log    func(string, interface{})
	id     rune
}

type Config struct {
	Name string
}

func discard(string, interface{}) {}

func New(config Config) Registry {
	return &registry{
		config: config,
		ee:     event.NewEmitter(),
		vars:   make([]Variable, 0, 100),
		log:    discard,
	}
}

func (r *registry) Float(name string) Float {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.vars {
		if v, ok := v.(*float); ok && v.name == name {
			return v
		}
	}

	res := &float{
		variable: variable{
			ee:   r.ee,
			name: name,
			log:  r.log,
			id:   string(r.id),
		},
	}

	r.vars = append(r.vars, res)
	r.id++

	res.value.Store(float64(0))

	return res
}

func (r *registry) Bool(name string) Bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, v := range r.vars {
		if v, ok := v.(*boolean); ok && v.name == name {
			return v
		}
	}

	res := &boolean{
		variable: variable{
			ee:   r.ee,
			name: name,
			log:  r.Log,
			id:   string(r.id),
		},
	}

	r.vars = append(r.vars, res)
	r.id++

	res.value.Store(false)

	return res
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
		res[v.Name()] = v.Value()
	}

	return res
}
