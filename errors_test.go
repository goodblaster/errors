package errors

import (
	oerrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("test")
	assert.NotNil(t, err)
	assert.Equal(t, `["test"]`, err.Error())

	err = New("test %v", 2)
	assert.NotNil(t, err)
	assert.Equal(t, `["test 2"]`, err.Error())
}

func TestIs(t *testing.T) {
	err := oerrors.New("test")
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
