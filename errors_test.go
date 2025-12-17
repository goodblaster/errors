package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("test")
	assert.NotNil(t, err)
	assert.Equal(t, `test`, err.Error())

	err = Newf("test %v", 2)
	assert.NotNil(t, err)
	assert.Equal(t, `test 2`, err.Error())
}

func TestIs(t *testing.T) {
	err := errors.New("test")
	assert.True(t, Is(err, err))

	err2 := Wrap(err, "wrap")
	assert.True(t, Is(err2, err))

	err3 := Wrap(err2, "wrap")
	assert.True(t, Is(err3, err2))
	assert.True(t, Is(err3, err))
}

type CustomError struct {
	msg string
}

func (e *CustomError) Error() string {
	return e.msg
}

func TestAs(t *testing.T) {
	// Create a custom error
	customErr := &CustomError{msg: "this is a custom error"}
	wrappedErr := New("some other error")

	// Wrap the custom error
	err := Join(customErr, wrappedErr)

	var target *CustomError
	if !As(err, &target) {
		t.Fatalf("errors.As failed to find CustomError in wrapped error")
	}

	if target != customErr {
		t.Fatalf("errors.As did not return the expected CustomError instance")
	}
}

func TestMarshalJSON(t *testing.T) {
	t.Run("single error", func(t *testing.T) {
		err := New("simple error")
		data, marshalErr := err.MarshalJSON()
		assert.NoError(t, marshalErr)
		assert.Equal(t, `["simple error"]`, string(data))
	})

	t.Run("wrapped errors", func(t *testing.T) {
		err1 := errors.New("base error")
		err2 := Wrap(err1, "context")
		data, marshalErr := err2.MarshalJSON()
		assert.NoError(t, marshalErr)
		// errors.Join separates errors with newlines
		assert.Contains(t, string(data), "context")
		assert.Contains(t, string(data), "base error")
	})

	t.Run("error with natural newlines", func(t *testing.T) {
		err := New("line1\nline2\nline3")
		data, marshalErr := err.MarshalJSON()
		assert.NoError(t, marshalErr)
		assert.Equal(t, `["line1","line2","line3"]`, string(data))
	})

	t.Run("empty error message", func(t *testing.T) {
		err := New("")
		data, marshalErr := err.MarshalJSON()
		assert.NoError(t, marshalErr)
		assert.Equal(t, `[""]`, string(data))
	})
}

func TestUnwrap(t *testing.T) {
	t.Run("unwrap Error returns wrapped error", func(t *testing.T) {
		inner := errors.New("inner")
		outer := Wrap(inner, "outer")

		unwrapped := outer.Unwrap()
		assert.NotNil(t, unwrapped)
		// Unwrap should return the joined error
		assert.True(t, errors.Is(unwrapped, inner))
	})

	t.Run("unwrap function with Error", func(t *testing.T) {
		inner := errors.New("inner")
		outer := Wrap(inner, "outer")

		unwrapped := Unwrap(outer)
		assert.NotNil(t, unwrapped)
		assert.True(t, errors.Is(unwrapped, inner))
	})

	t.Run("unwrap nil returns nil", func(t *testing.T) {
		unwrapped := Unwrap(nil)
		assert.Nil(t, unwrapped)
	})
}

func TestStdlibCompatibility(t *testing.T) {
	t.Run("stdlib errors.Is works with Error", func(t *testing.T) {
		base := errors.New("base error")
		wrapped := Wrap(base, "wrapped")

		// stdlib errors.Is should work because *Error implements Unwrap()
		assert.True(t, errors.Is(wrapped, base))
	})

	t.Run("stdlib errors.As works with Error", func(t *testing.T) {
		customErr := &CustomError{msg: "custom"}
		wrapped := Wrap(customErr, "context")

		var target *CustomError
		// stdlib errors.As should work because *Error implements Unwrap()
		assert.True(t, errors.As(wrapped, &target))
		assert.Equal(t, customErr, target)
	})

	t.Run("stdlib works with deeply nested errors", func(t *testing.T) {
		base := errors.New("base")
		err1 := Wrap(base, "level1")
		err2 := Wrap(err1, "level2")
		err3 := Wrap(err2, "level3")

		assert.True(t, errors.Is(err3, base))
		assert.True(t, errors.Is(err3, err1))
		assert.True(t, errors.Is(err3, err2))
	})
}

func TestErrorChainPreservation(t *testing.T) {
	t.Run("wrapping Error preserves type in chain", func(t *testing.T) {
		err1 := New("error1")
		err2 := Wrap(err1, "error2")
		err3 := Wrap(err2, "error3")

		// All *Error types should be findable in the chain
		var target1 *Error
		assert.True(t, errors.As(err3, &target1))

		// Should be able to find intermediate *Error types
		assert.True(t, Is(err3, err1))
		assert.True(t, Is(err3, err2))
	})

	t.Run("mixed error types in chain", func(t *testing.T) {
		stdErr := errors.New("std error")
		err1 := Wrap(stdErr, "wrapped once")
		customErr := &CustomError{msg: "custom"}
		err2 := Wrap(customErr, "wrapped custom")
		joined := Join(err1, err2)

		// Should find both standard and custom errors
		assert.True(t, errors.Is(joined, stdErr))

		var custom *CustomError
		assert.True(t, errors.As(joined, &custom))
	})
}

func TestJoin(t *testing.T) {
	t.Run("join multiple errors", func(t *testing.T) {
		err1 := errors.New("error1")
		err2 := errors.New("error2")
		err3 := errors.New("error3")

		joined := Join(err1, err2, err3)

		// All errors should be findable
		assert.True(t, Is(joined, err1))
		assert.True(t, Is(joined, err2))
		assert.True(t, Is(joined, err3))
	})

	t.Run("join with nil errors", func(t *testing.T) {
		err1 := errors.New("error1")
		joined := Join(err1, nil, nil)

		assert.True(t, Is(joined, err1))
	})

	t.Run("join returns Error type", func(t *testing.T) {
		err1 := errors.New("error1")
		err2 := errors.New("error2")

		joined := Join(err1, err2)

		// Should return *Error, not just error
		_, ok := joined.(*Error)
		assert.True(t, ok, "Join should return *Error type")
	})
}

func TestWrapNil(t *testing.T) {
	t.Run("Wrap nil creates error", func(t *testing.T) {
		err := Wrap(nil, "context")
		assert.NotNil(t, err)
		assert.Equal(t, "context", err.Error())
	})

	t.Run("Wrapf nil creates error", func(t *testing.T) {
		err := Wrapf(nil, "context: %d", 42)
		assert.NotNil(t, err)
		assert.Equal(t, "context: 42", err.Error())
	})
}

func TestDeepErrorChain(t *testing.T) {
	// Test very deep error chains don't cause stack overflow
	base := errors.New("base")
	current := Wrap(base, "level0")

	// Create a chain of 1000 errors
	for i := 1; i < 1000; i++ {
		current = Wrap(current, "level")
	}

	// Should still be able to find base error
	assert.True(t, Is(current, base))

	// Should not panic on As
	var target *Error
	assert.True(t, As(current, &target))
}
