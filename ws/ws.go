package ws

import (
	"github.com/gorilla/websocket"
)

//Message - message object model
type Message struct {
	Type int
	Body []byte
}

//Read - websocket read function
func Read(conn *websocket.Conn) (*Message, error) {
	var err error
	m := new(Message)
	m.Type, m.Body, err = conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return m, nil
}

//Send - websocket send function
func Send(conn *websocket.Conn, m *Message) error {
	err := conn.WriteMessage(m.Type, m.Body)
	if err != nil {
		return err
	}
	return nil
}
