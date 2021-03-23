package q

import "testing"

func TestKthElement(t *testing.T) {
	a := []int{5, 7, 9}
	b := []int{6, 7, 8}
	k := 5

	if ele := KthElement(a, b, k); ele != 8 {
		t.Fatalf("got %d want %d", ele, 8)
	}
}
