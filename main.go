// https://rafallorenz.com/go/go-multidomain-host-switch/

package main

import (
	"fmt"
	"log"
	"net/http"
)

type HostSwitch map[string]http.Handler

func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if handler, ok := hs[r.Host]; ok && handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the first home page!")
	})

	muxTwo := http.NewServeMux()
	muxTwo.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the second home page!")
	})

	hs := make(HostSwitch)
	hs["example-one.local:8080"] = mux
	hs["example-two.local:8080"] = muxTwo

	log.Fatal(http.ListenAndServe(":8080", hs))
}
