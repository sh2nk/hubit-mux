package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"hubit-mux/utils"
	"hubit-mux/view"
)

var (
	send   chan *bytes.Buffer
	camera *view.Camera

	//Пул клиентов, получающий поток
	pool = &utils.Pool{
		Streams: make(map[string]chan *bytes.Buffer, 12),
	}
)

func main() {
	err := utils.Config.Parse("config.json")
	if err != nil {
		log.Fatal("Error parsing config ", err)
	}

	fmt.Println(utils.Config)

	camera, err := view.NewCamera(utils.Config.Device, utils.Config.WB)
	if err != nil {
		log.Fatal("Camera init error ", err)
	}
	defer camera.Close()

	if err = camera.Setup(utils.Config.Format, utils.Config.Width, utils.Config.Height); err != nil {
		log.Fatal(err)
	}

	go camera.ReadAndStream(pool)

	log.Printf("Listening on %s...\n", utils.Config.Addr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/view", viewPage)
	http.HandleFunc("/about", aboutPage)
	http.HandleFunc("/settings", settingsPage)
	http.HandleFunc("/stream", stream)
	http.ListenAndServe(utils.Config.Addr, nil)
}
