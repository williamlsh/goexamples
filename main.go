package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"os"
	"path/filepath"
)

func main() {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("dumped request:\n %q\n\noriginal request:\n %+v\n", dump, r)
		io.WriteString(w, "ok")
	}))
	defer ts.Close()
	
	file, err := os.Open("WechatIMG1.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
	}
	writer.Close()
	
	r, err := http.NewRequest(http.MethodPost, ts.URL, body)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	_, err = ts.Client().Do(r)
	if err != nil {
		log.Fatal(err)
	}
}
