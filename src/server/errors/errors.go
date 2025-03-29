package serror

import "fmt"

type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("Internal error: %w", e.Err)
}

type SharedDocNotFound struct{}

func (e *SharedDocNotFound) Error() string {
	return "Shared document not found"
}

type InvalidReqError struct {
	spec string
}

func (e *InvalidReqError) Error() string {
	return fmt.Sprintf("Invalid request: %s", e.spec)
}

type DocIdError struct {
	Err error
}

func (e *DocIdError) Error() string {
	return fmt.Sprintf("Request doc id is invalid: %w", e.Err)
}
