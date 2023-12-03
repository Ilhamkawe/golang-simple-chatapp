package handlers

import (
	"belajar-go-websocket/helpers"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	helpers.RenderHTML(w, "views/index.html", nil)
}
