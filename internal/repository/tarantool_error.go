// Errors used internally for distinguishing between faulty cases.

package repository

import "fmt"

type TarantoolError struct {
	Code    int    // http code
	Message string // error message
}

// Satisfy Error interface.
func (err *TarantoolError) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

func NewTarantoolError(code int, message string) *TarantoolError {
	return &TarantoolError{Code: code, Message: message}
}

var (
	ErrNotFound      = NewTarantoolError(404, "key not found")
	ErrAlreadyExists = NewTarantoolError(409, "key already exists")
)
