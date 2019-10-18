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

	// Универсальная обработка ошибок
	defer func() {
		if rec := recover(); rec != nil {
			// паника!
			var ok bool
			if err, ok = rec.(error); !ok {
				log.Printf("Read panic: %#v", rec)
			}
		}

		if err != nil {
			log.Printf("Read error: %#v", err)
		}
	}()

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

		f := framebuf{new(bytes.Buffer), new(bytes.Buffer)}
		f.post.Grow(len(frame) + 100)
		f.stream.Grow(len(frame) + 100)

		func() {
			pool.RLock()
			defer pool.RUnlock()

			switch {
			case strings.Contains("yuyv", strings.ToLower(utils.Config.Format)):
				bytes, err := ioutil.ReadAll(utils.Yuyv2jpeg(frame, utils.Config.Width, utils.Config.Height))
				if err != nil {
					return
				}
				f.stream.Write(bytes)
				f.stream.Write(bytes)
				for name := range pool.Streams {
					pool.Streams[name] <- f.stream
				}
			default:
				f.stream.Write(frame)
				f.post.Write(frame)
				for name := range pool.Streams {
					pool.Streams[name] <- f.stream
				}
			}
		}()

		if !utils.Config.Debug {
			go post(utils.Config.StreamURL, f.post, utils.Config.Width, utils.Config.Height, int(utils.Config.Resize))
		}
	}
}
