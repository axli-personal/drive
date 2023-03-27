package errors

import (
	"fmt"
)

type Error struct {
	code    string // Classification of error
	message string // Detailed information about error
	err     error  // Optional original error
}

func New(code string, message string, err error) *Error {
	return &Error{
		code:    code,
		message: message,
		err:     err,
	}
}

func (e *Error) Error() string {
	msg := fmt.Sprintf("%s: %s", e.code, e.message)

	if e.err != nil {
		msg = fmt.Sprintf("%s\ncaused by: %s", msg, e.err.Error())
	}

	return msg
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Code() string {
	return e.code
}

func (e *Error) Message() string {
	return e.message
}
