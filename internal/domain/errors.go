package domain

import (
	"errors"
)

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

var ErrInvalidCredentials = errors.New("invalid credentials")
