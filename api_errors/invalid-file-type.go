package apierrors

import (
	"context"
	"fmt"
	"net/http"
)

type InvalidFileTypeError struct {
	expectedFileType string
}

func (err *InvalidFileTypeError) Error() string {
	return fmt.Sprintf("invalid file type for uploaded file,expected %s", err.expectedFileType)
}

func NewInvalidFileTypeError(expectedFileType string) error {
	return &InvalidFileTypeError{
		expectedFileType: expectedFileType,
	}
}

func HandleInvalidFileTypeError(ctx context.Context, err *InvalidFileTypeError) (int, any) {
	return http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	}
}
