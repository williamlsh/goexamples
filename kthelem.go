package q

import "sort"

// KthElement merges two sorted integer slices and returns the kth elements.
// Note: k begins with 1.
func KthElement(a []int, b []int, k int) int {
	a = append(a, b...)
	sort.Ints(a)

	if k < 1 || k > len(a) {
		return -1
	}
	return a[k-1]
}
