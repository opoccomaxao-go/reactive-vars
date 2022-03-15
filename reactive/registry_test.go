package reactive

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	r := New(Config{Name: "test"})
	r.Debug()

	changes := []string{}
	bindUpdateChanges := func(change string) func(interface{}) {
		return func(interface{}) {
			changes = append(changes, change)
		}
	}

	b := r.Bool("a")
	f := r.Float("a")

	b.OnChange(bindUpdateChanges("b"))
	f.OnChange(bindUpdateChanges("f"))

	b.Set(true)
	time.Sleep(time.Millisecond)
	b.Set(false)
	time.Sleep(time.Millisecond)
	f.Set(1)
	time.Sleep(time.Millisecond)
	f.Set(100)
	time.Sleep(time.Millisecond)

	f2 := r.Float("a")
	assert.True(t, f == f2)
	assert.Equal(t, float64(100), f2.Get())

	f2.Set(0)
	time.Sleep(time.Millisecond)

	assert.Equal(t, []string{"b", "b", "f", "f", "f"}, changes)

	r.Bool("b").Set(true)

	assert.Equal(t, map[string]interface{}{
		"f:a": float64(0),
		"b:a": false,
		"b:b": true,
	}, r.Dump())
}
