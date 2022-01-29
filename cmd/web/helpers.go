package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *viewData) {
	ts, ok := app.templates[name+".go.html"]

	if !ok {
		app.serverError(w, fmt.Errorf("the template %q does not exist", name))
		return
	}

	buffer := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buffer, name+".go.html", app.addDefaults(data, r))

	if err != nil {
		app.serverError(w, err)
		return
	}

	buffer.WriteTo(w)
}

func (app *application) addDefaults(data *viewData, r *http.Request) *viewData {
	if data == nil {
		data = &viewData{}
	}

	data.Year = time.Now().Year()
	data.CSRFToken = nosurf.Token(r)
	data.IsAuthenticated = app.isAuthenticated(r)

	if data.IsAuthenticated {
		user, err := app.userModel.GetById(app.sessionManager.GetInt(r.Context(), "userID"))
		if err != nil {
			app.errLogger.Println(err)
		} else {
			data.User = user
		}
	}

	return data
}

func (app *application) isAuthenticated(r *http.Request) bool {
	value, ok := r.Context().Value(contextKeyIsUserAuthenticated).(bool)
	if !ok {
		return false
	}

	return value
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLogger.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (app *application) renewSessionToken(r *http.Request) error {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	return nil
}
