package validate

import (
	"reflect"

	"github.com/tralus/koala/errors"
)

func newArgumentError(err error) error {
	return errors.NewIllegalArgumentError(err)
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
		return newArgumentError(errors.New("Unsupported type."))
	}

	if !valid {
		return newArgumentError(errors.New("Non zero value required."))
	}

	return nil
}
