package main

import (
	"net/http"
	"time"

	"github.com/JeremyCoco/kubernetes-app/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := getConfigFlags()
	if err != nil {
		errLogger.Fatalln(err)
	}

	templates, err := createTemplatesMap("ui/html")
	if err != nil {
		errLogger.Fatalln(err)
	}

	db, err := createDatabasePool(config.dsn)
	if err != nil {
		errLogger.Fatalln(err)
	}
	defer db.Close()

	app := &application{
		templates:      templates,
		infoLogger:     infoLogger,
		errLogger:      errLogger,
		sessionManager: createSessionManager(db),
		userModel:      &mysql.UserModel{DB: db},
		todoModel:      &mysql.TodoModel{DB: db},
	}

	server := &http.Server{
		Addr:     config.addr,
		Handler:  app.getRouter(),
		ErrorLog: errLogger,

		// Server connection timeout settings.
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	infoLogger.Printf("Listening on %s\n", config.addr)
	errLogger.Fatalln(server.ListenAndServe())
}
