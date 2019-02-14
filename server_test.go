package main

import (
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	req, _ := http.NewRequest("GET", "localhost:8080/multipart", nil)
	req.Header.Set("Accept", "multipart/form-data; charset=utf-8")
	resp, _ := http.DefaultClient.Do(req)
	_, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	mr := multipart.NewReader(resp.Body, params["boundary"])
	for part, err := mr.NextPart(); err == nil; part, err = mr.NextPart() {
		value, _ := ioutil.ReadAll(part)
		log.Printf("Value: %s", value)
	}
}
