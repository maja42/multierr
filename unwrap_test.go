package multierr

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_Unwrap(t *testing.T) {
	mErr3 := Append(
		errors.New("c1"),
	)
	mErr2 := Append(
		errors.New("b1"),
		errors.New("b2"),
		mErr3,
	)
	mErr1 := Append(
		errors.New("a1"),
		errors.New("a2"),
		mErr2,
		errors.New("a3"),
	)

	err := errors.Unwrap(mErr1)
	assert.EqualError(t, err, "a1")

	err = errors.Unwrap(err)
	assert.EqualError(t, err, "a2")

	err = errors.Unwrap(err)
	assert.EqualError(t, err, "b1")

	err = errors.Unwrap(err)
	assert.EqualError(t, err, "b2")

	err = errors.Unwrap(err)
	assert.EqualError(t, err, "c1")

	err = errors.Unwrap(err)
	assert.EqualError(t, err, "a3")

	err = errors.Unwrap(err)
	assert.Nil(t, err)
}

func TestError_Unwrap_nothing(t *testing.T) {
	assert.Nil(t, (&Error{}).Unwrap())

	var typedNil *Error
	//goland:noinspection GoNilness
	assert.Nil(t, typedNil.Unwrap())
}

func TestError_Is(t *testing.T) {
	mErr3 := Append(
		errors.New("c1"),
	)
	mErr2 := Append(
		errors.New("b1"),
		errors.New("b2"),
		mErr3,
		io.EOF,
	)
	mErr1 := Append(
		errors.New("a1"),
		errors.New("a2"),
		mErr2,
		errors.New("a3"),
	)

	assert.False(t, errors.Is(mErr3, io.EOF))
	assert.True(t, errors.Is(mErr2, io.EOF))
	assert.True(t, errors.Is(mErr1, io.EOF))
}

type testErr struct{}

func (t *testErr) Error() string {
	return ""
}

func TestError_As(t *testing.T) {
	needle := testErr{}
	mErr3 := Append(
		errors.New("c1"),
	)
	mErr2 := Append(
		errors.New("b1"),
		errors.New("b2"),
		mErr3,
		&needle,
	)
	mErr1 := Append(
		errors.New("a1"),
		errors.New("a2"),
		mErr2,
		errors.New("a3"),
	)

	var res *testErr
	assert.False(t, errors.As(mErr3, &res))
	assert.Nil(t, res)
	assert.True(t, errors.As(mErr2, &res))
	assert.Equal(t, &needle, res)
	assert.True(t, errors.As(mErr1, &res))
	assert.Equal(t, &needle, res)
}
