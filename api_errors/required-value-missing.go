package apierrors

import (
	"context"
	"fmt"
	"net/http"
)

type RequiredValueMissingError struct {
	missingValue string
}

func (err *RequiredValueMissingError) Error() string {
	return fmt.Sprintf("required Value %s missing", err.missingValue)
}

func NewRequiredValueMissingError(missingValue string) error {
	return &RequiredValueMissingError{
		missingValue: missingValue,
	}
}

func HandleRequiredValueMissingError(ctx context.Context, err *RequiredValueMissingError) (int, any) {
	return http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	}
}
