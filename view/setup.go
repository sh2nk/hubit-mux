package view

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/blackjack/webcam"
)

//Setup - настройка камеры
func (cam *Camera) Setup(s string, maxwidth uint32, maxheight uint32) error {
	formatDesc := cam.GetSupportedFormats()
	var format webcam.PixelFormat
	for f, desc := range formatDesc {
		if strings.Contains(desc, s) {
			format = f
		}
	}
	if format == 0 {
		err := fmt.Errorf("camera setup: format %s not supported by device", s)
		return err
	}

	frames := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frames)

	f, w, h, err := cam.SetImageFormat(format, maxwidth, maxheight)
	if err != nil {
		return err
	}
	log.Printf("Resulting image format: %s (%dx%d)\n", formatDesc[f], w, h)
	return nil
}
