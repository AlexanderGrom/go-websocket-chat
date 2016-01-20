package main

import (
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags)

	server := NewServer()
	log.Println("Listening server...")
	go server.Listen()

	http.Handle("/", server.Handler())
	http.ListenAndServe(":5213", nil)
}
