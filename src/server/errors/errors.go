package serror

import (
	"fmt"
)

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
	Err error
}

func (e *InvalidReqError) Error() string {
	return fmt.Sprintf("Invalid request: %s", e.Err)
}

type DocIdError struct {
}

func (e *DocIdError) Error() string {
	return fmt.Sprintf("Request doc id is invalid")
}
