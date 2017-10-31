package Excel

import (
	"time"
	"reflect"
	"fmt"
)

type ExcelFunc struct {}

func (ef *ExcelFunc) TimeFormat(sec interface{}, format string) (result string, err error) {
	switch reflect.ValueOf(sec).Kind() {
	case reflect.Int64:
		result = time.Unix(sec.(int64), 0).Format(format)
		break
	case reflect.Int:
		result = time.Unix(int64(sec.(int)), 0).Format(format)
		break
	case reflect.Int32:
		result = time.Unix(int64(sec.(int32)), 0).Format(format)
		break
	default:
		err = fmt.Errorf("[Excel] %s", "参数类型错误")
	}

	return
}

