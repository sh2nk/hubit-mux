package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

var tmpl *template.Template

func init() {
	var err error

	if tmpl, err = template.ParseGlob("templates/*.html"); err != nil {
		log.Fatal("Parsing error ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view", 301)
}

func viewPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index", nil)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "about", nil)
}

func settingsPage(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "settings", nil)
}

func stream(w http.ResponseWriter, r *http.Request) {
	multipartWriter := multipart.NewWriter(w)

	w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+multipartWriter.Boundary())

	stream := make(chan *bytes.Buffer)
	name := uuid.Must(uuid.NewV4()).String()

	func() {
		pool.Lock()
		defer pool.Unlock()
		pool.Streams[name] = stream
	}()
	defer func() {
		pool.Lock()
		defer pool.Unlock()
		delete(pool.Streams, name)
	}()

	for buf := range stream {
		iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
			"Content-type":   []string{"image/jpeg"},
			"Content-length": []string{strconv.Itoa(buf.Len())},
		})
		if err != nil {
			log.Printf("Потеря связи с клиентом: %+v", err)
			return
		}

		if _, err = io.Copy(iw, buf); err != nil {
			log.Printf("Потеря связи с клиентом: %+v", err)
			return
		}
	}
}
