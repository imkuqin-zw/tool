package reflectfunc

import (
	"fmt"
	"testing"
	"reflect"
)

type reflectFuncTest int

func (c reflectFuncTest) Test() {
	fmt.Println("Test")
}

func (c reflectFuncTest) TestParams(in string) {
	fmt.Printf("TestParams %s\n", in)
}

func (c reflectFuncTest) TestResult(in string) string {
	return in
}

func TestInvokeFunc(t *testing.T) {
	var test reflectFuncTest = 8
	var err error
	var result []reflect.Value
	funcs := NewReflectFunc(test)
	_, err = funcs.invoke("Test", test)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = funcs.invoke("TestParams", test, "test")
	if err != nil {
		fmt.Println(err.Error())
	}
	result, err = funcs.invoke("TestResult", test, "TestResult")
	if err != nil {
		fmt.Println(err.Error())
	} else{
		fmt.Println(result[0].String())
	}
}
