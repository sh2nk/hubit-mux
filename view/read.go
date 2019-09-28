package view

import (
	"bytes"

	"github.com/blackjack/webcam"
)

//Read - читаем поток с камеры
func (cam *Camera) Read(c chan *bytes.Buffer) {
	if err := cam.StartStreaming(); err != nil {
		return
	}
	defer cam.StopStreaming()

	for {
		//Таймаут 5 секунд
		if err := cam.WaitForFrame(5); err != nil {
			if _, ok := err.(*webcam.Timeout); ok {
				continue
			}
			return
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			return
		}

		if len(frame) == 0 {
			continue
		}

		buf := new(bytes.Buffer)
		buf.Grow(len(frame) + 100)
		buf.Write(frame)

		c <- buf
	}
}
