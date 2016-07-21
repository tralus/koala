package errors

import (
	"bytes"
	"runtime"
) 

type RuntimeError interface {
	GetStack() string
	Error() string
}

func IsRuntimeError(e error) bool {
	_, ok := e.(RuntimeError)
	return ok
}

type Error struct {
	Msg string
	Stack string
}

func (e Error) GetStack() string {
	return e.Stack
}

func (e Error) Error() string {
	return e.Msg
}

func New(m string) error {
	s, _ := StackTrace()
	return Error{m, s} 
}

func IsError(e error) bool {
	_, ok := e.(Error)
	return ok
}

type NotFoundError struct {
	Msg string
	Stack string
}

func (e NotFoundError) GetStack() string {
	return e.Stack
}

func (e NotFoundError) Error() string {
	return e.Msg
}

func NewNotFoundError(m string) error {
	s, _ := StackTrace()
	return NotFoundError{m, s} 
}

func IsNotFoundError(e error) bool {
	_, ok := e.(NotFoundError)
	return ok
}

type IllegalStateError struct {
	Msg string
	Stack string
}

func (e IllegalStateError) GetStack() string {
	return e.Stack
}

func (e IllegalStateError) Error() string {
	return e.Msg
}

func NewIllegalStateError(m string) error {
	s, _ := StackTrace()
	return IllegalStateError{m, s} 
}

func IsIllegalStateError(e error) bool {
	_, ok := e.(IllegalStateError)
	return ok
}

type RelationshipError struct {
	Msg string
	Stack string
}

func (e RelationshipError) GetStack() string {
	return e.Stack
}

func (e RelationshipError) Error() string {
	return e.Msg
}

func NewRelationshipError(m string) error {
	s, _ := StackTrace()
	return RelationshipError{m, s} 
}

func IsRelationshipError(e error) bool {
	_, ok := e.(RelationshipError)
	return ok
}

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

func StackTrace() (current, context string) {
	return stackTrace(3)
}