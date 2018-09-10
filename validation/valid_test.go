package validation

import (
	"testing"
	"fmt"
)

type TestValid struct {
	Test 	string 		 `valid:"Required;Regex(/^1[0-9]*9$/)"`
	Phone 	string 		 `valid:"Required;Mobile"`
	Email 	string 		 `valid:"Required;Email"`
	Ip 		string 		 `valid:"Required;Ip"`
	S 		string 		 `valid:"Required;AlphaNumeric"`
	N 		string 		 `valid:"Required;Numeric"`
	A 		string 		 `valid:"Required;Alpha"`
	Min 	int 		 `valid:"Required;Min(3)"`
	Max 	int 		 `valid:"Required;Max(-1)"`
	MinSize string 		 `valid:"Required;MinSize(2)"`
	Tgf 	[]TestValid  `valid:"Required"`
}

func (t TestValid) ValidMessage() map[string]map[string]string {
	return map[string]map[string]string{
		"Test": {
			"Required": "的沙发",
			"Regex": "格式错误",
		},
		"Phone": {
			"Required": "必须填写该字段",
			"Mobile": "必须是电话号码",
		},
	}
}

func TestValidation_Valid(t *testing.T) {
	valid := NewValidation()
	test := &TestValid{Test:"9", Phone:"1840844992", A:"5", N:"5", MinSize:"13"}
	b, err := valid.Valid(test)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if !b {
		for k, errs := range valid.ErrorsMap {
			fmt.Println(k + ": ")
			for _, v := range errs {
				fmt.Println("\t" + v.Msg)
			}
		}
	}
}

