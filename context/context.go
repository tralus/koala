package context

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/tralus/koala/errors"
)

// Add adds a value to the context
func Add(r *http.Request, key, val interface{}) {
	context.Set(r, key, val)
}

// Get gets a value from the context
func Get(r *http.Request, key interface{}) (interface{}, error) {
	u := context.Get(r, key)

	if u == nil {
		err := errors.Errorf("Key %s is not in the context.", key)
		return nil, errors.NewIllegalStateError(err)
	}

	return u, nil
}
