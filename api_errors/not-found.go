package apierrors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/google/uuid"
)

type RecordNotFoundError struct {
	ID uuid.UUID
}

func (err *RecordNotFoundError) Error() string {
	return fmt.Sprintf("Record with ID {%s} not found", err.ID.String())
}

func NewNotFoundError(ID uuid.UUID) error {
	return &RecordNotFoundError{
		ID: ID,
	}
}

func HandleRecordNotFoundError(ctx context.Context, err *RecordNotFoundError) (int, any) {
	return http.StatusNotFound, ErrorResponse{
		Error: err.Error(),
	}
}

type RecordsNotFoundWithIDs struct {
	IDs []uuid.UUID
}

func (err *RecordsNotFoundWithIDs) Error() string {
	return fmt.Sprintf("Record with IDs {%s} not found", strings.Join(utils.IDStringSlice(err.IDs), ","))
}

func NewRecordsNotFoundWithIDs(IDs []uuid.UUID) error {
	return &RecordsNotFoundWithIDs{
		IDs: IDs,
	}
}
