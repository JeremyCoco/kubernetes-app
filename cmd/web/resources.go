package main

import (
	"net/http"
	"strings"
)

func (app *application) serveStaticFiles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			app.notFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
