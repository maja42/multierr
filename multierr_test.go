package multierr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Titled(t *testing.T) {
	t.Run("simple error", func(t *testing.T) {
		orig := errors.New("err")
		err := Titled(orig, "title")
		assert.Equal(t, "title\n  - err", err.Error())
	})

	t.Run("multi error", func(t *testing.T) {
		orig := &Error{
			Errors: []error{errors.New("err")},
		}
		err := Titled(orig, "title")
		assert.Equal(t, "title\n  - err", err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		err := Titled(nil, "title")
		assert.NoError(t, err)
	})

	t.Run("typed nil error", func(t *testing.T) {
		var orig *Error
		err := Titled(orig, "title")
		assert.NotNil(t, err)
		assert.Nil(t, err.(*Error).Errors)

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

func Test_Titledf(t *testing.T) {
	orig := errors.New("err")
	err := Titledf(orig, "formatted %s %d", "title", 42)
	assert.Equal(t, "formatted title 42\n  - err", err.Error())
}

func Test_Append_appendSimpleError(t *testing.T) {
	simpleErr := errors.New("err")
	formatter := func([]error) string {
		return "test formatter"
	}

	t.Run("into existing error", func(t *testing.T) {
		original := &Error{
			Formatter: formatter,
			Errors:    []error{simpleErr},
		}

		result := Append(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 2)
		assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))

		original = &Error{
			Formatter: formatter,
		}
		result = Append(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
		assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))
	})

	t.Run("into typed nil", func(t *testing.T) {
		var original *Error
		result := Append(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
	})

	t.Run("into nil", func(t *testing.T) {
		var original error
		result := Append(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
	})
}

func Test_Append_appendMultiError(t *testing.T) {
	err := errors.New("err")

	multi1 := Append(err, err, err)
	multi2 := Append(err, err, err, err, err)

	result := Append(nil, multi1)
	assert.Len(t, result.(*Error).Errors, 1)

	result = Append(err, multi1)
	assert.Len(t, result.(*Error).Errors, 2)

	formatter := func([]error) string {
		return "test formatter"
	}
	multi1.(*Error).Formatter = formatter
	result = Append(multi1, multi2)
	assert.Len(t, result.(*Error).Errors, 4)
	assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))

	assert.Equal(t, result, multi1)
	assert.Len(t, multi2.(*Error).Errors, 5)
}

func Test_Append_appendMultipleErrors(t *testing.T) {
	err := errors.New("err")

	var original error
	result := Append(original, err, err, err)
	assert.Len(t, result.(*Error).Errors, 3)

	multi := Append(err, err, err)
	result = Append(nil, err, multi, nil, err, multi)
	assert.Len(t, result.(*Error).Errors, 4)
}

func Test_Merge_mergeSimpleError(t *testing.T) {
	simpleErr := errors.New("err")
	formatter := func([]error) string {
		return "test formatter"
	}

	t.Run("into existing error", func(t *testing.T) {
		original := &Error{
			Formatter: formatter,
			Errors:    []error{simpleErr},
		}

		result := Merge(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 2)
		assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))

		original = &Error{
			Formatter: formatter,
		}
		result = Merge(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
		assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))
	})

	t.Run("into typed nil", func(t *testing.T) {
		var original *Error
		result := Merge(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
	})

	t.Run("into nil", func(t *testing.T) {
		var original error
		result := Merge(original, simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
	})
}

func Test_Merge_mergeMultiError(t *testing.T) {
	err := errors.New("err")

	multi1 := Merge(err, err, err)
	multi2 := Merge(err, err, err, err, err)

	result := Merge(nil, multi1)
	assert.Len(t, result.(*Error).Errors, 3)

	result = Merge(err, multi1)
	assert.Len(t, result.(*Error).Errors, 4)

	formatter := func([]error) string {
		return "test formatter"
	}
	multi1.(*Error).Formatter = formatter
	result = Merge(multi1, multi2)
	assert.Len(t, result.(*Error).Errors, 8)
	assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))

	assert.Equal(t, result, multi1)
	assert.Len(t, multi2.(*Error).Errors, 5)
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

	assert.NoError(t, result)
}

func Test_Merge_mergeMultipleErrors(t *testing.T) {
	err := errors.New("err")

	var original error
	result := Merge(original, err, err, err)
	assert.Len(t, result.(*Error).Errors, 3)

	multi := Merge(err, err, err)
	result = Merge(nil, err, multi, nil, err, multi)
	assert.Len(t, result.(*Error).Errors, 8)
}

func TestInspect(t *testing.T) {

	t.Run("simple error", func(t *testing.T) {
		err := errors.New("err")
		errs := Inspect(err)
		assert.Equal(t, []error{err}, errs)
	})

	t.Run("multi error", func(t *testing.T) {
		err := errors.New("err")
		multiErr := &Error{
			Errors: []error{err, err},
		}
		errs := Inspect(multiErr)
		assert.Equal(t, []error{err, err}, errs)
	})

	t.Run("nil error", func(t *testing.T) {
		errs := Inspect(nil)
		assert.Nil(t, errs)
	})

	t.Run("typed nil error", func(t *testing.T) {
		var err *Error
		errs := Inspect(err)
		assert.Nil(t, errs)
	})

	t.Run("empty error", func(t *testing.T) {
		errs := Inspect(&Error{})
		assert.Nil(t, errs)
	})

}
