package algorithm

import (
	"fmt"
	"testing"
)

func compare(before interface{}, after interface{}) bool {
	if before.(int) > after.(int) {
		return true
	}
	return false
}

func TestInsertSort(t *testing.T) {
	in := []interface{}{5, 9, 1, 1, 2, 3, 4}
	InsertSort(in, compare)
	fmt.Println(in)
}

func TestShellSort(t *testing.T) {
	in := []interface{}{5, 2, 75, 5, 6, 9, 7, 1, 0, 2, 5, 0, 9, 1, 3, 4}
	ShellSort(in, compare)
	fmt.Println(in)
}

func TestBubbleSort(t *testing.T) {
	in := []interface{}{5, 2, 75, 5, 6, 9, 7, 1, 0, 2, 5, 0, 9, 1, 3, 4}
	BubbleSort(in, compare)
	fmt.Println(in)
}

func TestQuickSort(t *testing.T) {
	in := []interface{}{5, 2, 75, 5, 6, 9, 7, 1, 0, 2, 5, 0, 9, 1, 3, 4}
	QuickSort(in, 0, len(in)-1, compare)
	fmt.Println(in)
}

func TestSelectSort(t *testing.T) {
	in := []interface{}{5, 2, 75, 5, 6, 9, 7, 1, 0, 2, 5, 0, 9, 1, 3, 4}
	SelectSort(in, compare)
	fmt.Println(in)
}
