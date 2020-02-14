package utils

import (
	"bytes"
	"sync"

	"github.com/gorilla/websocket"
)

//StreamPool - структура пула видеопотоков
type StreamPool struct {
	sync.RWMutex
	Streams map[string]chan *bytes.Buffer
}

//WSPool - структура пула клиентов вебсокет сервера
type WSPool struct {
	Clients map[string]*websocket.Conn
}
