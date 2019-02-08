package main

import (
	"fmt"
	"sync"
)

// Consider a pipeline with three stages.

// The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list.
func gen(nums ...int) <-chan int {
	out := make(chan int, len(nums)) // use buffered channel to avoid creating new goroutine.
	for _, n := range nums {
		out <- n
	}
	close(out)
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
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(in)
	c2 := sq(in)

	// Consume the first value from output.
	done := make(chan struct{}, 2)
	out := merge(done, c1, c2)
	fmt.Println(<-out)

	// Tell the remianing senders we are leaving.
	done <- struct{}{}
	done <- struct{}{}
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed or it receives a value
	// from done, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			select {
			case out <- n:
			case <-done: // stop sending to out channel in case of its blocking with non-receiving when main exists
			}
		}
		wg.Done()
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
