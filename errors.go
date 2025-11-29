package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

type Error struct {
	Err error
}

func New(msg string) *Error {
	return &Error{
		Err: errors.New(msg),
	}
}

func Newf(msg string, args ...any) *Error {
	return &Error{
		Err: fmt.Errorf(msg, args...),
	}
}

func Wrap(err error, msg string) *Error {
	if err == nil {
		return &Error{
			Err: fmt.Errorf(msg),
		}
	}

	if e, ok := err.(*Error); ok {
		err = e.Err
	}

	return &Error{
		Err: errors.Join(fmt.Errorf(msg), err),
	}
}

func Wrapf(err error, msg string, args ...any) *Error {
	if err == nil {
		return &Error{
			Err: fmt.Errorf(msg, args...),
		}
	}

	if e, ok := err.(*Error); ok {
		err = e.Err
	}

	return &Error{
		Err: errors.Join(fmt.Errorf(msg, args...), err),
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
		Err: errors.Join(errs...),
	}
}

func Is(err, target error) bool {
	// If the source error is formatted, unwrap to the parent
	if e := Unformatted(err); e != nil {
		err = e
	}

	if e, ok := err.(*Error); ok {
		err = e.Err
	}

	if e, ok := target.(*Error); ok {
		target = e.Err
	}

	return errors.Is(err, target)
}

// Unformatted - If this is formatted, return the unformatted parent error.
func Unformatted(err error) *Error {
	if e, ok := err.(*Formatted); ok {
		return e.parent
	}
	return nil
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

type iface struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

func IsNil(err error) bool {
	if err == nil {
		return true
	}

	i := *(*iface)(unsafe.Pointer(&err))
	return i.data == nil
}
