package snowflake

import (
	"testing"
	"fmt"
)

func TestGetUUID(t *testing.T) {
	fmt.Println(maxCenterId)
}

func TestGetNewToken(t *testing.T) {
	Init(1,2)
	fmt.Println(GetNewToken())
}