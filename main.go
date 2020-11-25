package main

import (
	"io"
	"log"
	"net/http"

	"github.com/coreos/go-systemd/v22/activation"
)

func main() {
	listeners, err := activation.Listeners()
	if err != nil {
		log.Fatal(err)
	}
	if len(listeners) != 1 {
		panic("Unexpected number of socket activation fds")
	}

	http.HandleFunc("/", HelloServer)
	http.Serve(listeners[0], nil)
}

// HelloServer is an hello handler.
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello socket activated world!\n")
}
