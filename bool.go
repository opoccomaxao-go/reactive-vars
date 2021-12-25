package reactive

type Bool interface {
	Variable
	Get() bool
	Set(bool)
}
type boolean struct {
	variable
}

func (b *boolean) Get() bool {
	res, _ := b.value.Load().(bool)
	return res
}

func (b *boolean) Set(value bool) {
	b.set(value)
}

func (boolean) Prefix() string {
	return "b:"
}
