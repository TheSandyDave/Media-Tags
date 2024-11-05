package apierrors

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/google/uuid"
)

type InvalidTagsError struct {
	tagIDs []uuid.UUID
}

func (err *InvalidTagsError) Error() string {
	return fmt.Sprintf("Tags with following IDs could not be found: %s", strings.Join(utils.IDStringSlice(err.tagIDs), ","))
}

func NewInvalidTagsError(IDs []uuid.UUID) error {
	return &InvalidTagsError{
		tagIDs: IDs,
	}
}

func HandleInvalidTagsError(ctx context.Context, err *InvalidTagsError) (int, any) {
	return http.StatusBadRequest, ErrorResponse{
		Error: err.Error(),
	}
}
