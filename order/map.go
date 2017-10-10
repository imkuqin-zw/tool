package order

import (
	"errors"
	"sort"
)

var ErrNoData = errors.New("[Excel] no data found")

func GetOrderKey(data map[string]interface{}) (keys []string, err error){
	if len(data) == 0 {
		return nil, ErrNoData
	}
	keys = []string{}
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}