package utils

import (
	"bytes"
	"sync"
)

//Pool - структура пула стрима
type Pool struct {
	sync.RWMutex
	Streams map[string]chan *bytes.Buffer
}
