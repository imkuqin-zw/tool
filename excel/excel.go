package excel

import (
	"github.com/tealeg/xlsx"
	"reflect"
	"strings"
	"errors"
	"github.com/imkuqin-zw/tool/order"
	"fmt"
)

var ErrNoData = errors.New("[Excel] no data found")
var ErrParams = errors.New("[Excel] params error")
var ErrNotFunc = errors.New("[Excel] no func error")

//create excel file by struct array
func CreateByStructs(data []interface{}) (err error, file *xlsx.File ) {
	file = xlsx.NewFile()
	err = createSheetByStructs(file, "Sheet1", data)
	return
}

// 通过结构体KVMap创建多个sheet
func CreateByKVMapStructs(data *order.KVMap) (err error, file *xlsx.File) {
	if data.Count() == 0 {
		return ErrNoData, nil
	}
	file = xlsx.NewFile()
	keys, values := data.GetAllKV()
	for i, key := range keys {
		err = createSheetByStructs(file, key, values[i].([]interface{}))
		if err != nil {
			return
		}
	}

	return
}

// 通过结构体KVLink创建多个sheet
func CreateByKVLinkStructs(data *order.KVLinkHead) (err error, file *xlsx.File) {
	if data.GetCount() == 0 {
		return ErrNoData, nil
	}
	file = xlsx.NewFile()
	keys, values := data.GetAllKV()
	for i, key := range keys {
		err = createSheetByStructs(file, key, values[i].([]interface{}))
		if err != nil {
			return
		}
	}

	return
}

// 通过结构体创建表头
func CreateHeadByStruct(m interface{}, key string) (file *xlsx.File, err error) {
	file = xlsx.NewFile()
	sheet, err := file.AddSheet(key)
	if err != nil {
		return
	}
	t := reflect.TypeOf(m)
	row := sheet.AddRow()
	err , _ = createRowHeaderByStruct(row, t)
	return
}

// 通过KVMap中值的不同类型创建多个sheep
func CreateByKVMapInterface(data *order.KVMap) (err error, file *xlsx.File) {
	file = xlsx.NewFile()
	data.SortAsc()
	keys, values := data.GetAllKV()
	for i, key := range keys {
		t := reflect.TypeOf(values[i])
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		name := t.String()
		if name == "[]*order.KVMap" { //KVMap创建有数据的sheet
			err = createSheetByKVMap(file, key, values[i].([]*order.KVMap))
			if err != nil {
				return
			}
		} else if name == "order.KVLinkHead" { //创建没数据的sheet
			err = createSheetKVLink(file, key, values[i].(*order.KVLinkHead), []map[string]interface{}{})
			if err != nil {
				return
			}
		} else { //Struct创建有数据的表
			err = createSheetByStructs(file, key, values[i].([]interface{}))
			if err != nil {
				return
			}
		}
	}
	return
}

// 通过KVLink创建表
func CreateByKVLink(head *order.KVLinkHead, data []map[string]interface{}, key string) (file *xlsx.File, err error) {
	file = xlsx.NewFile()
	if key == "" {
		key = "sheet1"
	}
	err = createSheetKVLink(file, key, head, data)
	return
}

// 通过KVMap创建表
func CreateByKVMap(data []*order.KVMap, key string) (file *xlsx.File, err error) {
	file = xlsx.NewFile()
	if key == "" {
		key = "sheet1"
	}
	err = createSheetByKVMap(file, "sheet1", data)
	return
}

// 创建sheet通过结构体
func createSheetByStructs(file *xlsx.File, key string, values []interface{}) (err error) {
	if len(values) == 0 {
		return ErrNoData
	}
	sheet, err := file.AddSheet(key)
	if err != nil {
		return err
	}

	//header
	t := reflect.TypeOf(values[0])
	row := sheet.AddRow()
	err , exports := createRowHeaderByStruct(row, t)
	//data
	for i := 0; i < len(values); i++ {
		//val := reflect.ValueOf(values[i])
		row = sheet.AddRow()
		if err = createRowDataByStructs(row, reflect.ValueOf(values[i]), exports); err != nil {
			return
		}
	}

	return
}

