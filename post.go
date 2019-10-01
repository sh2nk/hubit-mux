package main

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/image/draw"

	uuid "github.com/satori/go.uuid"
)

func post(url string, width uint32, height uint32, level uint32) {
	stream := make(chan *bytes.Buffer)
	name := uuid.Must(uuid.NewV4()).String()

	func() {
		pool.Lock()
		defer pool.Unlock()
		pool.Streams[name] = stream
	}()
	defer func() {
		pool.Lock()
		defer pool.Unlock()
		delete(pool.Streams, name)
	}()

	for {
		frame := <-stream

		src, _, err := image.Decode(frame)
		if err != nil {
			log.Fatal(err)
		}
		dest := image.NewRGBA(image.Rect(0, 0, int(width/level), int(height/level)))
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
}
