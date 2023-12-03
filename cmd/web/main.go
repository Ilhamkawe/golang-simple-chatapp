package main

import (
	"belajar-go-websocket/handlers"
	"log"
	"net/http"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting channel listener")

	_ = http.ListenAndServe(":8080", mux)
}
