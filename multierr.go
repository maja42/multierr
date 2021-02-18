package multierr

import (
	"fmt"
)

// DefaultFormatter specifies the error formatter that is used for errors that
// don't have a dedicated formatter function specified.
var DefaultFormatter = ListFormatterFunc

// Error is an error type to track multiple errors. This is used to
// accumulate errors in cases and return them as a single "error".
type Error struct {
	Formatter FormatterFunc
	Errors    []error
}

// Error converts the error into a human readable string.
// Uses the error-specific formatter or, if none is specified, the DefaultFormatter.
func (e *Error) Error() string {
	formatter := e.Formatter
	if formatter == nil {
		formatter = DefaultFormatter
	}
	return formatter(e.Errors)
}

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (e *Error) ErrorOrNil() error {
	if e == nil || len(e.Errors) == 0 {
		return nil
	}
	return e
}

// Titled sets the error formatter to a TitledListFormatter.
// The given title is used when calling Error.Error().
//
// If the error is not a multierr.Error, it will be converted.
// Returns nil if if the error is nil.
// This is equivalent of setting Error.Formatter directly.
func Titled(err error, title string) *Error {
	if err == nil {
		return nil
	}
	formatter := TitledListFormatter(title)

	mErr, ok := err.(*Error)
	if mErr == nil {
		mErr = &Error{}
	}
	if !ok {
		mErr.Errors = []error{err}
	}
	mErr.Formatter = formatter

	return mErr
}

// Titledf sets the error formatter to a TitledListFormatter.
// See Titled for more information.
func Titledf(err error, format string, args ...interface{}) *Error {
	title := fmt.Sprintf(format, args...)
	return Titled(err, title)
}

// Append combines all errors into a single multi-error.
// Any nil-error will be ignored. Returns nil if there are no errors.
//
// If err is a multierr.Error, it will be reused (the title and error-slice are kept).
// Otherwise a new multierr.Error is created.
func Append(err error, errs ...error) *Error {
	return combine(false, err, errs...)
}

// Merge combines all errors into a single multi-error.
// Any nil-error will be ignored. Returns nil if there are no errors.
//
// If any errs is a multierr.Error, it will be flattened.
//
// If err is a multierr.Error, it will be reused (the title and error-slice are kept).
// Otherwise a new multierr.Error is created.
func Merge(err error, errs ...error) *Error {
	return combine(true, err, errs...)
}

func combine(flatten bool, err error, errs ...error) *Error {
	result, ok := err.(*Error)
	if result == nil {
		result = &Error{
			Errors: make([]error, 0, len(errs)+1),
		}
	}
	if !ok && err != nil { // err was not a multi error
		result.Errors = append(result.Errors, err)
	}

	for _, e := range errs {
		if e == nil {
			continue
		}
		multiErr, ok := e.(*Error)
		if ok && (multiErr == nil || len(multiErr.Errors) == 0) {
			continue
		}

		if ok && flatten {
			result.Errors = append(result.Errors, multiErr.Errors...)
		} else {
			result.Errors = append(result.Errors, e)
		}
	}
	if len(result.Errors) == 0 {
		return nil
	}
	return result
}
