package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
)

func TestBlogList(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/blogs", nil)

	req.Header.Set(headerAccept, jsonapi.MediaType)

	w := httptest.NewRecorder()

	fmt.Println("============ start list ===========")
	mux := http.NewServeMux()
	exampleHandler := &ExampleHandler{}
	mux.Handle("/blogs", exampleHandler)
	mux.ServeHTTP(w, req)
	fmt.Println("============ stop list ===========")

	jsonReply, _ := ioutil.ReadAll(w.Body)

	fmt.Println("============ jsonapi response from list ===========")
	fmt.Println(string(jsonReply))
	fmt.Println("============== end raw jsonapi from list =============")
}
