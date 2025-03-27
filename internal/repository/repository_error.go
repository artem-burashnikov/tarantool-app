// Errors used internally by repository for distinguishing between faulty cases.

package repository

import "fmt"

type RepositoryError struct {
	Message string // error message
}

// Satisfy Error interface.
func (err *RepositoryError) Error() string {
	return fmt.Sprintf("%s", err.Message)
}

func NewRepositoryError(message string) *RepositoryError {
	return &RepositoryError{Message: message}
}

var (
	ErrNotFound            = NewRepositoryError("key not found")
	ErrAlreadyExists       = NewRepositoryError("key already exists")
	ErrInsertOperationFail = NewRepositoryError("insert operation failed")
	ErrSelectOperationFail = NewRepositoryError("select operation failed")
	ErrUpdateOperationFail = NewRepositoryError("update operation failed")
	ErrDeleteOperationFail = NewRepositoryError("delete operation failed")
)
