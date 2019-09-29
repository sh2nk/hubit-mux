package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<img src="/stream">`)
}

func stream(w http.ResponseWriter, r *http.Request) {
	const boundary = `frame`
	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
	multipartWriter := multipart.NewWriter(w)
	multipartWriter.SetBoundary(boundary)

	name := uuid.Must(uuid.NewV4()).String()

	func() {
		pool.Lock()
		defer pool.Unlock()
		pool.Streams[name] = send
	}()
	defer func() {
		pool.Lock()
		defer pool.Unlock()
		delete(pool.Streams, name)
	}()

	for buf := range send {
		image := buf.Bytes()
		iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-type":   []string{"image/jpeg"},
			"Content-length": []string{strconv.Itoa(len(image))},
		})
		if err != nil {
			log.Println(err)
		}
		_, err = iw.Write(image)
		if err != nil {
			log.Println(err)
		}
	}
}
