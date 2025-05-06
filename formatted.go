package errors

import "fmt"

// Formatted represents a formatted instance of an Error.
type Formatted struct {
	parent *Error
	msg    string
}

// Error implements error for Formatted.
func (f *Formatted) Error() string {
	return f.msg
}

// Unwrap allows Formatted to unwrap back to the parent *Error.
func (f *Formatted) Unwrap() error {
	return f.parent
}

// Format applies fmt.Sprintf to the current error's Err.Error() and returns a formatted instance.
func (e Error) Format(args ...any) error {
	return &Formatted{
		parent: &e,
		msg:    fmt.Sprintf(e.Err.Error(), args...),
	}
}
