package Excel

import (
	"testing"
	"fmt"
	"strconv"
	"tool/order"
)

type StructTest struct {
	IntVal     int     `excel:"cellName(整数)"`
	StringVal  string  `excel:"cellName(-)"`
}

func TestCreateByStructs(t *testing.T) {
	test := []StructTest{
		StructTest{
			IntVal:    16,
			StringVal: "heyheyhey :)!",
		},
		StructTest{
			IntVal:    45,
			StringVal: "fdsaf :)!",
		},
	}
	test2 := []interface{}{}
	for _, item := range test {
		test2 = append(test2, interface{}(item))
	}
	err, file := CreateByStructs(test2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func TestCreateByKVLink(t *testing.T) {
	data := []map[string]interface{}{}
	for i := 0; i < 50 ; i++ {
		temp := map[string]interface{}{
			"filed_1": "test"+ strconv.Itoa(i) + "_1",
			"filed_2": "test"+ strconv.Itoa(i) + "_2",
		}
		data = append(data, temp)
	}

	kv := order.NewKVLink()
	kv.Add2Head("filed_1", "第一列")
	kv.Add2Last("filed_2", "第二列")

	file, err := CreateByKVLink(kv, data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func TestCreateByKVMap(t *testing.T) {
	data := make([]*order.KVMap, 0, 51)
	kv := order.NewKVMap(2)
	kv.Set("filed_1", "第一列")
	kv.Set("filed_2", "第二列")
	keys := kv.Keys()
	data = append(data, kv)
	temps := make([]map[string]interface{}, 0, 50)
	for i := 0; i < 50 ; i++ {
		temp := map[string]interface{}{
			"filed_1": "test"+ strconv.Itoa(i) + "_1",
			"filed_2": "test"+ strconv.Itoa(i) + "_2",
		}
		temps = append(temps, temp)
	}
	for _, val := range temps{
		kv := order.NewKVMap(2)
		for _, k := range keys{
			kv.Set(k, val[k])
		}
		data = append(data, kv)
	}

	file, err := CreateByKVMap(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("kvmap.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}