// 通过结构体创建sheet头
func createRowHeaderByStruct(row *xlsx.Row, t reflect.Type) (err error, exports []interface{}) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return ErrParams, nil // Should not panic here ?
	}
	for i := 0; i < t.NumField(); i++ {
		err, cellName := getCellNameByTag(t.Field(i).Tag.Get("excel"))
		if err != nil {
			return err, nil
		}
		if cellName == "-" {
			continue
		}
		if cellName == "struct" {
			err, appendExports := createRowHeaderByStruct(row,t.Field(i).Type)
			if err != nil {
				return err, nil
			}
			exports = append(exports, map[int][]interface{}{i:appendExports})
		} else {
			exports = append(exports, i)
			row.AddCell().SetValue(cellName)
		}
	}
	return
}

// 通过结构体创建数据
func createRowDataByStructs(row *xlsx.Row, val reflect.Value, exports []interface{})  error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for _, index := range exports {
		indexType := reflect.TypeOf(index)
		if indexType.Kind() != reflect.Map {
			item := val.Field(index.(int))
			err, funcs, funcParams := getFuncByTag(val.Type().Field(index.(int)).Tag.Get("excel"))
			if err != ErrNotFunc {
				excelFunc := &ExcelFunc{}
				functions := reflect.ValueOf(excelFunc)
				for i, funcName := range funcs {
					params := make([]reflect.Value, 1)
					params[0] = item
					for _, param := range funcParams[i] {
						params = append(params, reflect.ValueOf(param))
					}
					result := functions.MethodByName(funcName).Call(params)
					if len(result) == 2 {
						if result[1].Interface() != nil {
							err = fmt.Errorf("[Excel %s: %s]", funcName, result[1].Interface().(error).Error())
							return err
						}
					}
					item = result[0]
				}
			}
			row.AddCell().SetValue(item.Interface())
		} else {
			for k, v := range index.(map[int][]interface{}) {
				if err := createRowDataByStructs(row, val.Field(k), v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 通过KVMap创建sheet
func createSheetByKVMap(file *xlsx.File, key string, data []*order.KVMap) ( err error){
	sheet, err := file.AddSheet(key)
	if err != nil {
		return
	}
	for _, item := range data {
		row := sheet.AddRow()
		for _, v := range item.Values() {
			row.AddCell().SetValue(v)
		}
	}

	return
}

// 通过KVLink创建sheet
func createSheetKVLink(file *xlsx.File, key string, head *order.KVLinkHead, data []map[string]interface{}) ( err error ) {
	if head.GetCount() == 0 {
		return ErrParams
	}
	sheet, err := file.AddSheet(key)
	if err != nil {
		return
	}
	keys, vals := head.GetAllKV()
	row := sheet.AddRow()
	for _, val := range vals {
		row.AddCell().SetValue(val)
	}
	for _, item := range data {
		row = sheet.AddRow()
		for _, key := range keys {
			row.AddCell().SetValue(item[key])
		}
	}
	return
}

// 解析tag，获得表头
func getCellNameByTag(tag string) (err error, cellName string) {
	for _, v := range strings.Split(tag, ";") {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v) -1 {
			if v[:i] == "cellName" {
				return nil, v[i+1:len(v) -1]
			}
		}
	}
	return errors.New("[Excel] not found cell name from tag"), ""
}

// 解析tag，获得处理该字段所需的函数名和参数
func getFuncByTag(tag string) (err error, funcName []string, params[][]string) {
	for _, item := range strings.Split(tag, ";") {
		if item == "" {
			continue
		}
		item = strings.TrimSpace(item)
		if i := strings.Index(item, "("); i > 0 && strings.Index(item, ")") == len(item) -1 {
			if item[:i] == "func" {
				for _, val := range strings.Split(item[i+1:len(item) -1], "|") {
					val = strings.TrimSpace(val)
					temps := strings.Split(val, ",")
					funcName = append(funcName, strings.TrimSpace(temps[0]))
					param := []string{}
					for _, v := range temps[1:] {
						param = append(param, strings.TrimSpace(v))
					}
					params = append(params, param)
				}

				return
			}
		}
	}

	return ErrNotFunc, nil, nil
}


