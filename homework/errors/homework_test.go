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
