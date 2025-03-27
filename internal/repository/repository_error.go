// Errors used internally by repository for distinguishing between faulty cases.

package repository

import "fmt"

type RepositoryError struct {
	Code    int    // http code
	Message string // error message
}

// Satisfy Error interface.
func (err *RepositoryError) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

func NewRepositoryError(code int, message string) *RepositoryError {
	return &RepositoryError{Code: code, Message: message}
}

var (
	ErrNotFound      = NewRepositoryError(404, "key not found")
	ErrAlreadyExists = NewRepositoryError(409, "key already exists")
)
