package helpers

import (
	"net/http"
	"text/template"
)

func RenderHTML(w http.ResponseWriter, view string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles(view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
