package view

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/image/draw"
)

func post(url string, frame io.Reader, w int, h int, level int) {
	src, _, err := image.Decode(frame)
	if err != nil {
		log.Fatal(err)
	}
	dest := image.NewRGBA(image.Rect(0, 0, int(w/level), int(h/level)))
	draw.ApproxBiLinear.Scale(dest, dest.Bounds(), src, src.Bounds(), draw.Over, nil)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, dest, nil); err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/octet-stream", buf)
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Falied reading resp ", err)
	}
	log.Println(string(bytes))
}