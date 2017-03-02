package context

import (
	"net/http"

	"fmt"

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
		errMsg := fmt.Sprintf(
			"Key %s is not in the context.", key,
		)
		return nil, errors.NewIllegalStateError(errMsg)
	}

	return u, nil
}
