package main

import (
	"bytes"
	"log"
	"net/http"

	"hubit-mux/config"
	"hubit-mux/streamer"
	"hubit-mux/view"
)

var (
	conf   *config.Config
	send   chan *bytes.Buffer
	camera *view.Camera

	//Пул клиентов, получающий поток
	pool = &streamer.Pool{
		Streams: make(map[string]chan *bytes.Buffer, 12),
	}
)

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
	go camera.Read(send, pool)

	log.Printf("Listening on %s...\n", conf.Addr)
	http.HandleFunc("/stream", stream)
	http.HandleFunc("/", index)
	http.ListenAndServe(conf.Addr, nil)
}
