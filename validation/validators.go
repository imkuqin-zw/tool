package validation

import (
	"time"
	"reflect"
	"regexp"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	wordsize = 32 << (^uint(0) >> 32 & 1)
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

type Min struct {
	Lang string
	Field string
	Min int
}

func (m Min) IsValid(field interface{}) (b bool) {
	var v int
	switch field.(type) {
	case int64:
		if wordsize == 32 {
			return false
		}
		v = int(field.(int64))
	case int:
		v = field.(int)
	case int32:
		v = int(field.(int32))
	case int16:
		v = int(field.(int16))
	case int8:
		v = int(field.(int8))
	default:
		return false
	}

	return v >= m.Min
}

func (m Min) GetMsg() string {
	if m.Lang == "" {
		m.Lang = DefMsgLang
	}
	 return fmt.Sprintf(MsgTmplMap[m.Lang]["Min"], m.Min)
}

func (m Min) GetField() string {
	return m.Field
}

type Max struct {
	Lang string
	Field string
	Max int
}

func (m Max) IsValid(field interface{}) (b bool) {
	var v int
	switch field.(type) {
	case int64:
		if wordsize == 32 {
			return false
		}
		v = int(field.(int64))
	case int:
		v = field.(int)
	case int32:
		v = int(field.(int32))
	case int16:
		v = int(field.(int16))
	case int8:
		v = int(field.(int8))
	default:
		return false
	}

	return v <= m.Max
}

func (m Max) GetMsg() string {
	if m.Lang == "" {
		m.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[m.Lang]["Max"], m.Max)
}

func (m Max) GetField() string {
	return m.Field
}

type Range struct {
	Lang string
	Field string
	Max int
	Min int
}

func (m Range) IsValid(field interface{}) (b bool) {
	var v int
	switch field.(type) {
	case int64:
		if wordsize == 32 {
			return false
		}
		v = int(field.(int64))
	case int:
		v = field.(int)
	case int32:
		v = int(field.(int32))
	case int16:
		v = int(field.(int16))
	case int8:
		v = int(field.(int8))
	default:
		return false
	}

	return v <= m.Max && v >= m.Min
}

func (m Range) GetMsg() string {
	if m.Lang == "" {
		m.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[m.Lang]["Range"], m.Min, m.Max)
}

func (m Range) GetField() string {
	return m.Field
}

type MinSize struct {
	Lang string
	Field string
	Min int
}

func (m MinSize) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return utf8.RuneCountInString(str) >= m.Min
	}
	v := reflect.ValueOf(field)
	if v.Kind() == reflect.Slice {
		return v.Len() >= m.Min
	}
	return false
}

func (m MinSize) GetMsg() string {
	if m.Lang == "" {
		m.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[m.Lang]["MinSize"], m.Min)
}

func (m MinSize) GetField() string {
	return m.Field
}

type MaxSize struct {
	Lang string
	Field string
	Max int
}

func (m MaxSize) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return utf8.RuneCountInString(str) <= m.Max
	}
	v := reflect.ValueOf(field)
	if v.Kind() == reflect.Slice {
		return v.Len() <= m.Max
	}
	return false
}

func (m MaxSize) GetMsg() string {
	if m.Lang == "" {
		m.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[m.Lang]["MaxSize"], m.Max)
}

func (m MaxSize) GetField() string {
	return m.Field
}

type Length struct {
	Lang string
	Field string
	Num int
}

func (l Length) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return utf8.RuneCountInString(str) == l.Num
	}
	v := reflect.ValueOf(field)
	if v.Kind() == reflect.Slice {
		return v.Len() == l.Num
	}
	return false
}

func (l Length) GetMsg() string {
	if l.Lang == "" {
		l.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[l.Lang]["Length"], l.Num)
}

func (l Length) GetField() string {
	return l.Field
}

type Alpha struct {
	Lang string
	Field string
}

func (a Alpha) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		reg := regexp.MustCompile(`^[a-zA-Z]+$`)
		return reg.MatchString(str)
	}
	return false
}

func (a Alpha) GetMsg() string {
	if a.Lang == "" {
		a.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[a.Lang]["Alpha"])
}

func (a Alpha) GetField() string {
	return a.Field
}

type Numeric struct {
	Lang string
	Field string
}

func (n Numeric) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		reg := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
		return reg.MatchString(str)
	}
	return IsNumType(field)
}

func (n Numeric) GetMsg() string {
	if n.Lang == "" {
		n.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[n.Lang]["Numeric"])
}

func (n Numeric) GetField() string {
	return n.Field
}

type AlphaNumeric struct {
	Lang string
	Field string
}

func (n AlphaNumeric) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		reg := regexp.MustCompile(`^[0-9a-zA-Z]+$`)
		return reg.MatchString(str)
	}
	return IsNumType(field)
}

func (n AlphaNumeric) GetMsg() string {
	if n.Lang == "" {
		n.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[n.Lang]["AlphaNumeric"])
}

func (n AlphaNumeric) GetField() string {
	return n.Field
}

type Email struct {
	Lang string
	Field string
}

var emailPattern = regexp.MustCompile(`^[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[\w!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[a-zA-Z0-9](?:[\w-]*[\w])?$`)

func (e Email) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return emailPattern.MatchString(str)
	}
	return false
}

func (e Email) GetMsg() string {
	if e.Lang == "" {
		e.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[e.Lang]["Email"])
}

func (e Email) GetField() string {
	return e.Field
}

var ipPattern = regexp.MustCompile(`^(([1]\d\d?|2[0-4]\d|25[0-5])\.){3}([1]\d\d?|2[0-4]\d|25[0-5])`)

type Ip struct {
	Lang string
	Field string
}

func (e Ip) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return ipPattern.MatchString(str)
	}
	return false
}

func (e Ip) GetMsg() string {
	if e.Lang == "" {
		e.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[e.Lang]["Ip"])
}

func (e Ip) GetField() string {
	return e.Field
}

var mobilePattern = regexp.MustCompile(`^((\+86)|(86))?(1(([35][0-9])|[8][0-9]|[7][06789]|[4][579]))\d{8}$`)

type Mobile struct {
	Lang string
	Field string
}

func (m Mobile) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return mobilePattern.MatchString(str)
	}
	return false
}

func (m Mobile) GetMsg() string {
	if m.Lang == "" {
		m.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[m.Lang]["Mobile"])
}

func (m Mobile) GetField() string {
	return m.Field
}

var telPattern = regexp.MustCompile(`^(0\d{2,3}(\-)?)?\d{7,8}$`)

type Tel struct {
	Lang string
	Field string
}

func (t Tel) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return telPattern.MatchString(str)
	}
	return false
}

func (t Tel) GetMsg() string {
	if t.Lang == "" {
		t.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[t.Lang]["Tel"])
}

func (t Tel) GetField() string {
	return t.Field
}

type Phone struct {
	Lang string
	Field string
}

func (p Phone) IsValid(field interface{}) (b bool) {
	if str, ok := field.(string); ok {
		return telPattern.MatchString(str) || mobilePattern.MatchString(str)
	}
	return false
}

func (p Phone) GetMsg() string {
	if p.Lang == "" {
		p.Lang = DefMsgLang
	}
	return fmt.Sprintf(MsgTmplMap[p.Lang]["Phone"])
}

func (p Phone) GetField() string {
	return p.Field
}


