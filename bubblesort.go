package goexamples

import "fmt"

func bubbleSort(data []int) {
	for step := 0; step < len(data)-1; step++ {
		fmt.Printf("Step %d\n", step)
		var swapped bool
		for i := 0; i < len(data)-step-1; i++ {
			if data[i] > data[i+1] {
				temp := data[i]
				data[i] = data[i+1]
				data[i+1] = temp
				swapped = true
				fmt.Println("Numbers swapped.")
			}
		}
		if !swapped {
			fmt.Println("Numbers already swapped.")
		}
	}
}
