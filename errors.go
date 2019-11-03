package dotweb

import "errors"

// ErrValidatorNotRegistered error for not register Validator
var ErrValidatorNotRegistered = errors.New("validator not registered")

// ErrNotFound error for not found file
var ErrNotFound = errors.New("not found file")
