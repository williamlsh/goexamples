// Reference:
// - https://www.sanarias.com/blog/115LearningHTTPcachinginGo
// - https://devcenter.heroku.com/articles/increasing-application-performance-with-http-cache-headers
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/black/", handleBlack)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBlack(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)

	if matched := r.Header.Get("If-None-Match"); matched != "" {
		if strings.Contains(matched, etag()) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	if err := jpeg.Encode(w, image.Image(m), nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Etag", etag())
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
}

func etag() string {
	return "someKey"
}
