package models

import "errors"

var (
	ErrNoResults     = errors.New("models: no results found")
	ErrEmailTaken    = errors.New("models: email address already taken")
	ErrNotFound      = errors.New("models: resource could not be found")
	ErrWrongPassword = errors.New("models: password does not match hash")
	ErrNoEmail       = errors.New("models: email could not be found")
)
