package serror

import (
	"fmt"
)

type InternalError error

var NewInternalError = func(err error) error { return fmt.Errorf("Internal error: %w", err) }

type InvalidReqError error

var NewInvalidReqError = func(err error) error {
	return fmt.Errorf("Invalid request: %w", err)

}
