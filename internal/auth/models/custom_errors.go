package models

import "errors"

var (
	ErrInvalidUsername = errors.New("username must have only latin letters without special symbols")
	ErrUsernameExists = errors.New("etot username is zanyat")
	ErrUserNotFound = errors.New("user not found")
)