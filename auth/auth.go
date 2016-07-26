package auth

import (
	"net/http"
	"crypto/sha256"
	"encoding/hex"
	
	"github.com/gorilla/context"
	"github.com/tralus/koala/errors"
)

// Generic error type for user not found logic
type UsernameNotFoundError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e UsernameNotFoundError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e UsernameNotFoundError) Error() string {
	return e.Msg
}

// It creates an UsernameNotFoundError instance
func NewUsernameNotFoundError(m string) error {
	_, s := errors.StackTrace()
	return UsernameNotFoundError{m, s} 
}

// It verifies if error is a UsernameNotFoundError type
func IsUsernameNotFoundError(e error) bool {
	_, ok := e.(UsernameNotFoundError)
	return ok
}

// Generic error type for not authorized logic
type NotAuthorizedError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e NotAuthorizedError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e NotAuthorizedError) Error() string {
	return e.Msg
}

func NewNotAuthorizedError(m string) error {
	_, s := errors.StackTrace()
	return NotAuthorizedError{m, s} 
}

// It creates a NotAuthorizedError instance
func IsNotAuthorizedError(e error) bool {
	_, ok := e.(NotAuthorizedError)
	return ok
}

// UserDetails represents basic data of the user
// The username should be a unique value as an email field
type UserDetails struct {
  	Username string
  	Password string
  	IsActive bool
}

// It gets a UserDetails instance from request context
func ContextUser(r *http.Request) (UserDetails, error) {
	var details UserDetails
	
	if u := context.Get(r, "user"); u != nil {
		details, ok := u.(UserDetails)
		if (!ok) {
			message := "The user on context is not an UserDetails instance."
			return details, errors.New(message)
		}
		return details, nil
	}

	return details, errors.New("Key user not found on context.")
}

// UserDetailsService is an interface implemented to create an UserDetails
type UserDetailsService interface {
  	LoadUserByUsername(username string) (UserDetails, error)
}

// PasswordStrategy is a interface to create password logic
type PasswordStrategy interface {
  	Exec(password string) string
}

// Sha256PasswordStrategy implements the logic for hash sha256
type Sha256PasswordStrategy struct {}

// It hashs the password usgin sha256
func (s Sha256PasswordStrategy) Exec(password string) string {
	digest := sha256.New()
	digest.Write([]byte(password))
	return hex.EncodeToString(digest.Sum(nil))
}

// It creates Sha256PasswordStrategy instance
func NewSha256Password() Sha256PasswordStrategy {
  	return Sha256PasswordStrategy{}
}

// It a helper method to create a Sha256PasswordStrategy and to return a hashed password 
func Sha256Password(password string) string {
  	return NewSha256Password().Exec(password)
}

// AuthService executes the logic that performs the comparison of passwords
type AuthService struct {
  	UserDetailsService UserDetailsService
  	
	// It is used to get the hash of password
  	PasswordStrategy PasswordStrategy
}

// It creates an AuthService instance
func NewAuthService(u UserDetailsService, s PasswordStrategy) AuthService {
  	return AuthService{u, s}
}

// It authenticates using UserDetailsService and after it compares the passwords
func (auth AuthService) Authenticate(username string, password string) (UserDetails, error) {
  	user, err := auth.UserDetailsService.LoadUserByUsername(username)
  
  	if (err != nil) {
    	return user, err
  	}
  
  	old := auth.PasswordStrategy.Exec(password)
  
  	if (user.Password != old) {
  		return user, NewNotAuthorizedError("Not authorized for password")
  	}
    
  	return user, nil
}