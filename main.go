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
	w.Write([]byte("hello"))
	
	buf := make([]byte, len([]byte("hello")))
	n, err := r.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(buf[:n]))
}
