package reactive

type Float interface {
	Variable
	Get() float64
	Set(float64)
}

type float struct {
	variable
}

func (f *float) Get() float64 {
	res, _ := f.value.Load().(float64)
	return res
}

func (f *float) Set(value float64) {
	f.set(value)
}
