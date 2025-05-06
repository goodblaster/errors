package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Format(t *testing.T) {
	err := New("test %v")
	assert.NotNil(t, err)

	formatted := err.Format(2)
	assert.Equal(t, "test 2", formatted.Error())
}

func TestError_FormatIs(t *testing.T) {
	UnformattedError := New("test %v")
	assert.NotNil(t, UnformattedError)

	formatted := UnformattedError.Format(2)
	assert.Equal(t, "test 2", formatted.Error())

	// Must use your package's Is function (which correctly unwraps Formatted)
	assert.True(t, Is(formatted, UnformattedError), "errors.Is should match the original unformatted error")
}
