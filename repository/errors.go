package repository

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoUser            = errors.New("user not found")
)
