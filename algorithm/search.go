package algorithm

type searchCompare func(a, b interface{}) int8

func BinarySearch(in []interface{}, val interface{}, compare searchCompare) int {
	length := len(in)
	l, r := 0, length-1
	for l <= r {
		mid := (r + l) >> 1
		tmp := compare(val, in[mid])
		if tmp == 0 {
			return mid
		} else if tmp == 1 {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return -1
}
