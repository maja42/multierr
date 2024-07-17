package multierr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ListFormatterFunc_noError(t *testing.T) {
	DefaultFormatter = ListFormatterFunc
	assert.Equal(t, "no errors occurred", (&Error{}).Error())

	// Calling on a nil-error panics
	var err *Error
	assert.Panics(t, func() {
		_ = err.Error()
	})
}

func Test_ListFormatterFunc_singleError(t *testing.T) {
	DefaultFormatter = ListFormatterFunc
	err := Append(nil, errors.New("some error"))
	expected := "1 error occurred:\n  - some error"

	assert.Equal(t, expected, err.Error())

	// multiline-error:
	err = Append(nil, errors.New("some error\nsecond line"))
	expected = "1 error occurred:\n" +
		"  - some error\n" +
		"    second line"

	assert.Equal(t, expected, err.Error())
}

func Test_ListFormatterFunc_multipleErrors(t *testing.T) {
	DefaultFormatter = ListFormatterFunc
	err := Append(nil,
		errors.New("error 1"),
		errors.New("error 2"))

	expected := "2 errors occurred:\n" +
		"  - error 1\n" +
		"  - error 2"

	assert.Equal(t, expected, err.Error())

	// multiline-errors:
	err = Append(nil,
		errors.New("error 1\nsecond line 1"),
		errors.New("error 2\nsecond line 2"))
	expected = "2 errors occurred:\n" +
		"  - error 1\n" +
		"    second line 1\n" +
		"  - error 2\n" +
		"    second line 2"

	assert.Equal(t, expected, err.Error())
}

func Test_TitledListFormatterFunc_noError(t *testing.T) {
	DefaultFormatter = TitledListFormatter("the title")

	assert.Equal(t, "no errors occurred", (&Error{}).Error())

	// Calling on a nil-error panics
	var err *Error
	assert.Panics(t, func() {
		_ = err.Error()
	})
}

func Test_TitledListFormatterFunc_singleError(t *testing.T) {
	DefaultFormatter = TitledListFormatter("the title")

	err := Append(nil, errors.New("some error"))
	expected := "the title\n  - some error"

	assert.Equal(t, expected, err.Error())

	// multiline-error:
	err = Append(nil, errors.New("some error\nsecond line"))
	expected = "the title\n" +
		"  - some error\n" +
		"    second line"

	assert.Equal(t, expected, err.Error())
}

func Test_TitledListFormatterFunc_multipleErrors(t *testing.T) {
	DefaultFormatter = TitledListFormatter("the title")

	err := Append(nil,
		errors.New("error 1"),
		errors.New("error 2"))

	expected := "the title\n" +
		"  - error 1\n" +
		"  - error 2"

	assert.Equal(t, expected, err.Error())

	// multiline-errors:
	err = Append(nil,
		errors.New("error 1\nsecond line 1"),
		errors.New("error 2\nsecond line 2"))
	expected = "the title\n" +
		"  - error 1\n" +
		"    second line 1\n" +
		"  - error 2\n" +
		"    second line 2"

	assert.Equal(t, expected, err.Error())
}

func Test_PrefixedListFormatterFunc_noError(t *testing.T) {
	DefaultFormatter = PrefixedListFormatter("prefix: ")

	assert.Equal(t, "no errors occurred", (&Error{}).Error())

	// Calling on a nil-error panics
	var err *Error
	assert.Panics(t, func() {
		_ = err.Error()
	})
}

func Test_PrefixedListFormatterFunc_singleError(t *testing.T) {
	DefaultFormatter = PrefixedListFormatter("prefix: ")

	err := Append(nil, errors.New("some error"))
	expected := "prefix: some error"

	assert.Equal(t, expected, err.Error())

	// multiline-error:
	err = Append(nil, errors.New("some error\nsecond line"))
	expected = "" +
		"prefix: some error\n" +
		"        second line"

	assert.Equal(t, expected, err.Error())
}

func Test_PrefixedListFormatterFunc_multipleErrors(t *testing.T) {
	DefaultFormatter = PrefixedListFormatter("prefix: ")

	err := Append(nil,
		errors.New("error 1"),
		errors.New("error 2"))

	expected := "" +
		"prefix: error 1\n" +
		"prefix: error 2"

	assert.Equal(t, expected, err.Error())

	// multiline-errors:
	err = Append(nil,
		errors.New("error 1\nsecond line 1"),
		errors.New("error 2\nsecond line 2"))
	expected = "" +
		"prefix: error 1\n" +
		"        second line 1\n" +
		"prefix: error 2\n" +
		"        second line 2"

	assert.Equal(t, expected, err.Error())
}

func Test_CustomFormatterFunc(t *testing.T) {
	err := Append(nil,
		errors.New("error 1"),
		errors.New("error 2")).(*Error)

	err.Formatter = func(errs []error) string {
		assert.Len(t, errs, 2)
		assert.EqualError(t, errs[0], "error 1")
		assert.EqualError(t, errs[1], "error 2")
		return "custom format"
	}

	assert.Equal(t, "custom format", err.Error())
}
