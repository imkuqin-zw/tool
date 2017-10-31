package order

import "sort"

type KVMap struct {
	count	int
	keys	[]string
	hash	map[string]interface{}
}

func NewKVMap(iniCapacity uint) *KVMap {
	return &KVMap{
		count: 0,
		keys: make([]string, 0, int(iniCapacity)),
		hash: make(map[string]interface{}, int(iniCapacity)),
	}
}

func (kv *KVMap) Set(key string, val interface{}) {
	if _, ok := kv.hash[key]; !ok {
		kv.keys = append(kv.keys, key)
		kv.count++
	}
	kv.hash[key] = val
}

func (kv *KVMap) Get(k string) (interface{}, bool) {
	v, ok := kv.hash[k]
	return v, ok
}

func (kv *KVMap) ExistKey(k string) (bool) {
	_, ok := kv.hash[k]
	return ok
}

func (kv *KVMap) Keys() ([]string) {
	return kv.keys
}

func (kv *KVMap) Values() ([]interface{}) {
	vals := make([]interface{},0, kv.count)
	for i := 0; i < kv.count ; i++ {
		vals = append(vals, kv.hash[kv.keys[i]])
	}
	return vals
}

func (kv *KVMap) GetAllKV() ([]string, []interface{}) {
	vals := make([]interface{},0, kv.count)
	for i := 0; i < kv.count ; i++ {
		vals = append(vals, kv.hash[kv.keys[i]])
	}
	return kv.keys, vals
}

func (kv *KVMap) Del(key string) {
	for i, k := range kv.keys {
		if k == key {
			kv.keys = append(kv.keys[:i], kv.keys[i+1:]...)
			delete(kv.hash, key)
			kv.count--
		}
	}
}

func (kv *KVMap) DelAll() {
	kv.keys = make([]string, 0)
	kv.hash = make(map[string]interface{}, 0)
	kv.count = 0
}

func (kv *KVMap) Count() int {
	return kv.count
}

func (kv *KVMap) SortAsc() {
	sort.Strings(kv.keys)
}

func (kv *KVMap) SortDesc() {
	sort.Sort(sort.Reverse(sort.StringSlice(kv.keys)))
}