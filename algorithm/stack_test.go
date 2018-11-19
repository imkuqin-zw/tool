package algorithm

import (
	"fmt"
	"testing"
)

var stack *Stack

func TestNewStack(t *testing.T) {
	stack = NewStack()
}

func TestStack_Push(t *testing.T) {
	for i := 0; i < 10; i++ {
		stack.Push(i)
	}
}

func TestStack_Pop(t *testing.T) {
	result := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		result[i] = stack.Pop()
	}
	fmt.Println(result)
}

func TestStack_IsEmpty(t *testing.T) {
	fmt.Println(stack.IsEmpty())
}

func TestStack_Length(t *testing.T) {
	fmt.Println(stack.Length())
}

func TestStack(t *testing.T) {
	TestNewStack(t)
	TestStack_IsEmpty(t)
	TestStack_Length(t)
	TestStack_Push(t)
	TestStack_IsEmpty(t)
	TestStack_Length(t)
	TestStack_Pop(t)
	TestStack_IsEmpty(t)
}
