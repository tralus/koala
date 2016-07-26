package validate

import (
	"reflect"
	"errors"
	"strings"
	"sort"
	"fmt"
	"strconv"
	"encoding/json"
	"unicode/utf8"
)

var ErrUnsupported = errors.New("Unsupported type")

const NOT_ZERO_MSG = "The '%s' can not be zero"

const IS_EMAIL_MSG = "The '%s' value is not a valid email"

const MIN_LENGTH_MSG = "The '%s' length value is lesser than %d"

const EQ_CHOICE_MSG = "The '%s' value is not equal to %s"

const IS_JSON_MSG = "The '%s' value is not a valid json"

// Generic error for context validation aspects
type ContextValidationError struct {
	Errors []PropertyError
}

// Built-in error interface
func (v ContextValidationError) Error() string {
	var message string
	
	for _, propertyError := range v.Errors {
		if (len(v.Errors) > 1) {
			message = " :: "
		}
		message += propertyError.Error()
	}
	
	return message
}

// It verifies if erro is a ContextValidationError type
func IsValidationContextError(err error) bool {
	_, ok := err.(ContextValidationError)
	return ok
}

// Generic error type for property validations
type PropertyError struct {
	Property string
	Msg string
}

// It verifies if error ir a PropertyError type
func IsPropertyError(err error) bool {
	_, ok := err.(PropertyError)
	return ok
}

// Built-in error interface
func (v PropertyError) Error() string {
	return v.Msg
}

// Context groups Property Errors
type Context struct {
	Errors []PropertyError
}

// It creates a new Context instance 
func NewContext() Context {
	return Context{}
}

// It checks if there are erros on context
func (c Context) Check(errors ...error) error {
	for _, err := range errors {
		if (err == nil) {
			continue
		}
		
		c.Errors = append(c.Errors, err.(PropertyError))
	}
	
	if (len(c.Errors) >0) {
		return ContextValidationError{c.Errors}
	}
	
	return nil
}

// It creates a new PropertyError instance
func NewPropertyError(property string, msg string) PropertyError {
	return PropertyError{property, msg}
}

// ChoicesString represents field with two or more possibilities
type ChoicesString map[string]string

// It verifies if the value is equal to at least one choice list element
func EqChoiceString(property string, v string, choices ChoicesString) error {
	if (len(v) > 0) {
		found := false
		
		for _, c := range choices {
			if c == v {
		    	found = true
		    }
		}
		
		if (found) {
			return nil
		}
		
		var values []interface{}
	
		tmp := []string{}
		
		for k, c := range choices {
			tmp = append(tmp, k + ": " + c)
			
			sort.Strings(tmp)
		}
		
		values = append(values, property)
		
		if (len(tmp) > 0) {
			values = append(values, strings.Join(tmp, ", "))
		}
		
		return NewPropertyError(property, fmt.Sprintf(EQ_CHOICE_MSG, values...))
	}
	
	return nil
}

// ChoicesInt represents field with two or more possibilities
type ChoicesInt map[string]int

// It verifies if the value is equal to at least one choice list element
func EqChoiceInt(property string, v int, choices ChoicesInt) error {
	if (v > 0) {
		found := false
		
		for _, c := range choices {
			if c == v {
		    	found = true
		    }
		}
		
		if (found) {
			return nil
		}
		
		var values []interface{}
	
		tmp := []string{}
		
		for k, c := range choices {
			tmp = append(tmp, k + ": " + strconv.Itoa(c))
			
			sort.Strings(tmp)
		}
		
		values = append(values, property)
		
		if (len(tmp) > 0) {
			values = append(values, strings.Join(tmp, ", "))
		}
		
		return NewPropertyError(property, fmt.Sprintf(EQ_CHOICE_MSG, values...))
	}	
	
	return nil
}

// It verifies if the value is a zero value on Go 
func NotZero(property string, v interface{}) error {
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
		return ErrUnsupported
	}

	if !valid {
		return NewPropertyError(property, fmt.Sprintf(NOT_ZERO_MSG, property))
	}
	
	return nil
}

// It verifies if the value is a valid email
func IsEmail(property string, v string) error {
	if !emailRegex.MatchString(v) {
		return NewPropertyError(property, fmt.Sprintf(IS_EMAIL_MSG, property))
	}
	return nil
}

// It verifies if the value is an valid json
func IsJson(property string, v string) error {
	var js map[string]interface{}
	
	if err := json.Unmarshal([]byte(v), &js); err != nil {
		return NewPropertyError(property, fmt.Sprintf(IS_JSON_MSG, v)) 
	}
	
	return nil
}

// It verifies if the string length is lesser than min length allowed
func MinStrLength(property string, v string, min int) error {
	if (utf8.RuneCountInString(v) < min) {
		return NewPropertyError(property, fmt.Sprintf(MIN_LENGTH_MSG, property, min))
	}
	
	return nil
}