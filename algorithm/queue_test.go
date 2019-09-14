package algorithm

import (
	"fmt"
	"testing"
)

var queue *Queue

func TestNewQueue(t *testing.T) {
	queue = NewQueue()
}

func TestQueue_Push(t *testing.T) {
	for i := 0; i < 10; i++ {
		queue.Push(i)
	}
}

func TestQueue_Pop(t *testing.T) {
	result := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		result[i] = queue.Pop()
	}
	fmt.Println(result)
}

func TestQueue_IsEmpty(t *testing.T) {
	fmt.Println(queue.IsEmpty())
}

func TestQueue_Length(t *testing.T) {
	fmt.Println(queue.Length())
}

func TestQueue(t *testing.T) {
	TestNewQueue(t)
	TestQueue_IsEmpty(t)
	TestQueue_Length(t)
	TestQueue_Push(t)
	TestQueue_IsEmpty(t)
	TestQueue_Length(t)
	TestQueue_Pop(t)
	TestQueue_IsEmpty(t)


	
}
