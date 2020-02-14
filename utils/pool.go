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

//WSPool - пул соединений для вебсокета
var WSPool map[string]*websocket.Conn
