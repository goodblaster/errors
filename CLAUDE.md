# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go error handling library that extends Go's standard error capabilities with additional features:
- JSON marshalling support (errors serialize as string arrays)
- Error formatting with template variables via `Format()`
- `IsNil` function to handle typed nil errors properly
- Full stdlib compatibility (`*Error` implements `Unwrap()`)
- Type preservation in error chains
- Uses `errors.Join` internally

See README.md for complete usage documentation and examples.

## Development Commands

**Run all tests:**
```bash
go test -v ./...
```

**Run a specific test:**
```bash
go test -v -run TestName
```

**Run tests with coverage:**
```bash
go test -cover ./...
```

**Build/verify the module:**
```bash
go build
```

**Format code:**
```bash
go fmt ./...
```

## Architecture

### Core Types

**`Error` struct (errors.go:19-23)**: The main error wrapper type that wraps a standard Go error. Contains a single unexported `err error` field for encapsulation. Implements the `error` interface and `Unwrap()` method for compatibility with stdlib error functions.

**`FormattedError` struct (formatted.go:5-9)**: Represents a formatted instance of an Error with template variables applied. Contains a reference to the parent `*Error` and the formatted message string. Implements `Unwrap()` to return the parent.

### Key Design Patterns

1. **Error wrapping**: The library uses `errors.Join` internally in `Wrap()` and `Wrapf()` functions. These functions preserve the full error chain - if you wrap an `*Error`, it stays as `*Error` in the chain (not unwrapped).

2. **Unwrap method**: `*Error` implements `Unwrap() error` which returns the wrapped error. This allows stdlib `errors.Is` and `errors.As` to work correctly with the error chain.

3. **Compatibility functions**: `Is()` and `As()` are convenience wrappers around stdlib functions. `Is()` adds special handling for `FormattedError` types by unwrapping them to their parent first. Both functions work identically to their stdlib counterparts otherwise.

4. **Formatted errors**: Created via `Error.Format(args...)` which applies `fmt.Sprintf` to the error message. The error message MUST be a valid format template. This function will panic if format verbs don't match arguments or if non-template strings contain `%` (e.g., "100% complete"). Use a pointer receiver to preserve the parent reference.

5. **Typed nil detection**: `IsNil()` uses `unsafe.Pointer` to check if an error interface's underlying data pointer is nil, solving the common Go problem where `var e *SomeError = nil; return e` results in a non-nil error interface. This is useful for defensive programming when dealing with external libraries but relies on internal Go implementation details.

### JSON Marshalling

Errors are marshalled as arrays of strings by splitting on newlines. This allows joined/wrapped errors to serialize as multiple messages in order.

### Important Notes

1. **Encapsulation**: The `err` field in `Error` is unexported. Do not access it directly - use `Unwrap()` if you need the wrapped error.

2. **Wrapping nil**: `Wrap(nil, "msg")` returns a non-nil error with just the message. This is intentional to provide context even when the underlying error is nil.

3. **Stdlib compatibility**: Because `*Error` implements `Unwrap()`, you can use stdlib `errors.Is` and `errors.As` directly. The package's `Is()` and `As()` functions are provided for convenience and handle `FormattedError` properly.

4. **Format safety**: Only call `Format()` on errors that were created with format templates. Regular error messages containing `%` will produce error indicators (e.g., "%!(NOVERB)") in the output.

## Testing

Uses `github.com/stretchr/testify` for assertions. Test files follow the pattern `*_test.go` and are located in the same package.

### Test Coverage

**errors_test.go:**
- Basic error creation, wrapping, Is/As functionality
- JSON marshalling (single errors, wrapped errors, natural newlines, empty messages)
- Unwrap method and function behavior
- Stdlib compatibility (errors.Is, errors.As with nested chains)
- Error chain preservation (type information maintained)
- Join function with multiple/nil errors
- Nil wrapping behavior
- Deep error chains (1000+ levels)

**formatted_test.go:**
- Format method with valid templates
- Format behavior with edge cases (%, missing args, extra args)
- FormattedError unwrapping and Is() matching
- Unformatted() helper function

**isnil_test.go:**
- Comprehensive typed nil detection tests
- Various error types (pointer, value, alias)
- Stdlib error types
- Wrapped errors with typed nils