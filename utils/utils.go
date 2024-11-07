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
	for i, uuid := range source {
		slice[i] = uuid.String()
	}
	return slice
}

func StringSliceToUUID(source []string) ([]uuid.UUID, error) {
	slice := make([]uuid.UUID, len(source))
	for i, idStr := range source {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		slice[i] = id
	}
	return slice, nil
}
