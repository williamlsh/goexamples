package main

import (
	"fmt"
	"sync"
)

// Consider a pipeline with three stages.

// The first stage, gen, is a function that converts a list of integers to a channel that emits the integers in the list.
func gen(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int, len(nums)) // use buffered channel to avoid creating new goroutine.
	for _, n := range nums {
		select {
		case out <- n:
		case <-done:
			break
		}
	}
	close(out)
	return out
}

// The second stage, sq, receives integers from a channel and returns a channel that emits the square of each received integer.
func sq(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}

func main() {
	// Set up a done channel that's shared by the whole pipeline,
	// and close that channel when this pipeline exits, as a signal
	// for all the goroutines we started to exit.
	done := make(chan struct{})
	defer close(done)

	// Set up the pipeline and runs the final stage.
	in := gen(done, 2, 3)

	// Distribute the sq work across two goroutines that both read from in.
	c1 := sq(done, in)
	c2 := sq(done, in)

	// Consume the first value from output.
	out := merge(done, c1, c2)
	fmt.Println(<-out)
}

func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed or it receives a value
	// from done, then calls wg.Done.
	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done: // return as soon as done is closed without draining its inbound channel
				return
			}
		}
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
