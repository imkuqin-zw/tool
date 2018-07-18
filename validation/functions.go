package validation

import (
	"reflect"
	"strings"
	"regexp"
	"fmt"
	"strconv"
)

func isStruct(t reflect.Type) bool {
	return 	t.Kind() == reflect.Struct
}

func isStructPtr(t reflect.Type) bool {
	return 	t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func HasRequired(funcs []ValidFunc) bool {
	if len(funcs) == 0 {
		return false
	}
	for _, item := range funcs {
		if item.Name == "Required" {
			return true
		}
	}
	return false
}

func getFuncs(field reflect.StructField) (funcs []ValidFunc, err error) {
	tags, hasRegex, required := parseTag(field)
	if hasRegex {
		if required {
			funcs = make([]ValidFunc, 2)
			if funcs[0], err = parseFunc(tags[0], field.Name); err != nil {
				return
			}
			if funcs[1], err = parseRegFunc(tags[1], field.Name); err != nil {
				return
			}
		} else {
			funcs = make([]ValidFunc, 1)
			if funcs[0], err = parseRegFunc(tags[0], field.Name); err != nil {
				return
			}
		}
		return
	}
	funcs = make([]ValidFunc, 0, len(tags))
	for _, tag := range tags {
		var fs ValidFunc
		if fs, err = parseFunc(tag, field.Name); err != nil {
			return
		}
		funcs = append(funcs, fs)
	}
	return
}

func parseTag(field reflect.StructField) (result []string, hasRegex, required bool) {
	tag := field.Tag.Get(ValidTag)
	if tag == "" {
		return
	}
	tags := strings.Split(tag, ";")
	result = []string{}
	var regex string
	for _, item := range tags {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if item == "Required" {
			required = true
			if hasRegex {
				result = []string{"Required", regex}
				return
			}
		} else if isRegex(item) {
			hasRegex = true
			regex = item
			if required {
				result = []string{"Required", regex}
				return
			}
			continue
		}
		result = append(result, item)
	}
	if hasRegex {
		result = []string{regex}
	}
	return
}

func isRegex(in string) bool {
	index := strings.Index(in, "Regex(/")
	if index == -1 {
		return false
	}
	end := strings.Index(in, "/)")
	if end < index {
		return false
	}
	return true
}

func parseRegFunc(tag, field string) (validFunc ValidFunc, err error) {
	index := strings.Index(tag, "Regex(")
	if index == -1 {
		err = fmt.Errorf("%s : invalid Regex function", field)
		return
	}
	end := strings.LastIndex(tag, ")")
	if end < index {
		err = fmt.Errorf("%s : invalid Regex function", field)
		return
	}
	params := strings.Split(tag[index+len("Regex(") : end], ",")
	paramsNum := len(params)
	if paramsNum > 2 {
		err = fmt.Errorf("%s : Regex require 0 or 1 parameters", field)
		return
	}
	reg, err := regexp.Compile(strings.Trim(strings.TrimSpace(params[0]), "/"))
	if err != nil {
		return
	}
	funcParams := []interface{}{field, reg, true}
	if paramsNum == 2 {
		funcParams[2] = strings.TrimSpace(params[1])
	}
	validFunc = ValidFunc{Name: "Regex", Params: funcParams}
	return
}

func parseFunc(tag, field string) (validFunc ValidFunc, err error) {
	start := strings.Index(tag, "(")
	funcName := tag
	if start != -1 {
		funcName = tag[:start]
	}
	fn, b := validFuncs.GetFuncByName(funcName)
	if !b {
		err = fmt.Errorf("%s : not found the %s function", field, tag)
		return
	}

	paramsNum := validFuncs.GetParamsNumByName(funcName)
	if paramsNum == -1  {
		err = fmt.Errorf("%s : invalid %s function", field, tag)
		return
	}
	paramsNum -= 2
	params := []interface{}{field}
	if start == -1 {
		if paramsNum > 0 {
			err = fmt.Errorf("%s : %s require %d parameters", field, tag, paramsNum)
			return
		}
	} else {
		end := strings.LastIndex(tag, ")")
		if end < start {
			err = fmt.Errorf("%s : invalid Regex function", field)
			return
		}
		fParams := strings.Split(tag[start+1:end], ",")
		if paramsNum != len(params) {
			err = fmt.Errorf("%s : %s require %d parameters", field, tag, paramsNum)
			return
		}
		var temp interface{}
		for i, item := range fParams {
			temp, err = ParseParam(fn.Type().In(i + 3), strings.TrimSpace(item))
			params = append(params, temp)
		}
	}

	validFunc = ValidFunc{Name:funcName, Params:params}
	return
}

func ParseParam(t reflect.Type, param string) (i interface{}, err error) {
	switch t.Kind() {
	case reflect.Int:
		i, err = strconv.Atoi(param)
	case reflect.Int64:
		if wordsize == 32 {
			return nil, fmt.Errorf("not support int64 on 32-bit platform")
		}
		i, err = strconv.ParseInt(param, 10, 64)
	case reflect.Int32:
		var v int64
		v, err = strconv.ParseInt(param, 10, 32)
		if err == nil {
			i = int32(v)
		}
	case reflect.Int16:
		var v int64
		v, err = strconv.ParseInt(param, 10, 16)
		if err == nil {
			i = int16(v)
		}
	case reflect.Int8:
		var v int64
		v, err = strconv.ParseInt(param, 10, 8)
		if err == nil {
			i = int8(v)
		}
	case reflect.String:
		i = param
	default:
		err = fmt.Errorf("not support %s", t.Kind().String())
	}
	return
}

func mergeParam(v *Validation, obj interface{}, params []interface{}) []interface{} {
	return append([]interface{}{v, obj}, params...)
}

func IsNumType(field interface{}) bool {
	switch field.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float64:
		return true
	case float32:
		return true
	}
	return false
}

