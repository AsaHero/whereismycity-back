package inerr

import "errors"

var (
	ErrorIncorrectPassword = errors.New("incorrect password")
)

// error not found
type ErrNotFound struct {
	name string
}

func (e *ErrNotFound) Error() string {
	return e.name + " not found"
}

func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

func NewErrNotFound(text string) *ErrNotFound {
	return &ErrNotFound{text}
}

// error conflict
type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return e.name + " already exist"
}

func IsErrConflict(err error) bool {
	_, ok := err.(*ErrConflict)
	return ok
}

func NewErrConflict(text string) *ErrConflict {
	return &ErrConflict{text}
}

// error no changes
type ErrNoChanges struct {
	name string
}

func (e *ErrNoChanges) Error() string {
	return e.name + " is not changed"
}

func IsErrNoChanges(err error) bool {
	_, ok := err.(*ErrNoChanges)
	return ok
}

func NewErrNoChanges(text string) *ErrNoChanges {
	return &ErrNoChanges{text}
}

// ErrValidation represents different types of token validation errors
type ErrJwtValidation struct {
	Message string
}

func (e ErrJwtValidation) Error() string {
	return e.Message
}

func IsErrJwtValidation(err error) bool {
	_, ok := err.(*ErrJwtValidation)
	return ok
}

func NewErrJwtValidation(message string) *ErrJwtValidation {
	return &ErrJwtValidation{
		Message: message,
	}
}
