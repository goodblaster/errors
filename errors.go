// Package errors provides a wrapper around Go's standard error handling with additional features:
//   - JSON marshalling support (errors serialize as string arrays)
//   - Error formatting with template variables via Format()
//   - IsNil() function to detect typed nil errors
//   - Error wrapping using errors.Join internally
//
// The main Error type wraps standard Go errors and provides compatibility with
// errors.Is, errors.As, and errors.Unwrap while adding JSON serialization.
package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unsafe"
)

// Error wraps a standard Go error with additional functionality.
// The wrapped error is unexported to maintain encapsulation.
type Error struct {
	err error
}

func New(msg string) *Error {
	return &Error{
		err: errors.New(msg),
	}
}

func Newf(msg string, args ...any) *Error {
	return &Error{
		err: fmt.Errorf(msg, args...),
	}
}

// Wrap wraps an error with additional context.
// If err is nil, returns a new error with just the message.
// The returned error can be unwrapped to access the original error.
func Wrap(err error, msg string) *Error {
	if err == nil {
		return &Error{
			err: fmt.Errorf(msg),
		}
	}

	return &Error{
		err: errors.Join(fmt.Errorf(msg), err),
	}
}

// Wrapf wraps an error with additional formatted context.
// If err is nil, returns a new error with just the formatted message.
// The returned error can be unwrapped to access the original error.
func Wrapf(err error, msg string, args ...any) *Error {
	if err == nil {
		return &Error{
			err: fmt.Errorf(msg, args...),
		}
	}

	return &Error{
		err: errors.Join(fmt.Errorf(msg, args...), err),
	}
}

// Unwrap returns the result of calling the Unwrap method on err, if err's type contains
// an Unwrap method returning error. Otherwise, Unwrap returns nil.
// This is a convenience wrapper around errors.Unwrap.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Join(errs ...error) error {
	return &Error{
		err: errors.Join(errs...),
	}
}

// Is reports whether any error in err's chain matches target.
// It unwraps Formatted errors to their parent before checking.
// This function is compatible with errors.Is and can be used interchangeably.
func Is(err, target error) bool {
	// If the source error is formatted, unwrap to the parent
	if e := Unformatted(err); e != nil {
		err = e
	}

	return errors.Is(err, target)
}

// Unformatted returns the unformatted parent error if err is a FormattedError.
// Returns nil if err is not a FormattedError.
func Unformatted(err error) *Error {
	if e, ok := err.(*FormattedError); ok {
		return e.parent
	}
	return nil
}

// As finds the first error in err's chain that matches target.
// This function is compatible with errors.As and can be used interchangeably.
func As(err error, target any) bool {
	return errors.As(err, target)
}

func (e Error) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error, allowing errors.Is and errors.As to work correctly.
func (e *Error) Unwrap() error {
	return e.err
}

func (e Error) MarshalJSON() ([]byte, error) {
	strs := strings.Split(e.Error(), "\n")
	return json.Marshal(strs)
}

type iface struct {
	tab  unsafe.Pointer
	data unsafe.Pointer
}

// IsNil checks if an error is truly nil, even when wrapped in an interface.
//
// This function solves a common Go pitfall where a typed nil pointer,
// when returned as an error interface, is not equal to nil:
//
//	func buggy() error {
//	    var err *MyError = nil
//	    return err  // This is NOT nil when compared to error(nil)
//	}
//
// Using IsNil prevents bugs like:
//
//	if err := buggy(); err != nil {
//	    // This block executes even though the underlying error is nil!
//	}
//
// Instead, use:
//
//	if err := buggy(); !IsNil(err) {
//	    // This correctly identifies the typed nil
//	}
//
// IMPORTANT: This function uses unsafe.Pointer to inspect the error interface's
// internal structure. While this works with current Go implementations, it:
//   - May break in future Go versions if the interface representation changes
//   - Relies on internal implementation details
//   - Should be used sparingly
//
// The best practice is to fix code that returns typed nils rather than
// relying on IsNil. However, this function is useful for defensive programming
// when dealing with external libraries or legacy code.
func IsNil(err error) bool {
	if err == nil {
		return true
	}

	i := *(*iface)(unsafe.Pointer(&err))
	return i.data == nil
}
