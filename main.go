package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func main() {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request.")
		
		// reqBody := new(bytes.Buffer)
		// r.Body = ioutil.NopCloser(io.TeeReader(r.Body, reqBody))
		// var bodyBytes []byte
		// bodyBytes, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	log.Fatalf("Could not read request body: %v", err)
		// }
		// r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// fmt.Printf("Copied request body: %s\n\n\n", bodyBytes)
		
		// dump, err := httputil.DumpRequest(r, true)
		// if err != nil {
		// 	log.Fatalf("Could not dump request: %v", err)
		// }
		// fmt.Printf("Dumped request: \n%q\n\n\n\n", dump)
		// fmt.Printf("Original request: \n%+v\n\n\n", r)
		
		// reqBuf := bytes.NewBuffer(dump)
		// copiedReq, err := http.ReadRequest(bufio.NewReader(reqBuf))
		// if err != nil {
		// 	log.Fatalf("Could not read request: %v", err)
		// }
		err := r.ParseForm()
		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Fatalf("Could not parse multipart form of request: %v", err)
		}
		fmt.Printf("Parsed request multipart form: %#v\n", r.MultipartForm)
		
		io.WriteString(w, "OK.")
	}))
	defer ts.Close()
	
	file, err := os.Open("WechatIMG1.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()
	
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal(err)
	}
	
	r, err := http.NewRequest(http.MethodPost, ts.URL, body)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := ts.Client().Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Body: %s", b)
}
