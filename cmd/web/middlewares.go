package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/JeremyCoco/kubernetes-app/pkg/models"
	"github.com/justinas/nosurf"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func addSecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	CSRFHandler := nosurf.New(next)
	CSRFHandler.SetBaseCookie(http.Cookie{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	return CSRFHandler
}

func (app *application) verifyUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.sessionManager.Exists(r.Context(), "userID") {
			next.ServeHTTP(w, r)
			return
		}

		_, err := app.userModel.GetById(app.sessionManager.GetInt(r.Context(), "userID"))
		if err != nil {
			if errors.Is(err, models.ErrNoRecords) {
				app.sessionManager.Remove(r.Context(), "userID")
				next.ServeHTTP(w, r)
				return
			}

			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyIsUserAuthenticated, true)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
