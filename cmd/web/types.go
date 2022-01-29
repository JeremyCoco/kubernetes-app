package main

import (
	"html/template"
	"log"

	"github.com/JeremyCoco/kubernetes-app/pkg/forms"
	"github.com/JeremyCoco/kubernetes-app/pkg/models"
	"github.com/JeremyCoco/kubernetes-app/pkg/models/mysql"
)

type application struct {
	templates      map[string]*template.Template
	infoLogger     *log.Logger
	errLogger      *log.Logger
	sessionManager *scs.SessionManager
	userModel      *mysql.UserModel
	todoModel      *mysql.TodoModel
}

type viewData struct {
	Year            int
	IsAuthenticated bool
	CSRFToken       string
	Form            *forms.Form
	User            *models.User
	Todos           []*models.Todo
}

type configFlags struct {
	addr string // Server network address.
	dsn  string // Database data source name.
}
