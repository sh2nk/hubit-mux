package main

import (
	"bytes"
	"fmt"
	"html/template"
	"hubit-mux/view"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	"hubit-mux/ws"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

func wsServer(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Upgrade connection error:", err)
	}
	log.Printf("Upgraded new connection from %s\n", conn.RemoteAddr())

	name := fmt.Sprintf("%v", conn.RemoteAddr())
	wsPool.Clients[name] = conn

	go func() {
		for {
			message, err := ws.Read(conn)
			if err != nil {
				log.Println(err)
				conn.Close()
				delete(wsPool.Clients, name)
				return
			}
			for _, val := range wsPool.Clients {
				err := ws.Send(val, message)
				if err != nil {
					log.Println(err)
					conn.Close()
					delete(wsPool.Clients, name)
					return
				}
			}
			log.Printf("%s sent: %s\n", conn.RemoteAddr(), string(message.Body))
		}
	}()
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view", 301)
}

func viewPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Parsing error ", err)
	}
	tmpl.ExecuteTemplate(w, "index", nil)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Parsing error ", err)
	}
	tmpl.ExecuteTemplate(w, "about", nil)
}

func settingsPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal("Parsing error ", err)
	}
	tmpl.ExecuteTemplate(w, "settings", nil)
}

func getFaceData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string(view.RawFace))
}

func stream(w http.ResponseWriter, r *http.Request) {
	multipartWriter := multipart.NewWriter(w)

	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+multipartWriter.Boundary())

	stream := make(chan *bytes.Buffer)
	name := uuid.Must(uuid.NewV4()).String()

	func() {
		streamPool.Lock()
		defer streamPool.Unlock()
		streamPool.Streams[name] = stream
	}()
	defer func() {
		streamPool.Lock()
		defer streamPool.Unlock()
		delete(streamPool.Streams, name)
	}()

	for buf := range stream {
		iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-type":   []string{"image/jpeg"},
			"Content-length": []string{strconv.Itoa(buf.Len())},
		})
		if err != nil {
			log.Printf("Потеря связи с клиентом: %+v", err)
			return
		}

		if _, err = io.Copy(iw, buf); err != nil {
			log.Printf("Потеря связи с клиентом: %+v", err)
			return
		}
	}
}
