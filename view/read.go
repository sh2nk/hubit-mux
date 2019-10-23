package view

import (
	"bytes"
	"hubit-mux/utils"
	"io/ioutil"
	"log"
	"strings"

	"github.com/blackjack/webcam"
)

type framebuf struct {
	stream *bytes.Buffer
	post   *bytes.Buffer
}

//ReadAndStream - читаем поток с камеры
func (cam *Camera) ReadAndStream(pool *utils.Pool) {
	var err error
	var frame []byte

	// Универсальная обработка ошибок
	defer func() {
		if rec := recover(); rec != nil {
			// Поскольку нет особой обработки ошибки, cast к error без надобности, сразу принт
			log.Printf("Read panic: %+v", rec)
		}

		if err != nil {
			log.Printf("Read error: %+v", err)
		}
	}()

	if err = cam.StartStreaming(); err != nil {
		return
	}
	defer cam.StopStreaming()

	for {
		//Таймаут 5 секунд
		if err = cam.WaitForFrame(5); err != nil {
			if _, ok := err.(*webcam.Timeout); ok {
				continue
			}
			return
		}

		if frame, err = cam.ReadFrame(); err != nil {
			return
		}

		if len(frame) == 0 {
			continue
		}

		if strings.Contains("yuyv", strings.ToLower(utils.Config.Format)) {
			if frame, err = ioutil.ReadAll(utils.Yuyv2jpeg(frame, utils.Config.Width, utils.Config.Height)); err != nil {
				return
			}
		}

		func() {
			pool.RLock()
			defer pool.RUnlock()

			for name := range pool.Streams {
				pool.Streams[name] <- bytes.NewBuffer(frame)
			}
		}()

		if !utils.Config.Debug {
			go post(utils.Config.StreamURL, bytes.NewBuffer(frame), utils.Config.Width, utils.Config.Height, int(utils.Config.Resize))
		}
	}
}
