package jwtoken

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/tralus/koala/context"
	"github.com/tralus/koala/token"
)

const keyJwtClaimsContext = "koala.jwt.claims.0"

// Config represents the jwt settings
type Config struct {
	Exp    int
	Secret string
}

// NewConfig creates an instance for JwtConfig
func NewConfig(e int, s string) Config {
	if e == 0 {
		e = 72 // (7 (days) * 24 (hours)) - a week
	}

	return Config{e, s}
}

// Token represents a jwt token service
// It uses an AuthService for the authentication logic
// It generates a jwt token from UserDetails data
type Token struct {
	// AuthService   auth.DefaultService
	SigningMethod jwt.SigningMethod
	Config        Config
}

// New creates a new instance of TokenService
func New(m jwt.SigningMethod, c Config) Token {
	return Token{m, c}
}

// GenerateToken generates a token with UserDetails data
func (s Token) GenerateToken(claims jwt.Claims) (t token.Token, err error) {
	jwtToken := jwt.NewWithClaims(s.SigningMethod, claims)

	tokenStr, err := jwtToken.SignedString([]byte(s.Config.Secret))

	if err != nil {
		return t, err
	}

	return token.New(tokenStr), nil
}

// ClaimsToContext puts claims to the request context
func ClaimsToContext(r *http.Request, c *jwt.StandardClaims) {
	context.Add(r, keyJwtClaimsContext, c)
}

// ClaimsFromContext gets claims from the request context
func ClaimsFromContext(r *http.Request) (jwt.StandardClaims, error) {
	var claims jwt.StandardClaims

	value, err := context.Get(r, keyJwtClaimsContext)

	if err == nil {
		return claims, err
	}

	claims, ok := value.(jwt.StandardClaims)

	if !ok {
		errMsg := "The claims in the context is not a StandardClaims instance."
		return claims, errors.New(errMsg)
	}

	return claims, nil
}
