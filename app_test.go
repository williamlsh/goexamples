package main

import (
	"bytes"
	"encoding/json"
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

func TestBlogShow(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/blogs?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAccept, jsonapi.MediaType)

	w := httptest.NewRecorder()
	handler := &ExampleHandler{}

	fmt.Println("============ start show ===========")
	handler.ServeHTTP(w, req)
	fmt.Println("============ stop show ===========")

	jsonReply, _ := ioutil.ReadAll(w.Body)

	fmt.Println("============ jsonapi response from show ===========")
	fmt.Println(string(jsonReply))
	fmt.Println("============== end raw jsonapi from show =============")
}

func TestBlogCreate(t *testing.T) {
	blog := fixtureBlogCreate(1)
	in := new(bytes.Buffer)
	jsonapi.MarshalOnePayloadEmbedded(in, blog)

	req, err := http.NewRequest(http.MethodGet, "/blogs", in)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAccept, jsonapi.MediaType)

	w := httptest.NewRecorder()
	handler := &ExampleHandler{}

	fmt.Println("============ start create ===========")
	handler.ServeHTTP(w, req)
	fmt.Println("============ stop create ===========")

	jsonReply, _ := ioutil.ReadAll(w.Body)

	fmt.Println("============ jsonapi response from create ===========")
	fmt.Println(string(jsonReply))
	fmt.Println("============== end raw jsonapi from create =============")
}

func TestBlogEcho(t *testing.T) {
	blogs := []interface{}{
		fixtureBlogCreate(1),
		fixtureBlogCreate(2),
		fixtureBlogCreate(3),
	}
	in := new(bytes.Buffer)
	jsonapi.MarshalOnePayloadEmbedded(in, blogs)

	req, err := http.NewRequest(http.MethodGet, "/blogs", in)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(headerAccept, jsonapi.MediaType)

	w := httptest.NewRecorder()
	handler := &ExampleHandler{}

	fmt.Println("============ start create ===========")
	handler.ServeHTTP(w, req)
	fmt.Println("============ stop create ===========")

	jsonReply, _ := ioutil.ReadAll(w.Body)

	fmt.Println("============ jsonapi response from create ===========")
	fmt.Println(string(jsonReply))
	fmt.Println("============== end raw jsonapi from create =============")

	responseBlog := new(Blog)

	buf := new(bytes.Buffer)
	jsonapi.UnmarshalPayload(buf, responseBlog)

	out := bytes.NewBuffer(nil)
	json.NewEncoder(out).Encode(responseBlog)

	fmt.Println("================ Viola! Converted back our Blog struct =================")
	fmt.Println(string(out.Bytes()))
	fmt.Println("================ end marshal materialized Blog struct =================")
}
