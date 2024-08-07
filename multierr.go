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

// Titled sets the error formatter to a TitledListFormatter.
// The given title is used when calling Error.Error().
//
// If the error is not a multierr.Error, it will be converted.
// Returns nil if the error is nil. Otherwise, the result is always an *Error.
// This is equivalent of setting Error.Formatter directly.
func Titled(err error, title string) error {
	return withCustomFormatter(err, TitledListFormatter(title))
}

// Titledf sets the error formatter to a TitledListFormatter.
// See Titled for more information.
func Titledf(err error, format string, args ...interface{}) error {
	title := fmt.Sprintf(format, args...)
	return Titled(err, title)
}

// Prefixed sets the error formatter to a PrefixedListFormatter.
// The given prefix is used when calling Error.Error().
//
// If the error is not a multierr.Error, it will be converted.
// Returns nil if the error is nil. Otherwise, the result is always an *Error.
// This is equivalent of setting Error.Formatter directly.
func Prefixed(err error, prefix string) error {
	return withCustomFormatter(err, PrefixedListFormatter(prefix))
}

// Prefixedf sets the error formatter to a PrefixedListFormatter.
// See Prefixed for more information.
func Prefixedf(err error, format string, args ...interface{}) error {
	title := fmt.Sprintf(format, args...)
	return Prefixed(err, title)
}

func withCustomFormatter(err error, formatter FormatterFunc) error {
	if err == nil {
		return nil
	}

	mErr, ok := err.(*Error)
	if mErr == nil && ok {
		// typed nil-error
		// special case: do not return nil. Return empty *Error that contains the custom formatter.
	}
	if mErr == nil {
		mErr = &Error{}
	}
	if !ok { // convert error type
		mErr.Errors = []error{err}
	}
	mErr.Formatter = formatter
	return mErr
}

// Append combines all errors into a single multi-error.
// Any nil-error will be ignored. Returns nil if there are no errors.
// A returned error will always be of type *Error.
//
// If err is a multierr.Error, it will be reused (the title and error-slice are kept).
// Otherwise a new multierr.Error is created.
func Append(err error, errs ...error) error {
	return combine(false, err, "", errs...)
}

// Merge combines all errors into a single multi-error.
// Any nil-error will be ignored. Returns nil if there are no errors.
// A returned error will always be of type *Error.
//
// If any errs is a multierr.Error, it will be flattened.
//
// If err is a multierr.Error, it will be reused (the title and error-slice are kept).
// Otherwise, a new multierr.Error is created.
func Merge(err error, errs ...error) error {
	return combine(true, err, "", errs...)
}

// MergePrefixed combines all errors into a single multi-error.
// Any nil-error will be ignored. Returns nil if there are no errors.
// A returned error will always be of type *Error.
//
// If any errs is a multierr.Error, it will be flattened. Custom formatters are ignored.
// Every merged error in errs will be wrapped using the provided prefix.
//
// If err is a multierr.Error, it will be reused (the formatter and error-slice are kept).
// Otherwise, a new multierr.Error is created.
func MergePrefixed(err error, prefix string, errs ...error) error {
	return combine(true, err, prefix, errs...)
}

func combine(flatten bool, err error, errsPrefix string, errs ...error) error {
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
			subErrors := multiErr.Errors
			if errsPrefix != "" {
				subErrors = make([]error, len(multiErr.Errors))
				for i, err := range multiErr.Errors {
					subErrors[i] = fmt.Errorf("%s%w", errsPrefix, err)
				}
			}
			result.Errors = append(result.Errors, subErrors...)
		} else {
			if errsPrefix != "" {
				result.Errors = append(result.Errors, fmt.Errorf("%s%w", errsPrefix, e))
			} else {
				result.Errors = append(result.Errors, e)
			}
		}
	}
	if len(result.Errors) == 0 {
		return nil
	}
	return result
}

// Inspect returns all embedded sub-errors or nil if there are no errors.
// If err is not a multi-error, an error-slice with one element is returned.
func Inspect(err error) []error {
	if err == nil {
		return nil
	}
	result, ok := err.(*Error)
	if !ok {
		return []error{err}
	}
	if result == nil {
		return nil
	}
	return result.Errors
}
