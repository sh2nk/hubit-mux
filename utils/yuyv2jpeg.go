package utils

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
)

//Yuyv2jpeg - конвертирует кадр из YCbCr в JPEG
func Yuyv2jpeg(frame []byte, w, h uint32) *bytes.Buffer {
	buf := &bytes.Buffer{}
	yuyv := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
	for i := range yuyv.Cb {
		ii := i * 4
		yuyv.Y[i*2] = frame[ii]
		yuyv.Y[i*2+1] = frame[ii+2]
		yuyv.Cb[i] = frame[ii+1]
		yuyv.Cr[i] = frame[ii+3]
	}
	if err := jpeg.Encode(buf, yuyv, nil); err != nil {
		log.Fatal("Convert error ", err)
	}
	return buf
}
