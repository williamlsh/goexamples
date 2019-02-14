package main

import (
	"log"
	"mime"
	"mime/multipart"
	"net/http"
)

func multipartHandler(w http.ResponseWriter, r *http.Request) {
	mediatype, _, err := mime.ParseMediaType(r.Header.Get("Accept"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	if mediatype != "multipart/form-data" {
		http.Error(w, "set Accept: multipart/form-data", http.StatusMultipleChoices)
		return
	}
	mw := multipart.NewWriter(w)
	w.Header().Set("Content-Type", mw.FormDataContentType())
	for _, value := range getValues() {
		fw, err := mw.CreateFormField("value")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := fw.Write(value); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if err := mw.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/multipart", multipartHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
