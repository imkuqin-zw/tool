package validation

import (
	"reflect"
	"fmt"
	"github.com/imkuqin-zw/tool/reflectfunc"
	"regexp"
)

// Default error message language.
var DefMsgLang = "en"

// Validation functions.
var validFuncs reflectfunc.ReflectFunc

// ValidTag struct tag.
const ValidTag = "valid"

type Validation struct {
	Lang 		string
	Errors 		[]*Error
	ErrorsMap 	map[string][]*Error
	MsgMap 		map[string]map[string]string
}

func init() {
	ignoreMethod := []string{
		"apply",
		"SetError",
		"HasErrors",
		"Valid",
	}
	validFuncs = reflectfunc.NewReflectFunc(&Validation{}, ignoreMethod...)
}


func (v *Validation) Required(filed interface{}, name string) *Error {
	return v.apply(Required{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Regex(filed interface{}, name string, regex *regexp.Regexp, match bool) *Error {
	return v.apply(Regex{Field:name, Lang: v.Lang, Regex: regex, Match: match}, filed)
}

func (v *Validation) Min(filed interface{}, name string, min int) *Error {
	return v.apply(Min{Field:name, Lang: v.Lang, Min: min}, filed)
}

func (v *Validation) Max(filed interface{}, name string, max int) *Error {
	return v.apply(Max{Field:name, Lang: v.Lang, Max: max}, filed)
}

func (v *Validation) Range(filed interface{}, name string, min, max int) *Error {
	return v.apply(Range{Field:name, Lang: v.Lang, Max: max, Min: min}, filed)
}

func (v *Validation) MinSize(filed interface{}, name string, min int) *Error {
	return v.apply(MinSize{Field:name, Lang: v.Lang, Min: min}, filed)
}

func (v *Validation) MaxSize(filed interface{}, name string, max int) *Error {
	return v.apply(MaxSize{Field:name, Lang: v.Lang, Max: max}, filed)
}

func (v *Validation) Length(filed interface{}, name string, num int) *Error {
	return v.apply(Length{Field:name, Lang: v.Lang, Num: num}, filed)
}

func (v *Validation) Alpha(filed interface{}, name string) *Error {
	return v.apply(Alpha{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Numeric(filed interface{}, name string) *Error {
	return v.apply(Numeric{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) AlphaNumeric(filed interface{}, name string) *Error {
	return v.apply(AlphaNumeric{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Email(filed interface{}, name string) *Error {
	return v.apply(Email{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Ip(filed interface{}, name string) *Error {
	return v.apply(Ip{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Mobile(filed interface{}, name string) *Error {
	return v.apply(Mobile{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Tel(filed interface{}, name string) *Error {
	return v.apply(Tel{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) Phone(filed interface{}, name string) *Error {
	return v.apply(Phone{Field:name, Lang: v.Lang}, filed)
}

func (v *Validation) apply(validator Validator, filed interface{}) *Error {
	if validator.IsValid(filed) {
		return nil
	}
	err := &Error{
		Msg: validator.GetMsg(),
		Field: validator.GetField(),
	}
	validName := reflect.TypeOf(validator).Name()
	if msgMap, ok := v.MsgMap[validator.GetField()]; ok {
		if msg, ok := msgMap[validName]; ok {
			err.Msg = msg
		}
	}

	v.SetError(err)
	return err
}

func (v *Validation) SetError(err *Error) {
	v.Errors = append(v.Errors, err)
	if v.ErrorsMap == nil {
		v.ErrorsMap = make(map[string][]*Error)
	}
	if _, ok := v.ErrorsMap[err.Field]; !ok {
		v.ErrorsMap[err.Field] = []*Error{}
	}
	v.ErrorsMap[err.Field] = append(v.ErrorsMap[err.Field], err)
	return
}

func (v *Validation) HasErrors() bool {
	return len(v.Errors) > 0
}

func (v *Validation) Valid(obj interface{}) (b bool, err error) {
	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)

	// verify obj is valid.
	if !isStruct(objT) {
		if !isStructPtr(objT) {
			err = fmt.Errorf("%v must be a struct or a struct pointer", obj)
			return
		}
		objT = objT.Elem()
		objV = objV.Elem()
	}

	// Get custom error messages.
	if msg := objV.MethodByName("ValidMessage"); msg.IsValid() {
		v.MsgMap =  msg.Call(nil)[0].Interface().(map[string]map[string]string)
	}

	for i := 0; i < objT.NumField(); i++ {
		var funcs []ValidFunc
		if funcs, err = getFuncs(objT.Field(i)); err != nil {
			return
		}
		if !HasRequired(funcs) {
			if v.Required(objV.Field(i), objT.Field(i).Name) != nil {
				continue
			}
		}

		for _, item := range funcs {
			if _, err = validFuncs.Invoke(item.Name,
				mergeParam(v, objV.Field(i).Interface(), item.Params)...); err != nil {
				return
			}
		}
	}

	return !v.HasErrors(), nil
}



