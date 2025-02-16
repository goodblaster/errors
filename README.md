# go-errors

Simple wrapper around Golang's new error capabilities.

Uses Wrap instead of Join. Personal preference.

Allows JSON marshalling of errors. Errors are marshalled as a slice of strings, with the highest level error first.