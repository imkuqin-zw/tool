package reflectfunc

import (
	"reflect"
	"fmt"
)

// valid functions map.
type reflectFunc map[string]reflect.Value

func (f reflectFunc) Invoke(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := f[name]; !ok {
		err = NotFuncErr(fmt.Errorf("%s does not exist", name))
		return
	}

	if len(params) != f[name].Type().NumIn() {
		err = ParamsErr(fmt.Errorf("the number of parameters is not acceptable"))
		return
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f[name].Call(in)
	return
}

func (f reflectFunc) GetFuncByName(name string) (result reflect.Value, exist bool) {
	result, exist = f[name]
	return
}

func (f reflectFunc) GetParamsNumByName(name string) int {
	if _, ok := f[name]; !ok {
		return -1
	}
	return f[name].Type().NumIn() - 1
}

// Get a ReflectFunc.
func NewReflectFunc(obj interface{}, jump ...string) ReflectFunc {
	result := make(reflectFunc)
	t := reflect.TypeOf(obj)
	jumpMap := make(map[string]bool)
	for _, name := range jump {
		jumpMap[name] = true
	}
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if _, ok := jumpMap[m.Name]; ok {
			continue
		}
		result[m.Name] = m.Func
	}
	return result
}