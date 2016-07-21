package auth

import (
	"net/http"
	"crypto/sha256"
	"encoding/hex"
	
	"github.com/gorilla/context"
	"github.com/tralus/koala/errors"
)

type UsernameNotFoundError struct {
	Msg string
	Stack string
}

func (e UsernameNotFoundError) GetStack() string {
	return e.Stack
}

func (e UsernameNotFoundError) Error() string {
	return e.Msg
}

func NewUsernameNotFoundError(m string) error {
	_, s := errors.StackTrace()
	return UsernameNotFoundError{m, s} 
}

func IsUsernameNotFoundError(e error) bool {
	_, ok := e.(UsernameNotFoundError)
	return ok
}

type NotAuthorizedError struct {
	Msg string
	Stack string
}

func (e NotAuthorizedError) GetStack() string {
	return e.Stack
}

func (e NotAuthorizedError) Error() string {
	return e.Msg
}

func NewNotAuthorizedError(m string) error {
	_, s := errors.StackTrace()
	return NotAuthorizedError{m, s} 
}

func IsNotAuthorizedError(e error) bool {
	_, ok := e.(NotAuthorizedError)
	return ok
}

type UserDetails struct {
  	Username string
  	Password string
  	IsActive bool
}

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

type UserDetailsService interface {
  	LoadUserByUsername(username string) (UserDetails, error)
}

type PasswordStrategy interface {
  	Exec(password string) string
}

type Sha256PasswordStrategy struct {}

func (s Sha256PasswordStrategy) Exec(password string) string {
	digest := sha256.New()
	digest.Write([]byte(password))
	return hex.EncodeToString(digest.Sum(nil))
}

func NewSha256Password() Sha256PasswordStrategy {
  	return Sha256PasswordStrategy{}
}

func Sha256Password(password string) string {
  	return NewSha256Password().Exec(password)
}

type AuthService struct {
  	UserDetailsService UserDetailsService
  	PasswordStrategy PasswordStrategy
}

func NewAuthService(u UserDetailsService, s PasswordStrategy) AuthService {
  	return AuthService{u, s}
}

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