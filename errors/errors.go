// Package errors intentionally mirrors the standard "errors" module.
// All packages on Koala should use this.
package errors

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// BaseError represents the base error
type BaseError struct {
	err error
}

// Error implements error.Error
func (e BaseError) Error() string {
	return fmt.Sprintf("%s", e.err)
}

// GetStack implements RootError
func (e BaseError) GetStack() string {
	return fmt.Sprintf("%+v\n", e)
}

// Format formats the error message
func (e BaseError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v\n", e.err)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

// NewBaseError creates a BaseError instance
func NewBaseError(err error) BaseError {
	return BaseError{err}
}

// New only wraps errors.New
func New(m string) error {
	return errors.New(m)
}

// Cause only wraps errors.Cause
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrap only wraps errors.Wrap
func Wrap(err error, m string) error {
	return errors.Wrap(err, m)
}

// Errorf only wraps errors.Errorf
func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// RootError interface exposes additional information about the error
type RootError interface {
	// Returns the stack trace without the error message.
	GetStack() string

	Error() string
}

// IsRootError verifies if the error is a RuntimeError
func IsRootError(e error) bool {
	_, ok := e.(RootError)
	return ok
}

// NotFoundError represents a not found error
// It implements the RuntimeError interface
type NotFoundError struct {
	BaseError
}

// NewNotFoundError creates an error from NotFoundError
func NewNotFoundError(err error) error {
	return NotFoundError{NewBaseError(err)}
}

// IsNotFoundError verifies if error is a NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := errors.Cause(err).(NotFoundError)
	return ok
}

// IllegalStateError represents an illegal state error
// It implements the RuntimeError interface
type IllegalStateError struct {
	BaseError
}

// NewIllegalStateError creates an IllegalStateError instance
func NewIllegalStateError(err error) error {
	return IllegalStateError{NewBaseError(err)}
}

// IsIllegalStateError verifies if error is an IllegalStateError
func IsIllegalStateError(err error) bool {
	_, ok := errors.Cause(err).(IllegalStateError)
	return ok
}

// IllegalArgumentError represents an illegal argument error
// It implements the RuntimeError interface
type IllegalArgumentError struct {
	BaseError
}

// NewIllegalArgumentError creates an IllegalArgumentError instance
func NewIllegalArgumentError(err error) error {
	return IllegalArgumentError{NewBaseError(err)}
}

// IsIllegalArgumentError verifies if error is an IllegalStateError
func IsIllegalArgumentError(err error) bool {
	_, ok := errors.Cause(err).(IllegalStateError)
	return ok
}

// RelationshipError represents a relationship error
type RelationshipError struct {
	BaseError
}

// NewRelationshipError creates an error from RelationshipError
func NewRelationshipError(err error) error {
	return RelationshipError{NewBaseError(err)}
}

// IsRelationshipError verifies if error is an RelationshipError
func IsRelationshipError(err error) bool {
	_, ok := errors.Cause(err).(RelationshipError)
	return ok
}

// NotAuthorizedError represents the error for not authorized logic
type NotAuthorizedError struct {
	BaseError
}

// NewNotAuthorizedError creates a NotAuthorizedError instance
func NewNotAuthorizedError(err error) error {
	return NotAuthorizedError{NewBaseError(err)}
}

// IsNotAuthorizedError verifies if error is an NotAuthorizedError
func IsNotAuthorizedError(err error) bool {
	_, ok := errors.Cause(err).(NotAuthorizedError)
	return ok
}
