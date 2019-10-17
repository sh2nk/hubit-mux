package view

import (
	"bytes"
	"hubit-mux/utils"
	"log"

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
		f.post.Write(frame)
		f.stream.Write(frame)

		func() {
			pool.RLock()
			defer pool.RUnlock()

			for name := range pool.Streams {
				pool.Streams[name] <- f.stream
			}
		}()

		if !utils.Config.Debug {
			go post(utils.Config.StreamURL, f.post, int(utils.Config.Width), int(utils.Config.Height), int(utils.Config.Resize))
		}
	}
}
