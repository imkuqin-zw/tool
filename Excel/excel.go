package Excel

import (
	"github.com/tealeg/xlsx"
	"reflect"
	"strings"
	"errors"
	"tool/order"
)

var ErrNoData = errors.New("[Excel] no data found")
var ErrParams = errors.New("[Excel] params error")

//create excel file by struct array
func CreateByStructs(data []interface{}) (err error, file *xlsx.File ) {
	if len(data) == 0 {
		return ErrNoData, nil
	}
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return
	}

	exports := []int{}
	t := reflect.TypeOf(data[0])
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {


		return ErrParams, nil // Should not panic here ?
	}
	//first line
	row := sheet.AddRow()
	for i := 0; i < t.NumField(); i++ {
		err, cellName := getCellNameByTag(t.Field(i).Tag.Get("excel"))
		if err != nil {
			return err, nil
		}
		if cellName == "-" {
			continue
		}
		exports = append(exports, i)
		row.AddCell().SetString(cellName)
	}

	//data
	for i := 0; i < len(data); i++ {
		val := reflect.ValueOf(data[i])
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		row = sheet.AddRow()
		for _, index := range exports {
			row.AddCell().SetValue(val.Field(index).Interface())
		}
	}

	return
}

//get the cell name from tag
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

func CreateByKVLink(head *order.KVLinkHead, data []map[string]interface{}) (file *xlsx.File, err error) {
	if head.GetCount() == 0 {
		return nil, ErrParams
	}
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
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

func CreateByKVMap(data []*order.KVMap) (file *xlsx.File, err error) {
	file = xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
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


