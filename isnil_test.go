package errors

import (
	"errors"
	"fmt"
	"testing"
)

type ptrError struct {
	msg string
}

func (e *ptrError) Error() string {
	if e == nil {
		return "<nil ptrError>"
	}
	return e.msg
}

type valError struct {
	msg string
}

func (e valError) Error() string {
	return e.msg
}

// alias type just to make sure weird types don't break anything
type stringError string

func (e stringError) Error() string {
	return string(e)
}

func TestIsNil(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil interface",
			err:  nil,
			want: true,
		},
		{
			name: "typed nil pointer in interface",
			err: func() error {
				var e *ptrError = nil
				return e
			}(),
			want: true,
		},
		{
			name: "non-nil pointer error",
			err:  &ptrError{msg: "boom"},
			want: false,
		},
		{
			name: "non-nil value error",
			err:  valError{msg: "boom"},
			want: false,
		},
		{
			name: "standard library errors.New",
			err:  errors.New("std error"),
			want: false,
		},
		{
			name: "fmt.Errorf simple",
			err:  fmt.Errorf("fmt error"),
			want: false,
		},
		{
			name: "fmt.Errorf with typed nil wrapped",
			err: func() error {
				var e *ptrError = nil
				return fmt.Errorf("wrapped: %w", e)
			}(),
			// The wrapper itself is a real non-nil error value,
			// so we expect IsNil to be false.
			want: false,
		},
		{
			name: "fmt.Errorf with nil wrapped",
			err: func() error {
				// This is effectively the same as fmt.Errorf("x: %w", nil)
				// which produces a non-nil error.
				return fmt.Errorf("wrapped: %w", (error)(nil))
			}(),
			want: false,
		},
		{
			name: "string alias error type",
			err:  stringError("string error"),
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			got := IsNil(tt.err)
			if got != tt.want {
				t.Fatalf("IsNil(%T(%v)) = %v, want %v",
					tt.err, tt.err, got, tt.want)
			}
		})
	}
}

// This test exists only to ensure IsNil does not panic on
// various dynamic types that implement error.
func TestIsNil_NoPanicOnWeirdTypes(t *testing.T) {
	type ifaceError interface {
		error
		Extra() string
	}
	type weird struct{}

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("IsNil panicked with weird type: %v", r)
			}
		}()

		var _ error = weirdError("x")
		var err error = weirdError("x")
		_ = IsNil(err)
	}()

	_ = ifaceError(nil) // just to ensure ifaceError compiles
}

type weirdError string

func (w weirdError) Error() string { return string(w) }

// Optional: explicit regression test for the classic bug pattern.
func TestIsNil_TypedNilRegression(t *testing.T) {
	// Simulate a "buggy" function that returns a typed-nil error.
	buggy := func() error {
		var e *ptrError = nil
		return e
	}

	err := buggy()
	if err != nil && IsNil(err) {
		// This is the key guarantee we want:
		// even though err != nil, our IsNil should treat it as nil.
	} else if !IsNil(err) {
		t.Fatalf("typed nil error should be treated as nil by IsNil")
	}
}
