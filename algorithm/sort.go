package algorithm

type compareFunc func(before interface{}, after interface{}) bool

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

func ShellSort(in []interface{}, compare compareFunc) {
	skip := len(in)
	for skip > 1 {
		skip /= 2
		shellPass(in, skip, compare)
	}
	return
}

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

func QuickSort(in []interface{}, left, right int, compare compareFunc) {
	if left < right {
		mid := partition(in, left, right, compare)
		QuickSort(in, left, mid-1, compare)
		QuickSort(in, mid+1, right, compare)
	}
	return
}

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
