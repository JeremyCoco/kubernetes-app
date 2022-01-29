package models

import "time"

type User struct {
	Id              int
	Username        string
	hashed_password string
}

type Todo struct {
	Id        int
	Title     string
	Done      bool
	CreatedAt time.Time
	Expires   time.Time
}
