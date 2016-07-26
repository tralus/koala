package koala

import (
	"encoding/gob"
	"net/http"
	
	"github.com/gorilla/sessions"
)

// It creates a new session cookie store
// This session is used to add the token on cookie
func NewSessionCookie(c SessionConfig) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(c.Secret))
}

const AUTH_COOKIE = "sid"

// AuthTokenStore is an auth token storage
// It accesses the session cookie and wrap the needed logic to save, get and clear the storage
type AuthTokenStore struct {
	Store sessions.Store
}

// It creates a new auth token store
func NewAuthTokenStore(s sessions.Store) AuthTokenStore {
	return AuthTokenStore{s}
}

// It saves the token on session cookie
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

// It gets the token on session cookie
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

// It cleans the token on session cookie
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