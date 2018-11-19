package algorithm

import (
	"fmt"
	"testing"
)

func intCompare(before interface{}, after interface{}) int8 {
	if before.(int) > after.(int) {
		return 1
	} else if before.(int) == after.(int) {
		return 0
	} else {
		return -1
	}
}

func TestBinarySearch(t *testing.T) {
	in := []interface{}{1, 2, 3, 4, 5, 9}
	fmt.Println(BinarySearch(in, 9, intCompare))
}
