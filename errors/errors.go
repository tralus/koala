// Package errors intentionally mirrors the standard "errors" module.
// All packages on Koala should use this.
package errors

import (
	"bytes"
	"runtime"
)

// RuntimeError interface exposes additional information about the error.
type RuntimeError interface {
	// Returns the stack trace without the error message.
	GetStack() string

	Error() string
}

// IsRuntimeError verifies if the error is a RuntimeError
func IsRuntimeError(e error) bool {
	_, ok := e.(RuntimeError)
	return ok
}

// Error represents a generic error
// It implements RuntimeError interface
type Error struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e Error) GetStack() string {
	return e.Stack
}

// Error gets the error message
func (e Error) Error() string {
	return e.Msg
}

// New creates an error from Error
func New(m string) error {
	s, _ := StackTrace()
	return Error{m, s}
}

// IsError verifies if the error is an Error
func IsError(e error) bool {
	_, ok := e.(Error)
	return ok
}

// NotFoundError represents a not found error
// It implements the RuntimeError interface
type NotFoundError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e NotFoundError) GetStack() string {
	return e.Stack
}

// Error gets the error message
func (e NotFoundError) Error() string {
	return e.Msg
}

// NewNotFoundError creates an error from NotFoundError
func NewNotFoundError(m string) error {
	s, _ := StackTrace()
	return NotFoundError{m, s}
}

// IsNotFoundError verifies if error is a NotFoundError
func IsNotFoundError(e error) bool {
	_, ok := e.(NotFoundError)
	return ok
}

// IllegalStateError represents an illegal state error
// It implements the RuntimeError interface
type IllegalStateError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e IllegalStateError) GetStack() string {
	return e.Stack
}

// Error gets the error message
func (e IllegalStateError) Error() string {
	return e.Msg
}

// NewIllegalStateError creates an IllegalStateError instance
func NewIllegalStateError(m string) error {
	s, _ := StackTrace()
	return IllegalStateError{m, s}
}

// IsIllegalStateError verifies if error is an IllegalStateError
func IsIllegalStateError(e error) bool {
	_, ok := e.(IllegalStateError)
	return ok
}

// IllegalArgumentError represents an illegal argument error
// It implements the RuntimeError interface
type IllegalArgumentError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e IllegalArgumentError) GetStack() string {
	return e.Stack
}

// Error gets the error message
func (e IllegalArgumentError) Error() string {
	return e.Msg
}

// NewIllegalArgumentError creates an IllegalArgumentError instance
func NewIllegalArgumentError(m string) error {
	s, _ := StackTrace()
	return IllegalArgumentError{m, s}
}

// IsIllegalStateError verifies if error is an IllegalStateError
func IsIllegalArgumentError(e error) bool {
	_, ok := e.(IllegalArgumentError)
	return ok
}

// RelationshipError represents a relationship error
type RelationshipError struct {
	Msg   string
	Stack string
}

// GetStack gets the error stack trace
func (e RelationshipError) GetStack() string {
	return e.Stack
}

// Error gets the error messsage
func (e RelationshipError) Error() string {
	return e.Msg
}

// NewRelationshipError creates an error from RelationshipError
func NewRelationshipError(m string) error {
	s, _ := StackTrace()
	return RelationshipError{m, s}
}

// IsRelationshipError verifies if error is an RelationshipError
func IsRelationshipError(e error) bool {
	_, ok := e.(RelationshipError)
	return ok
}

// Returns a copy of the error with the stack trace field populated and any
// other shared initialization; skips 'skip' levels of the stack trace.
//
// NOTE: This panics on any error.
func stackTrace(skip int) (current, context string) {
	buf := make([]byte, 128)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, len(buf)*2)
	}

	indexNewline := func(b []byte, start int) int {
		if start >= len(b) {
			return len(b)
		}
		searchBuf := b[start:]
		index := bytes.IndexByte(searchBuf, '\n')
		if index == -1 {
			return len(b)
		}
		return (start + index)
	}

	var strippedBuf bytes.Buffer
	index := indexNewline(buf, 0)
	if index != -1 {
		strippedBuf.Write(buf[:index])
	}

	for i := 0; i < skip; i++ {
		index = indexNewline(buf, index+1)
		index = indexNewline(buf, index+1)
	}

	isDone := false
	startIndex := index
	lastIndex := index
	for !isDone {
		index = indexNewline(buf, index+1)
		if (index - lastIndex) <= 1 {
			isDone = true
		} else {
			lastIndex = index
		}
	}
	strippedBuf.Write(buf[startIndex:index])
	return strippedBuf.String(), string(buf[index:])
}

// StackTrace returns the current stack trace string.
func StackTrace() (current, context string) {
	return stackTrace(3)
}

// CatchStackTrace catchs the error and get the stack trace
func CatchStackTrace(err error) string {
	if IsRuntimeError(err) {
		return err.(RuntimeError).GetStack()
	}
	return err.Error()
}
