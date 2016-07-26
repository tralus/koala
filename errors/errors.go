// NOTE: This package intentionally mirrors the standard "errors" module.
// All packages on Koala should use this 
package errors

import (
	"bytes"
	"runtime"
) 

// RuntimeError interface exposes additional information about the error.
type RuntimeError interface {
	// This returns the stack trace without the error message.
	GetStack() string
	
	// Implements the built-in error interface.
	Error() string
}

// It verifies if error is a RuntimeError type 
func IsRuntimeError(e error) bool {
	_, ok := e.(RuntimeError)
	return ok
}

// Base error type
// It implements RuntimeError interface
type Error struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e Error) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e Error) Error() string {
	return e.Msg
}

// It creates a new Error instance
func New(m string) error {
	s, _ := StackTrace()
	return Error{m, s} 
}

// It verifies if error is an Error type
func IsError(e error) bool {
	_, ok := e.(Error)
	return ok
}

// Generic error type for not found logic
type NotFoundError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e NotFoundError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e NotFoundError) Error() string {
	return e.Msg
}

// It creates a new NotFoundError instance
func NewNotFoundError(m string) error {
	s, _ := StackTrace()
	return NotFoundError{m, s} 
}

// It verifies if error is a NotFoundError type
func IsNotFoundError(e error) bool {
	_, ok := e.(NotFoundError)
	return ok
}

// Generic error type for illegal data 
type IllegalStateError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e IllegalStateError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e IllegalStateError) Error() string {
	return e.Msg
}

// It creates a new NotFoundError instance
func NewIllegalStateError(m string) error {
	s, _ := StackTrace()
	return IllegalStateError{m, s} 
}

// It verifies if error is an IllegalStateError type
func IsIllegalStateError(e error) bool {
	_, ok := e.(IllegalStateError)
	return ok
}

// Generic error type for relationship aspects
type RelationshipError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e RelationshipError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e RelationshipError) Error() string {
	return e.Msg
}

// It creates a new RelationshipError instance
func NewRelationshipError(m string) error {
	s, _ := StackTrace()
	return RelationshipError{m, s} 
}

// It verifies if error is an RelationshipError type
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

// This returns the current stack trace string.
func StackTrace() (current, context string) {
	return stackTrace(3)
}