package serror

import "fmt"

type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return fmt.Sprint("Internal error: %w", e.Err)
}
