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

// Session represents a session
type Session struct {
	name    string
	store   sessions.Store
	options *sessions.Options
}

// New creates a Session instance
func New(name string, s sessions.Store, options *sessions.Options) Session {
	return Session{name, s, options}
}

// Values represents the values of the session
type Values map[interface{}]interface{}

// Save saves the token into session cookie
func (j Session) Save(r *http.Request, w http.ResponseWriter, values Values) error {
	s, err := j.store.Get(r, j.name)

	if err != nil {
		return err
	}

	s.Values = values

	return s.Save(r, w)
}

// Get gets the session values
func (j Session) Get(r *http.Request) (Values, error) {
	s, err := j.store.Get(r, j.name)

	if err != nil {
		return nil, err
	}

	return s.Values, nil
}

// Start starts the session with options
func (j Session) Start(r *http.Request) error {
	s, err := j.store.Get(r, j.name)

	s.Options = j.options

	if err != nil {
		return err
	}

	return nil
}

// Clear clears the token from session cookie
func (j Session) Clear(r *http.Request, w http.ResponseWriter) error {
	s, _ := j.store.Get(r, j.name)

	s.Values = nil

	s.Options.MaxAge = -1

	return s.Save(r, w)
}

var m map[string]interface{}

func init() {
	gob.Register(&m)
}
