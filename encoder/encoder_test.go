package encoder

import (
	"fmt"
	"testing"
)

func Test_Base64EncodeString(t *testing.T) {
	str := Base64EncodeString("fdsfsdf佛挡杀佛啥大发沙发f我发动机暗示开房间附近的啥卡==")
	fmt.Println(str)
	str = Base64DecodeString("ZmRzZnNkZuS9m+aMoeadgOS9m+WVpeWkp+WPkeaymeWPkWbmiJHlj5HliqjmnLrmmpfnpLrlvIDmiL/pl7TpmYTov5HnmoTllaXljaE9PQ")
	fmt.Println(str)
}
