package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gorilla/context"
	"github.com/tralus/koala/errors"
)

const keyUserContext = "koala.auth.0"

// UsernameExistsError represents error for username exists logic
type UsernameExistsError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e UsernameExistsError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e UsernameExistsError) Error() string {
	return e.Msg
}

// NewUsernameExistsError creates an UsernameExistsError instance
func NewUsernameExistsError(m string) error {
	_, s := errors.StackTrace()
	return UsernameExistsError{m, s}
}

// IsUsernameExistsError verifies if error is an UsernameExistsError
func IsUsernameExistsError(e error) bool {
	_, ok := e.(UsernameExistsError)
	return ok
}

// UsernameNotFoundError represents error for user not found logic
type UsernameNotFoundError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e UsernameNotFoundError) GetStack() string {
	return e.Stack
}

// Built-in interface
func (e UsernameNotFoundError) Error() string {
	return e.Msg
}

// NewUsernameNotFoundError creates an UsernameNotFoundError instance
func NewUsernameNotFoundError(m string) error {
	_, s := errors.StackTrace()
	return UsernameNotFoundError{m, s}
}

// IsUsernameNotFoundError verifies if error is an UsernameNotFoundError
func IsUsernameNotFoundError(e error) bool {
	_, ok := e.(UsernameNotFoundError)
	return ok
}

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

// UserDetails represents the user data
// The username should be a unique value as an email field
type UserDetails struct {
	Username string
	Password string
	IsActive bool
}

// ToContext puts an UserDetails instance to the request context
func ToContext(r *http.Request, u UserDetails) {
	context.Set(r, keyUserContext, u)
}

// FromContext gets an UserDetails instance from the request context
func FromContext(r *http.Request) (UserDetails, error) {
	var details UserDetails

	if u := context.Get(r, keyUserContext); u != nil {
		details, ok := u.(UserDetails)
		if !ok {
			m := "The user into context is not an UserDetails instance."
			return details, errors.New(m)
		}
		return details, nil
	}

	return details, errors.New("Key " + keyUserContext + " not found in the context.")
}

// UserDetailsService defines an interface to implement UserDetails logic
type UserDetailsService interface {
	LoadUserByUsername(username string) (UserDetails, error)
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

// DefaultService executes the logic that performs the comparison of passwords
type DefaultService struct {
	UserDetailsService UserDetailsService

	// It is used to get the hash of password
	PasswordStrategy PasswordStrategy
}

// NewDefaultService creates an DefaultService instance
func NewDefaultService(u UserDetailsService, s PasswordStrategy) DefaultService {
	return DefaultService{u, s}
}

// Authenticate uses UserDetailsService and after it compares the passwords
func (auth DefaultService) Authenticate(username string, password string) (UserDetails, error) {
	user, err := auth.UserDetailsService.LoadUserByUsername(username)

	if err != nil {
		return user, err
	}

	old := auth.PasswordStrategy.Exec(password)

	if user.Password != old {
		return user, NewNotAuthorizedError("Credentials not authorized.")
	}

	return user, nil
}
