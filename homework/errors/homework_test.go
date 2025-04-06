package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errors []error
}

func (e *MultiError) Error() string {
	n := len(e.errors)

	if n == 0 {
		return ""
	}

	if n == 1 {
		return e.errors[0].Error()
	}

	sb := strings.Builder{}
	sb.WriteString(strconv.FormatInt(int64(n), 10))
	sb.WriteString(" errors occured:\n")
	for _, err := range e.errors {
		sb.WriteString("\t* ")
		sb.WriteString(err.Error())
	}

	sb.WriteByte('\n')

	return sb.String()
}

func (e *MultiError) Is(target error) bool {
	if e == nil {
		return e == target
	}

	for _, err := range e.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (e *MultiError) As(target any) bool {
	if e == nil {
		return e == target
	}

	for _, err := range e.errors {
		if errors.As(err, target) {
			return true
		}
	}

	return false
}

func (e *MultiError) Unwrap() error {
	if e == nil || len(e.errors) == 0 {
		return nil
	}

	if len(e.errors) == 1 {
		return e.errors[0]
	}

	return &MultiError{
		errors: e.errors[1:],
	}
}

func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	if err == nil {
		return &MultiError{
			errors: errs,
		}
	}

	mErr, ok := err.(*MultiError)
	if !ok {
		mErr = &MultiError{
			errors: make([]error, 0, len(errs)+1),
		}
		mErr.errors = append(mErr.errors, err)
	}

	mErr.errors = append(mErr.errors, errs...)

	return mErr
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorUnwrap(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")
	err := Append(nil, err1, err2, err3)

	unwrappedOnce := errors.Unwrap(err)
	assert.False(t, errors.Is(unwrappedOnce, err1))
	assert.True(t, errors.Is(unwrappedOnce, err2))
	assert.True(t, errors.Is(unwrappedOnce, err3))

	unwrappedTwice := errors.Unwrap(unwrappedOnce)
	assert.False(t, errors.Is(unwrappedTwice, err1))
	assert.False(t, errors.Is(unwrappedTwice, err2))
	assert.True(t, errors.Is(unwrappedTwice, err3))

	unwrappedThrice := errors.Unwrap(unwrappedTwice)
	assert.False(t, errors.Is(unwrappedThrice, err1))
	assert.False(t, errors.Is(unwrappedThrice, err2))
	assert.True(t, errors.Is(unwrappedThrice, err3))

	assert.Nil(t, errors.Unwrap(unwrappedThrice))
}

func TestMultiErrorIs(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")
	err3 := errors.New("err3")
	err := Append(nil, err1, err2)

	assert.True(t, errors.Is(err, err1))
	assert.False(t, errors.Is(err, err3))
}

// Ниже типы ошибок исключительно для юнит теста поддержки errors.As

type error1 string

func (e error1) Error() string {
	return string(e)
}

type error2 string

func (e error2) Error() string {
	return string(e)
}

type error3 string

func (e error3) Error() string {
	return string(e)
}

// ptr - синтаксический сахар, поинтер на значение какого-то типа Т
func ptr[T any](t T) *T {
	return &t
}

func TestMultiErrorAs(t *testing.T) {
	err1 := error1("err1")
	err2 := error2("err2")
	err := Append(nil, err1, err2)

	assert.True(t, errors.As(err, ptr(error1(""))))
	assert.False(t, errors.As(err, ptr(error3(""))))
}
