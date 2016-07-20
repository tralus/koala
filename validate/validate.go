package validate

import (
	"strings"
	"sort"
	"fmt"
	"strconv"
	"encoding/json"
	
	"gopkg.in/go-playground/validator.v8"
)

const NOT_EMPTY_MSG = "The '%s' can not be empty"

const IS_EMAIL_MSG = "The '%s' value is not a valid email"

const MIN_LENGTH_MSG = "The '%s' length value is lesser than %d"

const EQ_CHOICE_MSG = "The '%s' value is not equal to %s"

const IS_JSON_MSG = "The '%s' value is not a valid json"

var validate *validator.Validate

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

func Property(property string, v interface{}, tags string, f string, e ...interface{}) error {
	err := validate.Field(v, tags)
	
	var values []interface{}
	
	values = append(values, property)
	
	message := fmt.Sprintf(f, append(values, e...)...)
	
	if (err != nil) {
		return NewPropertyError(property, message)
	}
	
	return nil
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

func NotEmpty(property string, v interface{}) error {
	return Property(property, v, "required", NOT_EMPTY_MSG)
}

func IsEmail(property string, v string) error {
	if (len(v) > 0) {
		return Property(property, v, "email", IS_EMAIL_MSG)
	}
	
	return nil
}

func IsJson(property string, v string) error {
	var js map[string]interface{}
	
	if err := json.Unmarshal([]byte(v), &js); err != nil {
		var values []interface{}
	
		values = append(values, property)
		
		message := fmt.Sprintf(IS_JSON_MSG, values...)
		
		return NewPropertyError(property, message) 
	}
	
	return nil
}

func MinLength(property string, v string, min int) error {
	if (len(v) > 0) {
		tags := "min=" + strconv.Itoa(min)
		return Property(property, v, tags, MIN_LENGTH_MSG, min)
	}
	
	return nil
}

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}