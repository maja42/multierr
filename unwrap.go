package multierr

import "errors"

// Unwrap returns the first error in Error or nil if there are no errors.
// By repeatedly calling Unwrap on the return value (until nil is returned),
// all (recursively) contained sub-errors can be obtained.
// Errors are unwrapped depth-first.
// This implements errors.Is/errors.As/errors.Unwrap methods from the standard library.
// Appending new errors while unwrapping has no effect (shallow copy).
func (e *Error) Unwrap() error {
	if e == nil || len(e.Errors) == 0 {
		return nil
	}
	if len(e.Errors) == 1 {
		return e.Errors[0]
	}

	// copy, to be independent if new errors are appended/merged while unwrapping
	errs := make([]error, len(e.Errors))
	copy(errs, e.Errors)
	return chain(errs)
}

type chain []error

// Unwrap returns the next error or nil if there are no errors left.
func (e chain) Unwrap() error {
	if len(e) == 1 {
		// current element is the last one, nothing to unwrap
		return nil
	}

	if err, ok := e[1].(*Error); ok {
		// multi-error -> depth-first search
		return chain(append(err.Errors, e[2:]...))
	}

	return e[1:] // remove first (=current) element
}

// Error implements the error interface
func (e chain) Error() string {
	return e[0].Error()
}

// Is implements errors.Is.
func (e chain) Is(target error) bool {
	return errors.Is(e[0], target)
}

// As implements errors.As.
func (e chain) As(target interface{}) bool {
	return errors.As(e[0], target)
}
