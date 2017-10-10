package order

import (
	"testing"
	"strconv"
	"fmt"
)

func TestKVMap(t *testing.T) {
	kv := NewKVMap(5)
	for i := 1; i < 5; i++ {
		kv.Set(strconv.Itoa(i), i)
	}
	kv.Set("0", 0)
	trasver(kv)
	v, ok := kv.Get("0")
	if ok == false {
		fmt.Println("不存在")
		fmt.Println("--------------------------------------")
		return
	}
	fmt.Printf("key: %v, val: %v\n", 0, v)
	fmt.Println("--------------------------------------")
	kv.SortAsc()
	trasver(kv)
	kv.SortDesc()
	trasver(kv)
	kv.Set("0", 1)
	trasver(kv)
	kv.DelAll()
	trasver(kv)
}

func trasver(kv *KVMap)  {
	Key, val := kv.GetAllKV()
	for i, k := range Key {
		fmt.Printf("key: %v, val: %v\n", k, val[i])
	}
	fmt.Println("--------------------------------------")
}