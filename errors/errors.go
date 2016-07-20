package errors

import (
	"github.com/dropbox/godropbox/errors"
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
	s, _ := errors.StackTrace()
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
	s, _ := errors.StackTrace()
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
	s, _ := errors.StackTrace()
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
	s, _ := errors.StackTrace()
	return RelationshipError{m, s} 
}

func IsRelationshipError(e error) bool {
	_, ok := e.(RelationshipError)
	return ok
}

func StackTrace() string {
	stack, _ := errors.StackTrace()
	return stack
}