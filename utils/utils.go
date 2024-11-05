package utils

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewLogger(ctx context.Context) *logrus.Entry {
	return logrus.WithContext(ctx)
}

func IDStringSlice(source []uuid.UUID) []string {
	slice := make([]string, len(source))
	for _, uuid := range source {
		slice = append(slice, uuid.String())
	}
	return slice
}
