package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	r, w, err := os.Pipe()
	if err != nil {
		log.Fatal(err)
	}
	origStdout := os.Stdout
	os.Stdout = w
	
	// This exceeds sys pipe buffers limit, see 'man pipe'
	for i := 0; i < 5000; i++ {
		fmt.Print("Hello to stdout")
	}
	
	buf := make([]byte, len([]byte("Hello to stdout")))
	n, err := r.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	
	// Restore original stdout.
	os.Stdout = origStdout
	fmt.Println("Written to stdout: ", string(buf[:n]))
}
