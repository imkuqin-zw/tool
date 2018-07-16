package validation

import (
	"time"
	"reflect"
	"regexp"
	"fmt"
	"strings"
)

// validator type
type ValidFunc struct {
	Name string
	Params []interface{}
}

// Validator interface
type Validator interface {
	IsValid(interface{}) bool
	GetMsg() string
	GetField() string
	//GetLimitValue() interface{}
}

// Required struct.
type Required struct {
	Lang string
	Field string
}

func (r Required) IsValid(field interface{}) bool {
	if field == nil {
		return false
	}

	if str, ok := field.(string); ok {
		return len(strings.TrimSpace(str)) > 0
	}
	if _, ok := field.(bool); ok {
		return true
	}
	if i, ok := field.(int); ok {
		return i != 0
	}
	if i, ok := field.(uint); ok {
		return i != 0
	}
	if i, ok := field.(int8); ok {
		return i != 0
	}
	if i, ok := field.(uint8); ok {
		return i != 0
	}
	if i, ok := field.(int16); ok {
		return i != 0
	}
	if i, ok := field.(uint16); ok {
		return i != 0
	}
	if i, ok := field.(uint32); ok {
		return i != 0
	}
	if i, ok := field.(int32); ok {
		return i != 0
	}
	if i, ok := field.(int64); ok {
		return i != 0
	}
	if i, ok := field.(uint64); ok {
		return i != 0
	}
	if t, ok := field.(time.Time); ok {
		return !t.IsZero()
	}


	v := reflect.ValueOf(field)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

func (r Required) GetMsg() string {
	if r.Lang == "" {
		r.Lang = DefMsgLang
	}
	return MsgTmplMap[r.Lang]["Required"]
}

func (r Required) GetField() string {
	return r.Field
}

// Required struct.
type Regex struct {
	Match bool
	Lang string
	Field string
	Regex *regexp.Regexp
}

func (r Regex) IsValid(field interface{}) (b bool) {
	var param = fmt.Sprintf("%v", field)
	if r.Match {
		b = r.Regex.MatchString(param)
	} else {
		b = !r.Regex.MatchString(param)
	}
	return
}

func (r Regex) GetMsg() (msg string) {
	if r.Lang == "" {
		r.Lang = DefMsgLang
	}
	if r.Match {
		msg = fmt.Sprintf(MsgTmplMap[r.Lang]["Match"], r.Regex.String())
	} else {
		msg = fmt.Sprintf(MsgTmplMap[r.Lang]["NoMatch"], r.Regex.String())
	}
	return
}

func (r Regex) GetField() string {
	return r.Field
}

