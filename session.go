package koala

import (
	"encoding/gob"
	"net/http"
	
	"github.com/gorilla/sessions"
)

func NewSessionCookie(c SessionConfig) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(c.Secret))
}

const AUTH_COOKIE = "sid"

type AuthTokenStore struct {
	Store sessions.Store
}

func NewAuthTokenStore(s sessions.Store) AuthTokenStore {
	return AuthTokenStore{s}
}

func (j AuthTokenStore) Save(r *http.Request, w http.ResponseWriter, token Token) error {
	s, err := j.Store.Get(r, AUTH_COOKIE)
	
	if (err != nil) {
		return err
	}
 
    s.Options = &sessions.Options{
    	Path: "/",
    	MaxAge: 86400 * 7,
    	HttpOnly: true,
	}
    
    s.Values["token"] = token.Value
    
	return s.Save(r, w)
}

func (j AuthTokenStore) Get(r *http.Request) (Token, bool) {
	var token Token
	
	s, err := j.Store.Get(r, AUTH_COOKIE)
	
	if (err != nil) {
		return token, false
	}
    
    tokenValue, ok := s.Values["token"].(string)
    
    token.Value = tokenValue
    
    return token, ok
}

func (j AuthTokenStore) Clear(r *http.Request, w http.ResponseWriter) error {
	s, _ := j.Store.Get(r, "sid")
    
    delete(s.Values, "token")
    
    s.Options.MaxAge = -1
    
    return s.Save(r, w)
}

type M map[string]interface{}

func init() {
    gob.Register(&M{})
}