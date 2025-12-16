package errors

import "fmt"

// FormattedError represents a formatted instance of an Error.
type FormattedError struct {
	parent *Error
	msg    string
}

// Error implements error for FormattedError.
func (f *FormattedError) Error() string {
	return f.msg
}

// Unwrap allows FormattedError to unwrap back to the parent *Error.
func (f *FormattedError) Unwrap() error {
	return f.parent
}

// Format applies fmt.Sprintf to the current error's Error() string and returns a formatted instance.
//
// The error message should contain valid fmt.Sprintf format verbs (e.g., %s, %d, %v).
// Example:
//
//	err := New("failed to process item %d: %s")
//	formatted := err.Format(42, "invalid input")
//	// formatted.Error() returns: "failed to process item 42: invalid input"
//
// IMPORTANT: While fmt.Sprintf doesn't panic, it will produce error indicators in the output for:
//   - Missing arguments: "error: %s %d" with Format("test") produces "error: test %!d(MISSING)"
//   - Extra arguments: "error: %s" with Format("test", "extra") produces "error: test%!(EXTRA ...)"
//   - Invalid/unknown verbs: "value: %z" with Format(42) produces "value: %!z(int=42)"
//   - Bare % characters: "100% complete" with Format() produces "100%!(NOVERB) complete"
//
// Best practice: Only call Format() on errors that were created with format templates.
// For non-template error messages, use the error directly without calling Format().
func (e *Error) Format(args ...any) error {
	return &FormattedError{
		parent: e,
		msg:    fmt.Sprintf(e.err.Error(), args...),
	}
}
