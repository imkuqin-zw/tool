package ali_sms

import (
	"testing"
	"fmt"
)

func TestSendVerifySms(t *testing.T) {
	err := SendVerifySms("18408249924", "456512")
	if err != nil {
		fmt.Println(err.Error())
	}

}
