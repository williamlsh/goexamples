// Reference: https://github.com/GoogleCloudPlatform/golang-samples/tree/master/appengine/go11x/helloworld
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "hello world!")
}
