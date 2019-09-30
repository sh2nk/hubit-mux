package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func post(url string) {
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
		resp, err := http.Post(url, "application/octet-stream", buf)
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Falied reading resp ", err)
		}
		log.Println(string(bytes))
	}
}
