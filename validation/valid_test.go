package validation

import (
	"testing"
	"fmt"
)

type TestValid struct {
	Test 	string 	`valid:"Required;Regex(/^1[0-9]*9$/)"`
	Phone 	string 	`valid:"Required;Mobile"`
}

func (t TestValid) ValidMessage() map[string]map[string]string {
	return map[string]map[string]string{
		"Test": {
			"Required": "必须填写该字段",
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
	test := &TestValid{Test:"19", Phone:"18408244992"}
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

