package main

import (
	"bytes"
	"fmt"
	"hubit-mux/config"
	"hubit-mux/view"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
)

var (
	conf   *config.Config
	send   chan *bytes.Buffer
	camera *view.Camera
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<img src="/stream">`)
}

func stream(w http.ResponseWriter, r *http.Request) {
	<-send
	const boundary = `frame`
	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
	multipartWriter := multipart.NewWriter(w)
	multipartWriter.SetBoundary(boundary)
	for {
		buf := <-send
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

func main() {
	conf = new(config.Config)
	err := conf.Parse("config.json")
	if err != nil {
		log.Fatal("Error parsing config ", err)
	}

	camera, err := view.NewCamera(conf.Device, conf.WB)
	if err != nil {
		log.Fatal("Camera init error ", err)
	}
	defer camera.Close()

	if err = camera.Setup(conf.Format, conf.Width, conf.Height); err != nil {
		log.Fatal(err)
	}

	send = make(chan *bytes.Buffer)
	go camera.Read(send)

	log.Printf("Listening on %s...\n", conf.Addr)
	http.HandleFunc("/stream", stream)
	http.HandleFunc("/", index)
	http.ListenAndServe(conf.Addr, nil)
}
