package koala

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionConfig represents the session settings
type SessionConfig struct {
	Secret string
}

// NewSessionConfig creates an instance for SessionConfig
func NewSessionConfig(s string) SessionConfig {
	return SessionConfig{s}
}

// NewCookieStore creates an instance *sessions.CookieStore
// This store is used to add the token into the cookie
func NewCookieStore(c SessionConfig) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(c.Secret))
}

const authCookie = "sid"

// AuthTokenStore is an auth token storage
// It accesses the session cookie and wrap the needed logic to save, get and clear the storage
type AuthTokenStore struct {
	Store sessions.Store
}

// NewAuthTokenStore creates a new instance of AuthTokenStore
func NewAuthTokenStore(s sessions.Store) AuthTokenStore {
	return AuthTokenStore{s}
}

// Save saves the token into session cookie
func (j AuthTokenStore) Save(r *http.Request, w http.ResponseWriter, token Token) error {
	s, err := j.Store.Get(r, authCookie)

	if err != nil {
		return err
	}

	s.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	s.Values["token"] = token.Value

	return s.Save(r, w)
}

// Get gets the token from session cookie
func (j AuthTokenStore) Get(r *http.Request) (Token, bool) {
	var token Token

	s, err := j.Store.Get(r, authCookie)

	if err != nil {
		return token, false
	}

	tokenValue, ok := s.Values["token"].(string)

	token.Value = tokenValue

	return token, ok
}

// Clear cleans the token from session cookie
func (j AuthTokenStore) Clear(r *http.Request, w http.ResponseWriter) error {
	s, _ := j.Store.Get(r, "sid")

	delete(s.Values, "token")

	s.Options.MaxAge = -1

	return s.Save(r, w)
}

var m map[string]interface{}

func init() {
	gob.Register(&m)
}
