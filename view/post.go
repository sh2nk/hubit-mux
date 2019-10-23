package view

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Face struct {
	Contains bool      `json:"contains"`
	Emotions *Emotions `json:"emotions"`
}

type Emotions struct {
	Angry     float64 `json:"angry"`
	Disgusted float64 `json:"disgusted"`
	Fearful   float64 `json:"fearful"`
	Happy     float64 `json:"happy"`
	Neutral   float64 `json:"neutral"`
	Sad       float64 `json:"sad"`
	Surprised float64 `json:"surprised"`
}

func post(url string, frame io.Reader, w, h uint32, level int) {
	// src, _, err := image.Decode(frame)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// dest := image.NewRGBA(image.Rect(0, 0, int(w/level), int(h/level)))
	// draw.ApproxBiLinear.Scale(dest, dest.Bounds(), src, src.Bounds(), draw.Over, nil)

	// buf := new(bytes.Buffer)
	// if err := jpeg.Encode(buf, dest, nil); err != nil {
	// 	log.Fatal(err)
	// }
	resp, err := http.Post(url, "application/octet-stream", frame)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var face Face

	if err = json.NewDecoder(resp.Body).Decode(&face); err != nil {
		log.Printf("Face reading error: %+v", err)
		return
	}

	if face.Emotions != nil {
		log.Printf("%.2f%% happy", face.Emotions.Happy*100.0)
	}
}
