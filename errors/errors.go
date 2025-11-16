// Package errors extends the standard library errors package with some
// additional convenience functions.
package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Relay to the methods in the standard library package, so callers don't have
// to import both.

func As(err error, target any) bool { return errors.As(err, target) }
func Is(err, target error) bool     { return errors.Is(err, target) }
func New(text string) error         { return errors.New(text) }
func Unwrap(err error) error        { return errors.Unwrap(err) }

func NewF(f string, a ...any) error { return fmt.Errorf(f, a...) }

type joinError struct {
	errs []error
}

func (je joinError) Error() string {
	var s []string
	for _, e := range je.errs {
		s = append(s, e.Error())
	}
	return strings.Join(s, "\n")
}

func (je joinError) Unwrap() []error {
	return je.errs
}

// Join is like errors.Join, but each successive call to it adds to a flat
// slice of errors rather than making nested slices.
func Join(errs ...error) error {
	var joined []error

	for _, err := range errs {
		switch err := err.(type) {
		case nil:
			// nothing
		case joinError:
			joined = append(joined, err.errs...)
		default:
			joined = append(joined, err)
		}
	}
	switch len(joined) {
	case 0:
		return nil
	case 1:
		return joined[0]
	default:
		return joinError{joined}
	}
}

// JoinString defines a new error (as in errors.New) and joins it to the
// provided error (as in errors.Join).
func JoinString(err error, str string) error {
	return Join(err, New(str))
}

// JoinF defines a new formatted error (as in fmt.Errorf) and joins it to the
// provide error (as in errors.Join).
func JoinF(err error, f string, a ...any) error {
	return Join(err, fmt.Errorf(f, a...))
}

// JoinCall joins the error return from a value-and-error function to the
// supplied err.  If f is a function that takes any arguments and returns a
// value of type Tand an error, you can call this like:
//
//	val, err = errors.JoinCall[T](err)(f(...))
//
// Unfortunately the type parameter T cannot be inferred and is required.
func JoinCall[T any](err error) func(T, error) (T, error) {
	return func(val T, fnerr error) (T, error) {
		return val, Join(err, fnerr)
	}
}

// UnwrapJoined returns the list of errors that were joined with Join.
func UnwrapJoined(err error) []error {
	if err, ok := err.(interface{ Unwrap() []error }); ok {
		return err.Unwrap()
	}
	return nil
}

// AddPrefix adds the specified prefix to the error.  If the error is a join of
// multiple errors, it adds the prefix to all of them.
func AddPrefix(err error, prefix string) error {
	switch err := err.(type) {
	case nil:
		return nil
	case joinError:
		for i := range err.errs {
			err.errs[i] = AddPrefix(err.errs[i], prefix)
		}
		return err
	default:
		return fmt.Errorf("%s: %w", prefix, err)
	}
}
