package token

import (
	"errors"
	"net/http"

	"github.com/tralus/koala/context"
)

const keyTokenContext = "koala.token.0"

// Token represents a Token
type Token struct {
	Value string `json:"token"`
}

// New creates an instance of Token
func New(value string) Token {
	return Token{value}
}

// ToContext puts a Token instance to the request context
func ToContext(r *http.Request, t Token) {
	context.Add(r, keyTokenContext, t)
}

// FromContext gets a Token instance from the request context
func FromContext(r *http.Request) (Token, error) {
	var token Token

	value, err := context.Get(r, keyTokenContext)

	if err != nil {
		return token, err
	}

	token, ok := value.(Token)

	if !ok {
		errMsg := "The token into context is not a Token instance."
		return token, errors.New(errMsg)
	}

	return token, nil
}
