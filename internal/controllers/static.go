package controllers

import (
	"net/http"
)

// StaticHandler takes a template, and returns a http.HandlerFunc. The func it returns runs execute, and passess no data to execute.
func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}
