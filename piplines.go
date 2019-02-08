package main

import "fmt"

// Consider a pipeline with three stages.

// The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list.
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// The second stage, sq, receives integers from a channel and returns a channel that emits the square of each received integer.
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func main() {
	// Set up the pipeline and runs the final stage.
	c := gen(2, 3)
	out := sq(c)

	// Consume the output.
	fmt.Println(<-out)
	fmt.Println(<-out)


	// Set up the pipline and consume the output.
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n)
	}
}
