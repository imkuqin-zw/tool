package reflectfunc

import (
	"reflect"
)

// Not found the function error.
type NotFuncErr error

// Function parameter error.
type ParamsErr error

// reflect invoke function interface.
type ReflectFunc interface{
	// Invoke a function by name.
	Invoke(name string, params ...interface{}) (result []reflect.Value, err error)

	// Get the number of the function`s parameters by function name.
	GetParamsNumByName(name string) int

	// Get the function by name.
	GetFuncByName(name string) (reflect.Value, bool)
}