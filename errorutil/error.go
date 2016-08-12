package errorutil

import (
	"errors"
	"fmt"
)

var _ error = &nestedError{}

type nestedError struct {
	message   string
	errorCode string
	cause     error
}

func (this *nestedError) ErrorCode() string {
	return this.errorCode
}

func (this *nestedError) Error() string {
	if this.cause == nil {
		return fmt.Sprintln(this.message)
	} else {
		return fmt.Sprintln(this.message, "\n caused by: ", this.cause.Error())
	}
}

func New(message string) error {
	return errors.New(message)
}

func NewWithErrorCode(errorCode, message string) error {
	return &nestedError{
		message:   message,
		errorCode: errorCode,
		cause:     nil,
	}
}

func NewNested(message string, cause error) error {
	n := &nestedError{
		message: message,
		cause:   cause,
	}
	return n
}
