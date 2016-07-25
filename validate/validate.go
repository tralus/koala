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

func IsValidationContextError(err error) bool {
	_, ok := err.(ContextValidationError)
	return ok
}

type ContextValidationError struct {
	Errors []PropertyError
}

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

type PropertyError struct {
	Property string
	Message string
}

func IsPropertyError(err error) bool {
	_, ok := err.(PropertyError)
	return ok
}

func (v PropertyError) Error() string {
	return v.Message
}

type Context struct {
	Errors []PropertyError
}

func NewContext() Context {
	return Context{}
}

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

func NewPropertyError(property string, msg string) PropertyError {
	return PropertyError{property, msg}
}

type ChoicesString map[string]string

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

type ChoicesInt map[string]int

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

func IsEmail(property string, v string) error {
	if !emailRegex.MatchString(v) {
		return NewPropertyError(property, fmt.Sprintf(IS_EMAIL_MSG, property))
	}
	return nil
}

func IsJson(property string, v string) error {
	var js map[string]interface{}
	
	if err := json.Unmarshal([]byte(v), &js); err != nil {
		return NewPropertyError(property, fmt.Sprintf(IS_JSON_MSG, v)) 
	}
	
	return nil
}

func MinStrLength(property string, v string, min int) error {
	if (utf8.RuneCountInString(v) < min) {
		return NewPropertyError(property, fmt.Sprintf(MIN_LENGTH_MSG, property, min))
	}
	
	return nil
}