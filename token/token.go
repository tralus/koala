package token

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/context"
)

const keyTokenContext = "koala.token.0"
const keyJwtClaimsContext = "koala.jwt.claims.0"

// Token represents a Token
type Token struct {
	Value string `json:"token"`
}

// New creates an instance of Token
func New(value string) Token {
	return Token{value}
}

// // Service defines an interface for a token service
// type Service interface {
// 	GenerateToken(details auth.UserDetails) (Token, error)
// }

// JwtConfig represents the jwt settings
type JwtConfig struct {
	Exp    int
	Secret string
}

// NewJwtConfig creates an instance for JwtConfig
func NewJwtConfig(e int, s string) JwtConfig {
	if e == 0 {
		e = 72 // (7 (days) * 24 (hours)) - a week
	}

	return JwtConfig{e, s}
}

// JwtToken represents a jwt token service
// It uses an AuthService for the authentication logic
// It generates a jwt token from UserDetails data
type JwtToken struct {
	// AuthService   auth.DefaultService
	SigningMethod jwt.SigningMethod
	JwtConfig     JwtConfig
}

// NewJwtToken creates a new instancer of JwtTokenService
func NewJwtToken(m jwt.SigningMethod, c JwtConfig) JwtToken {
	return JwtToken{m, c}
}

// GenerateToken generates a token with UserDetails data
func (s JwtToken) GenerateToken(claims jwt.Claims) (t Token, err error) {
	jwtToken := jwt.NewWithClaims(s.SigningMethod, claims)

	tokenStr, err := jwtToken.SignedString([]byte(s.JwtConfig.Secret))

	if err != nil {
		return t, err
	}

	return New(tokenStr), nil
}

// ToContext puts a Token instance to the request context
func ToContext(r *http.Request, t Token) {
	context.Set(r, keyTokenContext, t)
}

// FromContext gets a Token instance from the request context
func FromContext(r *http.Request) (Token, error) {
	var token Token

	u := context.Get(r, keyTokenContext)

	if u == nil {
		errMsg := "Token (" + keyTokenContext + ") is not into context."
		return token, errors.New(errMsg)
	}

	token, ok := u.(Token)

	if !ok {
		errMsg := "The token into context is not a Token instance."
		return token, errors.New(errMsg)
	}

	return token, nil
}

// JwtClaimsToContext puts claims to the request context
func JwtClaimsToContext(r *http.Request, c *jwt.StandardClaims) {
	context.Set(r, keyJwtClaimsContext, c)
}

// JwtClaimsFromContext gets claims from the request context
func JwtClaimsFromContext(r *http.Request) (jwt.StandardClaims, error) {
	var claims jwt.StandardClaims

	u := context.Get(r, keyJwtClaimsContext)

	if u == nil {
		errMsg := "StandardClaims (" + keyJwtClaimsContext + ") is not into context."
		return claims, errors.New(errMsg)
	}

	claims, ok := u.(jwt.StandardClaims)

	if !ok {
		errMsg := "The claims into context is not a StandardClaims instance."
		return claims, errors.New(errMsg)
	}

	return claims, nil
}
