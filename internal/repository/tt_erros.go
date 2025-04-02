// Errors used internally by repository for distinguishing between faulty cases.

package repository

import "fmt"

type RepositoryError struct {
	message string
}

var _ error = RepositoryError{} // RepositoryError must satisfy error

func (err RepositoryError) Error() string {
	return fmt.Sprintf("%s", err.message)
}

func NewRepositoryError(msg string) RepositoryError {
	return RepositoryError{message: msg}
}

var (
	ErrNotFound            = NewRepositoryError("404 key not found")
	ErrAlreadyExists       = NewRepositoryError("409 key already exists")
	ErrInsertOperationFail = NewRepositoryError("insert operation failed")
	ErrSelectOperationFail = NewRepositoryError("select operation failed")
	ErrUpdateOperationFail = NewRepositoryError("update operation failed")
	ErrDeleteOperationFail = NewRepositoryError("delete operation failed")
)
