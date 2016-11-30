package session

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

// Config represents the session config
type Config struct {
	Secret string
}

// NewConfig creates a Config instance
func NewConfig(s string) Config {
	return Config{s}
}

// NewCookieStore creates a *sessions.CookieStore instance
// The store is used to add values in the cookie
func NewCookieStore(c Config) *sessions.CookieStore {
	return sessions.NewCookieStore([]byte(c.Secret))
}

const sessionID = "koala_ssid"

// Store represents a session store
type Store struct {
	store    sessions.Store
	request  *http.Request
	response http.ResponseWriter
}

// NewStore creates a new instance of Store
func NewStore(s sessions.Store, r *http.Request, w http.ResponseWriter) Store {
	store := Store{s, r, w}

	store.SetOptions(&sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	})

	return store
}

// SetOptions sets the session options
func (j Store) SetOptions(options *sessions.Options) error {
	s, err := j.store.Get(j.request, sessionID)

	if err != nil {
		return err
	}

	s.Options = options

	return s.Save(j.request, j.response)
}

// Values represents the values of the session
type Values map[interface{}]interface{}

// Save saves the token into session cookie
func (j Store) Save(values Values) error {
	s, err := j.store.Get(j.request, sessionID)

	if err != nil {
		return err
	}

	s.Values = values

	return s.Save(j.request, j.response)
}

// Get gets the token from session cookie
func (j Store) Get() (Values, error) {
	s, err := j.store.Get(j.request, sessionID)

	if err != nil {
		return nil, err
	}

	return s.Values, nil
}

// Clear cleans the token from session cookie
func (j Store) Clear() error {
	s, _ := j.store.Get(j.request, sessionID)

	s.Values = nil

	s.Options.MaxAge = -1

	return s.Save(j.request, j.response)
}

var m map[string]interface{}

func init() {
	gob.Register(&m)
}
