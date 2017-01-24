package validate

import (
	"reflect"

	"github.com/tralus/koala/errors"
)

// ArgumentError wraps errors.IllegalArgumentError
type ArgumentError struct {
	err error
}

// Error gets the error message
func (a ArgumentError) Error() string {
	return a.err.Error()
}

// GetStack gets the error stack trace
func (a ArgumentError) GetStack() string {
	if e, ok := a.err.(errors.IllegalArgumentError); ok {
		return e.GetStack()
	}
	return ""
}

// NewArgumentError creates an ArgumentError instance
func NewArgumentError(m string) ArgumentError {
	return ArgumentError{
		errors.NewIllegalArgumentError(m),
	}
}

// IsArgumentError verifies if error is a ArgumentError
func IsArgumentError(err error) bool {
	_, ok := err.(ArgumentError)
	return ok
}

// NotZero verifies if the v is a zero value on Go
func NotZero(v interface{}) error {
	st := reflect.ValueOf(v)
	valid := true

	switch st.Kind() {
	case reflect.String:
		valid = len(st.String()) != 0

	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()

	case reflect.Slice, reflect.Map, reflect.Array:
		valid = st.Len() != 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() != 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() != 0

	case reflect.Float32, reflect.Float64:
		valid = st.Float() != 0

	case reflect.Bool:
		valid = st.Bool()

	case reflect.Invalid:
		valid = false // always invalid

	case reflect.Struct:
		valid = true // always valid since only nil pointers are empty

	default:
		return NewArgumentError("Unsupported type.")
	}

	if !valid {
		return NewArgumentError("Non zero value required.")
	}

	return nil
}
