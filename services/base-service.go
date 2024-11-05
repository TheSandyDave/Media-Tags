package services

import (
	"context"
	"errors"
	"reflect"

	apierrors "github.com/TheSandyDave/Media-Tags/api_errors"
	"github.com/TheSandyDave/Media-Tags/domain"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// compile time check for the struct implementing the interface
var _ BaseService[domain.IDbObject] = (*baseService[domain.IDbObject])(nil)

type baseService[T domain.IDbObject] struct {
	Database *gorm.DB
}

type Option[T domain.IDbObject] func(*gorm.DB) *gorm.DB

type BaseService[T domain.IDbObject] interface {
	Get(ctx context.Context, options ...Option[T]) ([]*T, error)
	GetWithID(ctx context.Context, id uuid.UUID, options ...Option[T]) (*T, error)
	GetWithIDs(ctx context.Context, ids []uuid.UUID, options ...Option[T]) ([]*T, error)
	Create(ctx context.Context, item ...*T) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func (service *baseService[T]) Get(ctx context.Context, options ...Option[T]) ([]*T, error) {
	logger := utils.NewLogger(ctx)

	dbQuery := service.Database.WithContext(ctx)
	for _, option := range options {
		dbQuery = option(dbQuery)
	}

	var result []*T
	if err := dbQuery.Find(&result).Error; err != nil {
		logger.
			WithField("model", reflect.TypeFor[T]().String()).
			WithError(err).
			Error("failed to get")

		return nil, err
	}

	return result, nil
}

func (service *baseService[T]) GetWithID(ctx context.Context, id uuid.UUID, options ...Option[T]) (*T, error) {
	logger := utils.NewLogger(ctx)

	dbQuery := service.Database.WithContext(ctx)
	for _, option := range options {
		dbQuery = option(dbQuery)
	}

	var result T
	if err := dbQuery.First(&result, id).Error; err != nil {

		logger.
			WithField("model", reflect.TypeFor[T]().String()).
			WithError(err).
			Error("failed to get with ID")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierrors.NewNotFoundError(id)
		}
		return nil, err
	}

	return &result, nil
}

func (service *baseService[T]) GetWithIDs(ctx context.Context, ids []uuid.UUID, options ...Option[T]) ([]*T, error) {
	logger := utils.NewLogger(ctx)

	dbQuery := service.Database.WithContext(ctx)
	for _, option := range options {
		dbQuery = option(dbQuery)
	}

	var result []*T
	if err := dbQuery.First(&result, ids).Error; err != nil {

		logger.
			WithField("model", reflect.TypeFor[T]().String()).
			WithError(err).
			Error("failed to get with ID")
		return nil, err
	}

	if len(result) == len(ids) {
		return result, nil
	}

	var missing []uuid.UUID
	for _, id := range ids {
		if !domain.ContainsID(id, result) {
			missing = append(missing, id)
		}
	}

	err := apierrors.NewRecordsNotFoundWithIDs(missing)
	return nil, err
}

func (service *baseService[T]) Create(ctx context.Context, item ...*T) error {
	logger := utils.NewLogger(ctx)

	if err := service.Database.WithContext(ctx).Create(item).Error; err != nil {
		logger.WithError(err).Error("failed creating")
		return err
	}

	return nil
}

func (service *baseService[T]) Delete(ctx context.Context, id uuid.UUID) error {
	logger := utils.NewLogger(ctx)
	if err := service.Database.WithContext(ctx).Delete(new(T), id).Error; err != nil {
		logger.WithError(err).Error("failed deleting")
		return err
	}

	return nil
}
