package order

import (
	"testing"
	"fmt"
)

func TestKVLink(t *testing.T) {
	kv := NewKVLink()
	kv.Add("test1", "1231", 0)
	kv.Add("test2", "1232", 1)
	kv.Add("test3", "1233", 2)
	kv.Add("test4", "1234", 3)

	//kv.Traverse()
	//fmt.Println("--------------------------------------")
	//val, err := kv.Get("test1")
	//if err != nil {
	//	fmt.Println("GET:", err.Error())
	//	return
	//}
	//fmt.Printf("GET key: %v, val: %v \n", "test1", val)
	//fmt.Println("--------------------------------------")
	//if err = kv.Set("test1", "2135"); err != nil {
	//	fmt.Println("Set:", err.Error())
	//	return
	//}
	//
	//kv.Traverse()
	//fmt.Println("--------------------------------------")
	//kv.Del("test4")
	//kv.Traverse()
	//fmt.Println("--------------------------------------")
	//kv.Del("test1")
	//kv.Traverse()
	//fmt.Println("--------------------------------------")
	//kv.Del("test2")
	//kv.Traverse()
	//fmt.Println("--------------------------------------")
	//kv.Del("test3")
	//kv.Traverse()
	//fmt.Println("--------------------------------------")

	kv.Traverse()
	fmt.Println("--------------------------------------")
	kv.DelAll()
	kv.Traverse()
	fmt.Println("--------------------------------------")
}


