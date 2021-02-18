package multierr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_ErrorOrNil(t *testing.T) {
	err := &Error{
		Errors: []error{errors.New("err")},
	}
	assert.Equal(t, err, err.ErrorOrNil())

	err = &Error{
		Errors: []error{},
	}
	assert.Nil(t, err.ErrorOrNil())

	err = &Error{}
	assert.Nil(t, err.ErrorOrNil())

	var typedErr *Error
	//goland:noinspection GoNilness
	assert.Nil(t, typedErr.ErrorOrNil())
}

func Test_Titled(t *testing.T) {
	t.Run("simple error", func(t *testing.T) {
		err := errors.New("err")
		err = Titled(err, "title")
		assert.Equal(t, "title\n  - err", err.Error())
	})

	t.Run("multi error", func(t *testing.T) {
		err := &Error{
			Errors: []error{errors.New("err")},
		}
		err = Titled(err, "title")
		assert.Equal(t, "title\n  - err", err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		err := Titled(nil, "title")
		assert.Nil(t, err)
	})

	t.Run("typed nil error", func(t *testing.T) {
		var err *Error
		err = Titled(err, "title")
		assert.NotNil(t, err)
		assert.Nil(t, err.Errors)

		assert.Equal(t, "no errors occurred", err.Error())

		err = Append(err, errors.New("err"))
		assert.Equal(t, "title\n  - err", err.Error())
	})

	t.Run("empty error", func(t *testing.T) {
		err := Titled(&Error{}, "title")
		assert.Equal(t, "no errors occurred", err.Error())

		err = Append(err, errors.New("err"))
		assert.Equal(t, "title\n  - err", err.Error())
	})
}

func Test_Append_appendSimpleError(t *testing.T) {
	err := errors.New("err")
	formatter := func([]error) string {
		return "test formatter"
	}

	t.Run("into existing error", func(t *testing.T) {
		original := &Error{
			Formatter: formatter,
			Errors:    []error{err},
		}

		result := Append(original, err)
		assert.Len(t, result.Errors, 2)
		assert.Equal(t, "test formatter", result.Formatter(nil))

		original = &Error{
			Formatter: formatter,
		}
		result = Append(original, err)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, "test formatter", result.Formatter(nil))
	})

	t.Run("into typed nil", func(t *testing.T) {
		var original *Error
		result := Append(original, err)
		assert.Len(t, result.Errors, 1)
	})

	t.Run("into nil", func(t *testing.T) {
		var original error
		result := Append(original, err)
		assert.Len(t, result.Errors, 1)
	})
}

func Test_Append_appendMultiError(t *testing.T) {
	err := errors.New("err")

	multi1 := Append(err, err, err)
	multi2 := Append(err, err, err, err, err)

	result := Append(nil, multi1)
	assert.Len(t, result.Errors, 1)

	result = Append(err, multi1)
	assert.Len(t, result.Errors, 2)

	formatter := func([]error) string {
		return "test formatter"
	}
	multi1.Formatter = formatter
	result = Append(multi1, multi2)
	assert.Len(t, result.Errors, 4)
	assert.Equal(t, "test formatter", result.Formatter(nil))

	assert.Equal(t, result, multi1)
	assert.Len(t, multi2.Errors, 5)
}

func Test_Append_appendMultipleErrors(t *testing.T) {
	err := errors.New("err")

	var original error
	result := Append(original, err, err, err)
	assert.Len(t, result.Errors, 3)

	multi := Append(err, err, err)
	result = Append(nil, err, multi, nil, err, multi)
	assert.Len(t, result.Errors, 4)
}

func Test_Merge_mergeSimpleError(t *testing.T) {
	err := errors.New("err")
	formatter := func([]error) string {
		return "test formatter"
	}

	t.Run("into existing error", func(t *testing.T) {
		original := &Error{
			Formatter: formatter,
			Errors:    []error{err},
		}

		result := Merge(original, err)
		assert.Len(t, result.Errors, 2)
		assert.Equal(t, "test formatter", result.Formatter(nil))

		original = &Error{
			Formatter: formatter,
		}
		result = Merge(original, err)
		assert.Len(t, result.Errors, 1)
		assert.Equal(t, "test formatter", result.Formatter(nil))
	})

	t.Run("into typed nil", func(t *testing.T) {
		var original *Error
		result := Merge(original, err)
		assert.Len(t, result.Errors, 1)
	})

	t.Run("into nil", func(t *testing.T) {
		var original error
		result := Merge(original, err)
		assert.Len(t, result.Errors, 1)
	})
}

func Test_Merge_mergeMultiError(t *testing.T) {
	err := errors.New("err")

	multi1 := Merge(err, err, err)
	multi2 := Merge(err, err, err, err, err)

	result := Merge(nil, multi1)
	assert.Len(t, result.Errors, 3)

	result = Merge(err, multi1)
	assert.Len(t, result.Errors, 4)

	formatter := func([]error) string {
		return "test formatter"
	}
	multi1.Formatter = formatter
	result = Merge(multi1, multi2)
	assert.Len(t, result.Errors, 8)
	assert.Equal(t, "test formatter", result.Formatter(nil))

	assert.Equal(t, result, multi1)
	assert.Len(t, multi2.Errors, 5)
}

func Test_Merge_Nothing(t *testing.T) {
	var original error
	var typedNil1 error
	var typedNil2 *Error

	result := Merge(original, nil)
	result = Merge(result, &Error{})
	result = Merge(result, typedNil1)
	result = Merge(result, typedNil2)
	result = Merge(result, nil)
	result = Merge(result, nil, &Error{}, typedNil1, typedNil2)

	assert.Nil(t, result)
}

func Test_Merge_mergeMultipleErrors(t *testing.T) {
	err := errors.New("err")

	var original error
	result := Merge(original, err, err, err)
	assert.Len(t, result.Errors, 3)

	multi := Merge(err, err, err)
	result = Merge(nil, err, multi, nil, err, multi)
	assert.Len(t, result.Errors, 8)
}
