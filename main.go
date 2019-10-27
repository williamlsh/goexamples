package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
)

const url = "https://example.com"

func main() {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	trace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			fmt.Printf("hostPort: %s\n", hostPort)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got conn: %+v\n", connInfo)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("Got dns info: %+v\n", dnsInfo)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}
}
