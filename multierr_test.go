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
		assert.Nil(t, err)
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

func Test_Prefixed(t *testing.T) {
	t.Run("simple error", func(t *testing.T) {
		orig := errors.New("err")
		err := Prefixed(orig, "prefix: ")
		assert.Equal(t, "prefix: err", err.Error())
	})

	t.Run("multi error", func(t *testing.T) {
		orig := &Error{
			Errors: []error{errors.New("err")},
		}
		err := Prefixed(orig, "prefix: ")
		assert.Equal(t, "prefix: err", err.Error())
	})

	t.Run("nil error", func(t *testing.T) {
		err := Prefixed(nil, "prefix: ")
		assert.Nil(t, err)
	})

	t.Run("typed nil error", func(t *testing.T) {
		var orig *Error
		err := Prefixed(orig, "prefix: ")
		assert.NotNil(t, err)
		assert.Nil(t, err.(*Error).Errors)

		assert.Equal(t, "no errors occurred", err.Error())

		err = Append(err, errors.New("err"))
		assert.Equal(t, "prefix: err", err.Error())
	})

	t.Run("empty error", func(t *testing.T) {
		err := Prefixed(&Error{}, "prefix: ")
		assert.Equal(t, "no errors occurred", err.Error())

		err = Append(err, errors.New("err"))
		assert.Equal(t, "prefix: err", err.Error())
	})
}

func Test_Prefixedf(t *testing.T) {
	orig := errors.New("err")
	err := Prefixedf(orig, "formatted %s %d: ", "prefix", 42)
	assert.Equal(t, "formatted prefix 42: err", err.Error())
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

func Test_MergePrefixed_mergeSimpleError(t *testing.T) {
	simpleErr := errors.New("err")
	formatter := func([]error) string {
		return "test formatter"
	}

	t.Run("into existing error", func(t *testing.T) {
		original := &Error{
			//Formatter: formatter,
			Errors: []error{simpleErr},
		}

		result := MergePrefixed(original, "prefix: ", simpleErr)
		assert.Len(t, result.(*Error).Errors, 2)
		assert.Error(t, result, "2 errors occurred:\n  - err\n  - prefix: err")

		original = &Error{
			Formatter: formatter,
		}
		result = MergePrefixed(original, "prefix: ", simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
		assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))
	})

	t.Run("into typed nil", func(t *testing.T) {
		var original *Error
		result := MergePrefixed(original, "prefix: ", simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
		assert.Error(t, result, "1 error occurred:\n  - prefix: err")
	})

	t.Run("into nil", func(t *testing.T) {
		var original error
		result := MergePrefixed(original, "prefix: ", simpleErr)
		assert.Len(t, result.(*Error).Errors, 1)
		assert.Error(t, result, "1 error occurred:\n  - prefix: err")
	})
}

func Test_MergePrefixed_mergeMultiError(t *testing.T) {
	err := errors.New("err")

	multi1 := MergePrefixed(err, "prefix 1: ", err, err)
	multi2 := MergePrefixed(err, "prefix 2: ", err, err, err, err)

	result := MergePrefixed(nil, "prefix 3: ", multi1)
	assert.Len(t, result.(*Error).Errors, 3)
	assert.Error(t, result, "3 errors occurred:\n"+
		"  - prefix 3: err\n"+
		"  - prefix 3: prefix 1: err\n"+
		"  - prefix 3: prefix 1: err")

	result = MergePrefixed(err, "prefix 4: ", multi1)
	assert.Len(t, result.(*Error).Errors, 4)
	assert.Error(t, result, "4 errors occurred:\n"+
		"  - err\n"+
		"  - prefix 4: err\n"+
		"  - prefix 4: prefix 1: err\n"+
		"  - prefix 4: prefix 1: err")

	formatter := func([]error) string {
		return "test formatter"
	}
	multi1.(*Error).Formatter = formatter
	result = MergePrefixed(multi1, "prefix 5: ", multi2)
	assert.Len(t, result.(*Error).Errors, 8)
	assert.Equal(t, "test formatter", result.(*Error).Formatter(nil))

	assert.Equal(t, result, multi1)
	assert.Len(t, multi2.(*Error).Errors, 5)
}

func Test_MergePrefixed_Nothing(t *testing.T) {
	var original error
	var typedNil1 error
	var typedNil2 *Error

	result := MergePrefixed(original, "prefix: ", nil)
	result = MergePrefixed(result, "prefix: ", &Error{})
	result = MergePrefixed(result, "prefix: ", typedNil1)
	result = MergePrefixed(result, "prefix: ", typedNil2)
	result = MergePrefixed(result, "prefix: ", nil)
	result = MergePrefixed(result, "prefix: ", nil, &Error{}, typedNil1, typedNil2)

	assert.NoError(t, result)
}

func Test_MergePrefixed_mergeMultipleErrors(t *testing.T) {
	err := errors.New("err")

	var original error
	result := MergePrefixed(original, "prefix: ", err, err, err)
	assert.Len(t, result.(*Error).Errors, 3)
	assert.Error(t, result, "3 errors occurred:\n"+
		"  - prefix: err\n"+
		"  - prefix: err\n"+
		"  - prefix: err")

	multi := MergePrefixed(err, "prefix: ", err, err)
	result = MergePrefixed(nil, "prefix 2: ", err, multi, nil, err, multi)
	assert.Len(t, result.(*Error).Errors, 8)

	assert.Error(t, result, "8 errors occurred:\n"+
		"  - prefix2: err\n"+ // direct
		"  - prefix2: err\n"+ // part of multi
		"  - prefix2: prefix: err\n"+
		"  - prefix2: prefix: err\n"+
		"  - prefix2: err\n"+ // direct
		"  - prefix2: err\n"+ // part of multi
		"  - prefix2: prefix: err\n"+
		"  - prefix2: prefix: err")
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
