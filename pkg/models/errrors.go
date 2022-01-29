package models

import "errors"

var (
	ErrNoRecords          = errors.New("models: no records found")
	ErrDuplicatedUsername = errors.New("models: duplicated username")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
)
