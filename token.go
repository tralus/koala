package koala

import (
	"errors"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/tralus/koala/auth"
)

// Token represents a Token
type Token struct {
	Value string `json:"token"`
}

// NewToken creates an instance of Token
func NewToken(value string) Token {
	return Token{value}
}

// TokenService defines an interface for a token service
type TokenService interface {
	GenerateToken(details auth.UserDetails) (Token, error)
}

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

// JwtTokenService represents a jwt token service
// It uses an AuthService for the authentication logic
// It generates a jwt token from UserDetails data
type JwtTokenService struct {
	AuthService   auth.DefaultService
	SigningMethod jwt.SigningMethod
	JwtConfig     JwtConfig
}

// NewJwtTokenService creates a new instancer of JwtTokenService
func NewJwtTokenService(a auth.DefaultService, m jwt.SigningMethod, c JwtConfig) JwtTokenService {
	return JwtTokenService{a, m, c}
}

// GetJwtClaimns gets the claimn param from request context
func GetJwtClaimns(r *http.Request, key string) (interface{}, error) {
	if c := context.Get(r, "jwtClaimns"); c != nil {
		claimns := c.(map[string]interface{})

		if v, ok := claimns[key]; ok {
			return v, nil
		}
	}

	return nil, errors.New("JWT Claimns with key " + key + " not found.")
}

// Authenticate authenticates via AuthService and generates a jwt token
func (s JwtTokenService) Authenticate(username string, pwd string) (Token, error) {
	var token Token

	user, err := s.AuthService.Authenticate(username, pwd)

	if err != nil {
		return token, err
	}

	token, err = s.GenerateToken(user)

	if err != nil {
		return token, err
	}

	return token, nil
}

// GenerateToken generates a token with UserDetails data
func (s JwtTokenService) GenerateToken(details auth.UserDetails) (Token, error) {
	var token Token

	duration := time.Hour * time.Duration(s.JwtConfig.Exp)

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(duration).Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        details.Username,
	}

	jwtToken := jwt.NewWithClaims(s.SigningMethod, claims)

	tokenString, err := jwtToken.SignedString([]byte(s.JwtConfig.Secret))

	if err != nil {
		return token, err
	}

	return NewToken(tokenString), nil
}
