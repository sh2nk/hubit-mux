package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Upgrade connection error:", err)
	}
	log.Printf("Upgraded new connection from %s\n", conn.RemoteAddr())

	name := fmt.Sprintf("%v", conn.RemoteAddr())
}
