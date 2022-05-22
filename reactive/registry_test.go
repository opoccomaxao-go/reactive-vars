package reactive

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func bindChanges[T any](dest *[]string) func(T) {
	return func(T) {
		*dest = append(*dest, getPrefix[T]())
	}
}

func TestRegistry(t *testing.T) {
	r := New(Config{
		Name:  "test",
		Debug: true,
	})

	changes := []string{}

	b := r.Bool("a")
	f := r.Float("a")

	b.OnChange(bindChanges[bool](&changes))
	f.OnChange(bindChanges[float64](&changes))

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

	assert.Equal(t, []string{
		"bool",
		"bool",
		"float64",
		"float64",
		"float64",
	}, changes)

	r.Bool("b").Set(true)

	assert.Equal(t, map[string]interface{}{
		"float64:a": float64(0),
		"bool:a":    false,
		"bool:b":    true,
	}, r.Dump())
}
