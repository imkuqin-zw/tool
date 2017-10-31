package order

import (
	"errors"
	"fmt"
)

type KVLinkHead struct {
	count	int
	first	*KVLinkData
	last	*KVLinkData
}

type KVLinkData struct {
	key		string
	val		interface{}
	next	*KVLinkData
}

var ErrNoKey = errors.New("[KVLink] No Key found")

func NewKVLink() *KVLinkHead {
	return &KVLinkHead{
		count: 0,
		first: nil,
		last: nil,
	}
}

func (kv *KVLinkHead) Add2Head(key string, val interface{}){
	kv.Add(key, val, 0)
}

func (kv *KVLinkHead) Add2Last(key string, val interface{}){
	kv.Add(key, val, -1)
}

func (kv *KVLinkHead) Add(key string, val interface{}, index int) {
	next := &KVLinkData{
		key: key,
		val: val,
		next: nil,
	}
	if kv.count == 0 {
		kv.first = next
		kv.last = next
	} else {
		if index < 0 || index >= kv.count {
			kv.last.next = next
			kv.last = next
		} else if index == 0 {
			next.next = kv.first
			kv.first = next
		} else {
			temp := kv.first
			for i := 1; i < index; i++ {
				temp = temp.next
			}
			next.next = temp.next
			temp.next = next
		}
	}
	kv.count++
	return
}

func (kv *KVLinkHead) Set(key string, val interface{}) error {
	temp := kv.first
	for i := 0; i < kv.count; i++ {
		if temp.key == key {
			temp.val = val
			return nil
		}
		temp = temp.next
	}

	return ErrNoKey
}

func (kv *KVLinkHead) Del(key string) error {
	var temp *KVLinkData
	for i := 0; i < kv.count; i++ {
		if i == 0 {
			if kv.first.key == key {
				kv.first = kv.first.next
				kv.count--
				if kv.count == 0 {
					kv.last = kv.first
				}
				return nil
			}
			temp = kv.first
		} else {
			if temp.next.key == key {
				temp.next = temp.next.next
				if kv.count == i+1 {
					kv.last = temp
				}
				kv.count--
				return nil
			}
			temp = temp.next
		}
	}

	return ErrNoKey
}

func (kv *KVLinkHead) DelAll() {
	kv.first = nil
	kv.last = nil
	kv.count = 0
}

func (kv *KVLinkHead) Get(key string) (interface{}, error) {
	temp := kv.first
	for i := 0; i < kv.count; i++ {
		if temp.key == key {
			return temp.val, nil
		}
		temp = temp.next
	}

	return nil, ErrNoKey
}

func (kv *KVLinkHead) GetLast() (interface{}, interface{}, error) {
	if kv.last == nil {
		return nil, nil, ErrNoData
	}
	return kv.last.key, kv.last.val, nil
}

func (kv *KVLinkHead) GetFirst() (interface{}, interface{}, error) {
	if kv.first == nil {
		return nil, nil, ErrNoData
	}
	return kv.first.key, kv.first.val, nil
}

func (kv *KVLinkHead) GetAllKV() ([]string, []interface{}) {
	temp := kv.first
	keys := make([]string, 0,  kv.count)
	vals := make([]interface{}, 0, kv.count)
	for i := 0; i < kv.count; i++ {
		keys = append(keys, temp.key)
		vals = append(vals, temp.val)
		temp = temp.next
	}
	return keys, vals
}

func (kv *KVLinkHead) GetCount() int {
	return kv.count
}

func (kv *KVLinkHead) Traverse() {
	temp := kv.first
	for i := 0; i < kv.count; i++ {
		fmt.Printf("key: %v, val: %v \n", temp.key, temp.val)
		temp = temp.next
	}
}
