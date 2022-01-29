package main

import (
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	_ "github.com/go-sql-driver/mysql"
)

type contextKey string

var (
	contextKeyIsUserAuthenticated = contextKey("IsUserAuthenticated")
)

func createSessionManager(db *sql.DB) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = time.Hour
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Secure = true
	sessionManager.Store = mysqlstore.New(db)

	return sessionManager
}
