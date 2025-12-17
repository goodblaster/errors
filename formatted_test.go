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

func TestError_FormatBehavior(t *testing.T) {
	t.Run("format with valid template", func(t *testing.T) {
		err := New("error: %s (code: %d)")
		formatted := err.Format("test", 42)
		assert.Equal(t, "error: test (code: 42)", formatted.Error())
	})

	t.Run("format with percent sign", func(t *testing.T) {
		// Note: fmt.Sprintf handles bare % by outputting %!(NOVERB)
		err := New("100% complete")
		formatted := err.Format()
		// fmt.Sprintf treats bare % as an error but doesn't panic
		assert.Contains(t, formatted.Error(), "%")
	})

	t.Run("format with missing arguments", func(t *testing.T) {
		err := New("error: %s %d")
		formatted := err.Format("test")
		// fmt.Sprintf handles missing args by outputting %!d(MISSING)
		assert.Contains(t, formatted.Error(), "test")
		assert.Contains(t, formatted.Error(), "MISSING")
	})

	t.Run("format with extra arguments", func(t *testing.T) {
		err := New("error: %s")
		formatted := err.Format("test", "extra", 42)
		// fmt.Sprintf outputs EXTRA indicators for unused arguments
		assert.Contains(t, formatted.Error(), "test")
		assert.Contains(t, formatted.Error(), "EXTRA")
	})
}

func TestError_FormatUnwrap(t *testing.T) {
	original := New("test %v")
	formatted := original.Format(42)

	// FormattedError should unwrap to the original Error
	unwrapped := Unwrap(formatted)
	assert.Equal(t, original, unwrapped)
}

func TestUnformatted(t *testing.T) {
	t.Run("formatted error returns parent", func(t *testing.T) {
		original := New("test %v")
		formatted := original.Format(42)

		parent := Unformatted(formatted)
		assert.NotNil(t, parent)
		assert.Equal(t, original, parent)
	})

	t.Run("non-formatted error returns nil", func(t *testing.T) {
		err := New("test")
		parent := Unformatted(err)
		assert.Nil(t, parent)
	})

	t.Run("nil error returns nil", func(t *testing.T) {
		parent := Unformatted(nil)
		assert.Nil(t, parent)
	})
}
