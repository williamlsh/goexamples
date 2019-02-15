package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"
)

func TestBlogsList(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/blogs", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAccept, jsonapi.MediaType)

	w := httptest.NewRecorder()
	handler := &ExampleHandler{}

	fmt.Println("============ start list ===========")
	handler.ServeHTTP(w, req)
	fmt.Println("============ stop list ===========")

	jsonReply, _ := ioutil.ReadAll(w.Body)

	fmt.Println("============ jsonapi response from list ===========")
	fmt.Println(string(jsonReply))
	fmt.Println("============== end raw jsonapi from list =============")
}
