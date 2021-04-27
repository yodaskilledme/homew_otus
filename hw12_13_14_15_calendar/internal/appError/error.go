package appError

import (
	"strings"
)

type AppError struct {
	// logical operation
	Op string
	// wrapped appError
	Err error
	// human readable message
	Message string
	// machine readable appError code
	Code string
}

// Error returns the string representation of the appError message.
func (e AppError) Error() string {
	var buf strings.Builder

	// Print the current operation in our stack, if any.
	if e.Op != "" {
		buf.WriteString(e.Op)
		buf.WriteString(": ")
	}

	// If wrapping an appError, print its Error() message.
	// Otherwise print the appError code & message.
	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != "" {
			buf.WriteRune('<')
			buf.WriteString(e.Code)
			buf.WriteRune('>')
		}
		if e.Code != "" && e.Message != "" {
			// add a space
			buf.WriteRune(' ')
		}
		buf.WriteString(e.Message)
	}

	return buf.String()
}

func OpError(op string, err error) *AppError {
	return &AppError{Op: op, Err: err}
}

func OpErrorOrNil(op string, err error) error {
	if err == nil {
		return nil
	}

	return &AppError{Op: op, Err: err}
}
