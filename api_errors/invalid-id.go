package apierrors

import (
	"context"
	"fmt"
	"net/http"
)

type InvalidUUIDError struct {
	invalidValue string
}

func (err *InvalidUUIDError) Error() string {
	return fmt.Sprintf("invalid value supplied {%s},expected a UUID", err.invalidValue)
}

func NewInvalidUUIDError(invalidValue string) error {
	return &InvalidUUIDError{
		invalidValue: invalidValue,
	}
}

func HandleInvalidUUIDError(ctx context.Context, err *InvalidUUIDError) (int, any) {
	return http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	}
}
