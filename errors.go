package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Error struct {
	Err error
}

func New(msg string, args ...any) *Error {
	return &Error{
		fmt.Errorf(msg, args...),
	}
}

func Wrap(err error, msg string, args ...any) *Error {
	if err == nil {
		return &Error{
			fmt.Errorf(msg, args...),
		}
	}

	if e, ok := err.(*Error); ok {
		err = e.Err
	}

	return &Error{
		errors.Join(fmt.Errorf(msg, args...), err),
	}
}

func Unwrap(err error) error {
	if e, ok := err.(*Error); ok {
		err = e.Err
	}
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return &Error{
		errors.Join(errs...),
	}
}

func Is(err, target error) bool {
	if e, ok := err.(*Error); ok {
		err = e.Err
	}

	if e, ok := target.(*Error); ok {
		target = e.Err
	}

	return errors.Is(err, target)
}

func As(err error, target any) bool {
	if e, ok := err.(*Error); ok {
		return As(e.Err, target)
	}
	return errors.As(err, target)
}

func (e Error) Error() string {
	return e.Err.Error()
}

func (e Error) MarshalJSON() ([]byte, error) {
	strs := strings.Split(e.Error(), "\n")
	return json.Marshal(strs)
}
