// Reference: https://www.sanarias.com/blog/115LearningHTTPcachinginGo
package main

import (
	"bytes"
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
)

var root = flag.String("root", ".", "File system path")

func main() {
	flag.Parse()

	http.HandleFunc("/black/", handleBlack)
	http.Handle("/", http.FileServer(http.Dir(*root)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleBlack(w http.ResponseWriter, r *http.Request) {
	m := image.NewRGBA(image.Rect(0, 0, 240, 240))
	black := color.RGBA{0, 0, 0, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{black}, image.Point{}, draw.Src)

	var img image.Image = m
	if err := writeImage(w, &img); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeImage(w http.ResponseWriter, img *image.Image) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, *img, nil); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
