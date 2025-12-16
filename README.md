# go-errors

A Go error handling library that extends the standard library's error capabilities with additional features.

## Features

- **Error wrapping** - Wrap errors with context using `Wrap()` and `Wrapf()`
- **JSON marshalling** - Errors serialize as string arrays for easy API responses
- **Error formatting** - Apply template variables to error messages with `Format()`
- **Typed nil detection** - `IsNil()` detects typed nil errors that slip through normal nil checks
- **Full stdlib compatibility** - Works seamlessly with `errors.Is` and `errors.As`
- **Type preservation** - Error chains maintain full type information for accurate error matching

## Installation

```bash
go get github.com/goodblaster/errors
```

## Basic Usage

### Creating Errors

```go
import "github.com/goodblaster/errors"

// Create a simple error
err := errors.New("something went wrong")

// Create an error with formatting
err := errors.Newf("failed to process %d items", count)
```

### Wrapping Errors

```go
// Wrap an error with additional context
if err := doSomething(); err != nil {
    return errors.Wrap(err, "failed to do something")
}

// Wrap with formatted context
if err := processItem(id); err != nil {
    return errors.Wrapf(err, "failed to process item %d", id)
}
```

### Error Formatting

Create error templates and format them later:

```go
// Define error template
var ErrInvalidInput = errors.New("invalid input: expected %s, got %s")

// Format with specific values
err := ErrInvalidInput.Format("number", "string")
// Output: "invalid input: expected number, got string"
```

**Note:** Only use `Format()` on errors created as templates. Calling `Format()` on regular error messages containing `%` will produce error indicators in the output.

### Checking Errors

```go
// Use package functions (handles FormattedError properly)
if errors.Is(err, ErrNotFound) {
    // handle not found
}

var customErr *MyCustomError
if errors.As(err, &customErr) {
    // handle custom error
}

// Or use stdlib directly (also works!)
if errors.Is(err, ErrNotFound) {
    // handle not found
}
```

### JSON Marshalling

```go
type Response struct {
    Error *errors.Error `json:"error,omitempty"`
}

err := errors.Wrap(originalErr, "operation failed")
resp := Response{Error: err}
json.Marshal(resp)
// Output: {"error": ["operation failed", "original error message"]}
```

### Typed Nil Detection

Detect the common Go pitfall where typed nil pointers are not equal to nil:

```go
func mightReturnTypedNil() error {
    var err *MyError = nil
    return err // This is NOT nil when compared to error(nil)!
}

err := mightReturnTypedNil()
if err != nil {
    // This block executes even though err is logically nil!
}

// Use IsNil instead
if !errors.IsNil(err) {
    // This correctly identifies the typed nil
}
```

## API Reference

### Core Functions

- `New(msg string) *Error` - Create a new error
- `Newf(msg string, args ...any) *Error` - Create a new formatted error
- `Wrap(err error, msg string) *Error` - Wrap an error with context
- `Wrapf(err error, msg string, args ...any) *Error` - Wrap with formatted context
- `Join(errs ...error) error` - Join multiple errors into one
- `Is(err, target error) bool` - Check if err matches target
- `As(err error, target any) bool` - Find first error of specific type
- `Unwrap(err error) error` - Unwrap one level of error
- `IsNil(err error) bool` - Check if error is truly nil (even typed nils)

### Methods

- `(*Error) Error() string` - Returns the error message
- `(*Error) Unwrap() error` - Returns the wrapped error
- `(*Error) Format(args ...any) error` - Format error message with arguments
- `(*Error) MarshalJSON() ([]byte, error)` - Serialize error to JSON

## Design Notes

- **Uses `errors.Join` internally** - Despite being named "Wrap", uses stdlib Join under the hood
- **Type preservation** - Wrapping an `*Error` preserves it in the chain, unlike some libraries
- **Nil wrapping** - `Wrap(nil, "msg")` returns a non-nil error with the message (for context preservation)
- **Stdlib compatible** - `*Error` implements `Unwrap()`, so stdlib `errors.Is/As` work directly
- **Encapsulated** - The wrapped error field is unexported to prevent direct modification

## Testing

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v -run TestName
```

## License

See LICENSE file for details.