package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/JeremyCoco/kubernetes-app/pkg/forms"
	"github.com/JeremyCoco/kubernetes-app/pkg/models"
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) getRouter() http.Handler {
	baseMiddlewares := alice.New(app.recoverPanic, app.logRequest, addSecureHeaders)
	dynamicMiddlewares := alice.New(app.sessionManager.LoadAndSave, noSurf, app.verifyUserAuth)
	dynamicwithRequireAuth := dynamicMiddlewares.Append(app.requireAuthentication)

	router := pat.New()

	router.Get("/", dynamicMiddlewares.ThenFunc(app.home))

	router.Get("/register", dynamicMiddlewares.ThenFunc(app.registerForm))
	router.Post("/auth/register", dynamicMiddlewares.ThenFunc(app.registerUser))

	router.Get("/login", dynamicMiddlewares.ThenFunc(app.loginForm))
	router.Post("/auth/login", dynamicMiddlewares.ThenFunc(app.loginUser))

	router.Get("/todos/list", dynamicwithRequireAuth.ThenFunc(app.todosManager))
	router.Get("/todos/create", dynamicwithRequireAuth.ThenFunc(app.createTodoForm))
	router.Post("/todos/create", dynamicwithRequireAuth.ThenFunc(app.createTodo))
	router.Post("/todos/complete", dynamicwithRequireAuth.ThenFunc(app.completeTodo))
	router.Post("/todos/delete", dynamicwithRequireAuth.ThenFunc(app.deleteTodo))

	router.Post("/auth/logout", dynamicwithRequireAuth.ThenFunc(app.logoutUser))

	fileServer := http.StripPrefix("/static", app.serveStaticFiles(http.FileServer(http.Dir("./ui/assets"))))
	router.Get("/static/", fileServer)

	return baseMiddlewares.Then(router)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home", nil)
}

func (app *application) registerForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register", &viewData{
		Form: forms.New(nil),
	})
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "password")
	form.NoWhiteSpace("username")
	form.MinLength("password", 10)

	if !form.IsValid() {
		app.render(w, r, "register", &viewData{
			Form: form,
		})

		return
	}

	err = app.userModel.Insert(form.Get("username"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedUsername) {
			form.Errors.Add("username", "This username is already taken")
			app.render(w, r, "register", &viewData{
				Form: form,
			})

			return
		}

		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) loginForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login", &viewData{Form: forms.New(nil)})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "password")
	form.NoWhiteSpace("username")

	if !form.IsValid() {
		app.render(w, r, "login", &viewData{
			Form: form,
		})

		return
	}

	id, err := app.userModel.Authenticate(form.Get("username"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Username or password is incorrect")
			app.render(w, r, "login", &viewData{
				Form: form,
			})

			return
		}

		app.serverError(w, err)
		return
	}

	err = app.renewSessionToken(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "userID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) todosManager(w http.ResponseWriter, r *http.Request) {
	userId := app.sessionManager.GetInt(r.Context(), "userID")
	todos, err := app.todoModel.GetAll(userId)
	if err != nil {
		app.serverError(w, err)
	}

	app.render(w, r, "todos", &viewData{
		Todos: todos,
	})
}

func (app *application) createTodoForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create", &viewData{
		Form: forms.New(nil),
	})
}

func (app *application) createTodo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.infoLogger.Printf("Creating new todo...")
	form := forms.New(r.PostForm)
	form.Required("title", "days", "hours", "minutes")
	form.MaxLength("title", 255)
	form.RequireTypeInt("days", "hours", "minutes")

	app.infoLogger.Printf("Checking if inputs are correct...")
	if !form.IsValid() {
		shouldAddExpiresErr := false
		for _, field := range []string{"days", "hours", "minutes"} {
			hasErrors := len(form.Errors[field]) > 0
			if hasErrors {
				shouldAddExpiresErr = true
				break
			}
		}

		if shouldAddExpiresErr {
			form.Errors.Add("expires", "Expiration should only contain numeric values")
		}

		app.render(w, r, "create", &viewData{
			Form: form,
		})

		return
	}

	userID := app.sessionManager.GetInt(r.Context(), "userID")

	app.infoLogger.Printf("Create creation time and expiry time...")
	days, _ := strconv.Atoi(form.Get("days"))
	hours, _ := strconv.Atoi(form.Get("hours"))
	minutes, _ := strconv.Atoi(form.Get("minutes"))

	createdAt := time.Now().UTC()
	expires := createdAt.AddDate(0, 0, days)
	expires = expires.Add((time.Hour * time.Duration(hours)) + (time.Minute * time.Duration(minutes)))

	app.infoLogger.Printf("Create the todo struct with the data that will be inserted into the database...")
	todo := models.Todo{
		Title:     form.Get("title"),
		CreatedAt: createdAt,
		Expires:   expires,
	}

	err = app.todoModel.Insert(userID, todo)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLogger.Printf("New todo created!")
	http.Redirect(w, r, "/todos/list", http.StatusSeeOther)
}

func (app *application) completeTodo(w http.ResponseWriter, r *http.Request) {
	app.infoLogger.Printf("Starting process of complete todo")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("todo_id")
	form.RequireTypeInt("todo_id")

	if !form.IsValid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userId := app.sessionManager.GetInt(r.Context(), "userID")
	todoId, _ := strconv.Atoi(form.Get("todo_id"))

	err = app.todoModel.MarkAsComplete(userId, todoId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		app.serverError(w, err)
		return
	}

	app.infoLogger.Printf("Todo marked as done")
	http.Redirect(w, r, "/todos/list", http.StatusSeeOther)
}

func (app *application) deleteTodo(w http.ResponseWriter, r *http.Request) {
	app.infoLogger.Printf("Starting process of delete todo")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("todo_id")
	form.RequireTypeInt("todo_id")

	if !form.IsValid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userId := app.sessionManager.GetInt(r.Context(), "userID")
	todoId, _ := strconv.Atoi(form.Get("todo_id"))

	err = app.todoModel.Delete(userId, todoId)
	if err != nil {
		if errors.Is(err, models.ErrNoRecords) {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		app.serverError(w, err)
		return
	}

	app.infoLogger.Printf("Todo deleted")
	http.Redirect(w, r, "/todos/list", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	err := app.renewSessionToken(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
