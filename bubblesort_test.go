package goexamples

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	data := []int{2, 1, 3, 4, 5}
	fmt.Printf("Before sorting: %+v\n", data)
	bubbleSort(data)
	fmt.Printf("After sorting: %+v\n", data)
	if reflect.DeepEqual(data, []int{1, 2, 3, 4, 5}) {
		fmt.Println("Success")
	} else {
		t.Fail()
	}
}
