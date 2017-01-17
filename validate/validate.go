package validate

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/tralus/koala/errors"
)

// ContextValidationError represents the context errors
type ContextValidationError struct {
	Errors []ArgumentError
}

// Built-in error interface
func (v ContextValidationError) Error() string {
	return v.Error()
}

// IsValidationContextError verifies if error is a ContextValidationError type
func IsValidationContextError(err error) bool {
	_, ok := err.(ContextValidationError)
	return ok
}

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

// Context groups errors
type Context struct{}

// NewContext creates a Context instance
func NewContext() Context {
	return Context{}
}

// Do checks if the context has errors
func (c Context) Do(errors ...error) error {
	var argErrors []ArgumentError

	for _, err := range errors {
		if e, ok := err.(ArgumentError); ok {
			argErrors = append(argErrors, e)
		}
	}

	if len(argErrors) > 0 {
		return ContextValidationError{argErrors}
	}

	return nil
}

// ChoicesInt represents an int choice field
type ChoicesInt []int

// EqChoiceInt verifies if v is equal to at least one option in the list
// If the value should not be zero, use the not zero validator
func EqChoiceInt(v int, choices ChoicesInt) error {
	if v == 0 || len(choices) == 0 {
		return nil
	}

	d := []string{}
	found := false

	sort.Ints(choices)
	for _, c := range choices {
		if c == v {
			found = true
		}

		d = append(d, strconv.Itoa(c))
	}

	if found {
		return nil
	}

	s := fmt.Sprintf("option(int:%s)", strings.Join(d, ","))

	return NewArgumentError(notValidateMsg(v, s))
}

// ChoicesString represents an str choice field
type ChoicesString []string

// EqChoiceString verifies if v is equal to at least one option in the list
// If the value should not be zero, use the not zero validator
func EqChoiceString(v string, choices ChoicesString) error {
	if v == "" || len(choices) == 0 {
		return nil
	}

	d := []string{}
	found := false

	sort.Strings(choices)
	for _, c := range choices {
		if c == v {
			found = true
		}

		d = append(d, c)
	}

	if found {
		return nil
	}

	s := fmt.Sprintf("option(str:%s)", strings.Join(d, ","))

	return NewArgumentError(notValidateMsg(v, s))
}

// NotZeroString verifies if s is not empty
func NotZeroString(s string) error {
	valid := utf8.RuneCountInString(s) != 0

	if !valid {
		return NewArgumentError("Non zero string required.")
	}

	return nil
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

func notValidateMsg(v interface{}, validatorName string) string {
	return fmt.Sprintf("%s does not validate as %s.", fmt.Sprint(v), validatorName)
}

// IsEmail verifies if the v is a valid email
func IsEmail(v string) error {
	if ok := govalidator.IsEmail(v); !ok {
		return NewArgumentError(notValidateMsg(v, "email"))
	}
	return nil
}

// IsJSON verifies if the v is an valid json
func IsJSON(v string) error {
	if ok := govalidator.IsJSON(v); !ok {
		return NewArgumentError(notValidateMsg(v, "json"))
	}
	return nil
}

// MinStrLength verifies if v length is lesser than the min length allowed
func MinStrLength(v string, m int) error {
	if utf8.RuneCountInString(v) < m {
		return NewArgumentError(notValidateMsg(v, fmt.Sprintf("length(min:%d)", m)))
	}
	return nil
}
