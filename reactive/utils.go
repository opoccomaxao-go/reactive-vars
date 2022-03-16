package reactive

import (
	"reflect"
)

func getPrefix[T any]() string {
	var temp T

	return reflect.TypeOf(temp).Name()
}

func getZeroValue[T any]() (res T) {
	return
}
