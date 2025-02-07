package errors

import (
	"errors"
	"fmt"
)

func New(msg string, args ...any) error {
	return errors.New(fmt.Sprintf(msg, args...))
}

func Wrap(err error, msg string, args ...any) error {
	return errors.Join(err, New(msg, args...))
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
