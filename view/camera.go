package view

import (
	"fmt"

	"github.com/blackjack/webcam"
)

//NewCamera - инициализация камеры
func NewCamera(d string, wb bool) (*Camera, error) {
	var err error
	cam := new(Camera)
	if cam.Webcam, err = webcam.Open(fmt.Sprintf("/dev/%s", d)); err != nil {
		return nil, err
	}
	if wb {
		if err := cam.SetAutoWhiteBalance(true); err != nil {
			return nil, err
		}
	}
	return cam, nil
}

//Camera - устройство камеры
type Camera struct {
	*webcam.Webcam
}
