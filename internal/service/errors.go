package service

import "fmt"

var (
	ErrUserNotFound      = fmt.Errorf("user not found")
	ErrUserAlreadyExists = fmt.Errorf("user already exists")

	ErrCannotSignToken  = fmt.Errorf("cannot sign token")
	ErrCannotParseToken = fmt.Errorf("cannot parse token")

	ErrActorNotFound = fmt.Errorf("actor not found")
	ErrFilmNotFound  = fmt.Errorf("film not found")
)
