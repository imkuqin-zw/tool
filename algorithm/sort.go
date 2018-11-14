package algorithm

type compareFunc func(before interface{}, after interface{}) bool

// 插入排序
func InsertSort(in []interface{}, compare compareFunc) {
	for i := 1; i < len(in); i++ {
		var j = i - 1
		for j >= 0 && !compare(in[i], in[j]) {
			j--
		}
		if j != i-1 {
			temp := in[i]
			for k := i; k > j+1; k-- {
				in[k] = in[k-1]
			}
			in[j+1] = temp
		}
	}
	return
}

func shellPass(in []interface{}, skip int, compare compareFunc) {
	length := len(in)
	for x := 0; x < skip; x++ {
		for i := x + skip; i < length; i += skip {
			var j = i - skip
			for j >= 0 && !compare(in[i], in[j]) {
				j -= skip
			}
			if j != i-skip {
				temp := in[i]
				for k := i; k > j+skip; k -= skip {
					in[k] = in[k-skip]
				}
				in[j+skip] = temp
			}
		}
	}
	return
}

//希尔排序
func ShellSort(in []interface{}, compare compareFunc) {
	skip := len(in)
	for skip > 1 {
		skip /= 2
		shellPass(in, skip, compare)
	}
	return
}

//冒泡排序
func BubbleSort(in []interface{}, compare compareFunc) {
	length := len(in)
	for i := 0; i < length-1; i++ {
		for j := 0; j < length-i-1; j++ {
			if compare(in[j], in[j+1]) {
				in[j], in[j+1] = in[j+1], in[j]
			}
		}
	}
	return
}

func partition(in []interface{}, left, right int, compare compareFunc) int {
	tmp := in[left]
	for left < right {
		for left < right && !compare(tmp, in[right]) {
			right--
		}
		if left < right {
			in[left] = in[right]
			left++
		}
		for left < right && !compare(in[left], tmp) {
			left++
		}
		if left < right {
			in[right] = in[left]
			right--
		}
	}
	in[left] = tmp
	return left
}

func quickSort(in []interface{}, left, right int, compare compareFunc) {
	if left < right {
		mid := partition(in, left, right, compare)
		quickSort(in, left, mid-1, compare)
		quickSort(in, mid+1, right, compare)
	}
	return
}

//快速排序
func QuickSort(in []interface{}, compare compareFunc) {
	quickSort(in, 0, len(in)-1, compare)
}

//选择排序
func SelectSort(in []interface{}, compare compareFunc) {
	length := len(in)
	for i := 0; i < length; i++ {
		k := i
		for j := i + 1; j < length; j++ {
			if compare(in[k], in[j]) {
				k = j
			}
		}
		if k != i {
			in[k], in[i] = in[i], in[k]
		}
	}
	return
}

func merge(in []interface{}, left, mid, right int, compare compareFunc) {
	temp := make([]interface{}, 0, right-left+1)
	var l, r = left, mid + 1
	for l <= mid && r <= right {
		if !compare(in[l], in[r]) {
			temp = append(temp, in[l])
			l++
		} else {
			temp = append(temp, in[r])
			r++
		}
	}
	if l <= mid {
		temp = append(temp, in[l:mid+1]...)
	} else {
		temp = append(temp, in[r:right+1]...)
	}
	in = append(append(in[0:left], temp...), in[right+1:]...)

}

func mergeSort(in []interface{}, left, right int, compare compareFunc) {
	if left < right {
		mid := (right + left) / 2
		mergeSort(in, left, mid, compare)
		mergeSort(in, mid+1, right, compare)
		merge(in, left, mid, right, compare)
	}
}

// 归并排序
func MergeSort(in []interface{}, compare compareFunc) {
	mergeSort(in, 0, len(in)-1, compare)
}
