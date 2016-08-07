package errors

import (
	goerrors "errors"
	"fmt"
)

var _ error = nestedError{}

type nestedError struct {
	message string
	cause   error
}

func (this nestedError) Error() string {
	if this.cause == nil {
		return fmt.Sprintln(this.message)
	} else {
		return fmt.Sprintln(this.message, "\n caused by: ", this.cause.Error())
	}
}

func New(message string) error {
	return goerrors.New(message)
}

func NewNested(message string, cause error) error {
	n := nestedError{
		message: message,
		cause:   cause,
	}
	return n
}
