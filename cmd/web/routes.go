package main

import (
	"belajar-go-websocket/handlers"
	"github.com/bmizerany/pat"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/websocket", http.HandlerFunc(handlers.WsEndpoint))
	return mux
}
