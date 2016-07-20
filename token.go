package koala

import (
	"errors"
	"net/http"
	"time"
	
	"github.com/gorilla/context"
	"github.com/tralus/koala/auth"
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	Value string `json:"token"`
}

func NewToken(value string) Token {
	return Token{value}
}

type TokenService interface {
	GenerateToken(details auth.UserDetails) (Token, error)
}

type JwtTokenService struct {
	AuthService auth.AuthService
	SigningMethod jwt.SigningMethod
	JwtConfig JwtConfig
}

func NewJwtTokenService(a auth.AuthService, m jwt.SigningMethod, j JwtConfig) JwtTokenService {
	return JwtTokenService{a, m, j}
}

func GetJwtClaimns(r *http.Request, key string) (interface{}, error) {
	if c := context.Get(r, "jwtClaimns"); c != nil {
		claimns := c.(map[string]interface{})
		
		if v, ok := claimns[key]; ok {
			return v, nil			
		}
	}

	return nil, errors.New("JWT Claimns with key " + key + " not found.")
}

func (s JwtTokenService) Authenticate(username string, pwd string) (Token, error) {
	var token Token
	
	user, err := s.AuthService.Authenticate(username, pwd)
	
	if (err != nil) {
		return token, err
	}
		
	token, err = s.GenerateToken(user)
		
	if (err != nil) {
		return token, err
	}
	
	return token, nil
}

func (s JwtTokenService) GenerateToken(details auth.UserDetails) (Token, error) {
	var token Token
	
	jwtToken := jwt.New(s.SigningMethod)
	
	duration := time.Hour * time.Duration(s.JwtConfig.Expire)
	
	jwtToken.Claims["exp"] = time.Now().Add(duration).Unix()
	jwtToken.Claims["iat"] = time.Now().Unix()
	jwtToken.Claims["jit"] = details.Username
	
	tokenString, err := jwtToken.SignedString([]byte(s.JwtConfig.Secret))
	
	if err != nil {
		return token, err
	}
	
	return NewToken(tokenString), nil
}