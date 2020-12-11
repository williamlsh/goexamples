// Reference: https://www.sanarias.com/blog/115LearningHTTPcachinginGo
package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/black/", handleBlack)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBlack(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	w.Header().Set("Content-Type", "image/jpeg")

	if err := jpeg.Encode(w, image.Image(m), nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
