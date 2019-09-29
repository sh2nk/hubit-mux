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
		req, err := http.NewRequest("POST", url, buf)
		req.Header.Set("Content-Type", "image/jpeg")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Falied POST ", err)
		}
		defer resp.Body.Close()

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Falied reading resp ", err)
		}
		log.Println(string(bytes))
	}
}
