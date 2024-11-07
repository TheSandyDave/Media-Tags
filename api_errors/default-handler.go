package apierrors

import (
	"context"
	"net/http"

	"github.com/TheSandyDave/Media-Tags/utils"
)

func DefaultErrorHandler(ctx context.Context, err error) (int, any) {
	logger := utils.NewLogger(ctx)

	logger.WithError(err).Warn("no error handler found for error")

	return http.StatusInternalServerError, ErrorResponse{
		Error: "Internal server error",
	}
}
