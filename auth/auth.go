package auth

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/tralus/koala/errors"
)

// NotAuthorizedError represents the error for not authorized logic
type NotAuthorizedError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e NotAuthorizedError) GetStack() string {
	return e.Stack
}

// Built-in interface
func (e NotAuthorizedError) Error() string {
	return e.Msg
}

// NewNotAuthorizedError creates a NotAuthorizedError instance
func NewNotAuthorizedError(m string) error {
	s, _ := errors.StackTrace()
	return NotAuthorizedError{m, s}
}

// IsNotAuthorizedError verifies if error is an NotAuthorizedError
func IsNotAuthorizedError(e error) bool {
	_, ok := e.(NotAuthorizedError)
	return ok
}

// PasswordStrategy defines an interface to create password logic
type PasswordStrategy interface {
	Exec(password string) string
}

// Sha256PasswordStrategy implements the logic for hash sha256
type Sha256PasswordStrategy struct{}

// Exec generates the password using sha256
func (s Sha256PasswordStrategy) Exec(password string) string {
	digest := sha256.New()
	digest.Write([]byte(password))
	return hex.EncodeToString(digest.Sum(nil))
}

// NewSha256Password creates Sha256PasswordStrategy instance
func NewSha256Password() Sha256PasswordStrategy {
	return Sha256PasswordStrategy{}
}

// Sha256Password is a helper method to create a Sha256PasswordStrategy and to return a hashed password
func Sha256Password(password string) string {
	return NewSha256Password().Exec(password)
}
