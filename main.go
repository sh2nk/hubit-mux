package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sh2nk/hubit-mux/utils"
	"github.com/sh2nk/hubit-mux/view"
)

var (
	camera *view.Camera

	//Пул клиентов, получающий поток
	streamPool = &utils.StreamPool{
		Streams: make(map[string]chan *bytes.Buffer, 12),
	}

	wsPool = &utils.WSPool{
		Clients: make(map[string]*websocket.Conn),
	}
)

func main() {
	err := utils.Config.Parse("config.json")
	if err != nil {
		log.Fatalf("Error parsing config: %+v", err)
	}

	fmt.Println(utils.Config)

	camera, err := view.NewCamera(utils.Config.Device, utils.Config.WB)
	if err != nil {
		log.Fatalf("Camera init error: %+v", err)
	}
	defer camera.Close()

	if err = camera.Setup(utils.Config.Format, utils.Config.Width, utils.Config.Height); err != nil {
		log.Fatalf("Camera init error: %+v", err)
	}

	go camera.ReadAndStream(streamPool)

	log.Printf("Listening on %s...\n", utils.Config.Addr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/view", viewPage)
	http.HandleFunc("/about", aboutPage)
	http.HandleFunc("/settings", settingsPage)
	http.HandleFunc("/stream", stream)
	http.HandleFunc("/getface", getFaceData)
	//http.HandleFunc("/ws", wsServer)

	if err = http.ListenAndServe(utils.Config.Addr, nil); err != nil {
		log.Fatalf("Listen error: %+v", err)
	}
}
