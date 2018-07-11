package excel

import (
	"testing"
	"fmt"
	"strconv"
	"github.com/imkuqin-zw/tool/order"
)
type GrandChild struct {
	Time int	   `excel:"cellName(时间3);func(TimeFormat, 2006-01-02)"`
}
type Children struct {
	Time int	   `excel:"cellName(时间2);func(TimeFormat, 2006-01-02 15:04:05)"`
	GrandChild *GrandChild `excel:"cellName(struct)"`
}
type StructTest struct {
	IntVal     int     `excel:"cellName(整数)"`
	StringVal  string  `excel:"cellName(-)"`
	Time 	   int64	   `excel:"cellName(时间);func(TimeFormat, 2006-01-02 15:04:05)"`
	Time2	   *Children `excel:"cellName(struct)"`
}

func TestCreateHeadByStruct(t *testing.T) {
	file, err := CreateHeadByStruct(new(StructTest), "test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func TestCreateByStructs(t *testing.T) {
	test := []StructTest{
		StructTest{
			IntVal:    16,
			StringVal: "heyheyhey :)!",
			Time: 1509419761,
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
		StructTest{
			IntVal:    45,
			StringVal: "fdsaf :)!",
			Time: 1509419761,
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
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
	err = file.Save("MyXLSXFile2.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func TestCreateByMapStructs(t *testing.T) {
	kvMap := order.NewKVMap(2)

	kvMap.Set("CPSA",[]interface{}{
		StructTest{
			IntVal:    16,
			StringVal: "heyheyhey :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
		StructTest{
			IntVal:    45,
			StringVal: "fdsaf :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
	})
	kvMap.Set("CPSG",[]interface{}{
		StructTest{
			IntVal:    17,
			StringVal: "heyheyhey :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
		StructTest{
			IntVal:    45,
			StringVal: "fdsaf :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
	})

	err, file := CreateByKVMapStructs(kvMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func TestCreateByKVLinkStructs(t *testing.T) {
	kvLink := order.NewKVLink()

	kvLink.Add2Last("CPSA",[]interface{}{
		StructTest{
			IntVal:    16,
			StringVal: "heyheyhey :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
		StructTest{
			IntVal:    45,
			StringVal: "fdsaf :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
	})
	kvLink.Add2Last("CPSG",[]interface{}{
		StructTest{
			IntVal:    79,
			StringVal: "heyheyhey :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
		StructTest{
			IntVal:    478,
			StringVal: "fdsaf :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
	})

	err, file := CreateByKVLinkStructs(kvLink)
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

	file, err := CreateByKVLink(kv, data,"")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("MyXLSXFile.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}

func TestCreateByKVMapInterface(t *testing.T) {
	KvMap := order.NewKVMap(3)
	KvMap.Set("CPSA",[]interface{}{
		StructTest{
			IntVal:    16,
			StringVal: "heyheyhey :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
		StructTest{
			IntVal:    45,
			StringVal: "fdsaf :)!",
			Time2: &Children{
				Time: 1509419761,
				GrandChild: &GrandChild{Time:1509419761},
			},
		},
	})

	kv := order.NewKVLink()
	kv.Add2Head("filed_1", "第一列")
	kv.Add2Last("filed_2", "第二列")
	KvMap.Set("CGD", kv)


	data := make([]*order.KVMap, 0, 21)
	kv2 := order.NewKVMap(2)
	kv2.Set("filed_1", "第一列")
	kv2.Set("filed_2", "第二列")
	keys := kv2.Keys()
	data = append(data, kv2)
	temps := make([]map[string]interface{}, 0, 50)
	for i := 0; i < 20 ; i++ {
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
	KvMap.Set("fdsa", data)

	err, file := CreateByKVMapInterface(KvMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("MyXLSXFile3.xlsx")
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

	file, err := CreateByKVMap(data, "fdsaf")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = file.Save("kvmap.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}
}