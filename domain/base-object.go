package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDbObject interface {
	GetID() uuid.UUID
}

func ContainsID[T IDbObject](id uuid.UUID, items []*T) bool {
	for _, item := range items {
		if (*item).GetID() == id {
			return true
		}
	}
	return false
}

type BaseObject struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (baseObject *BaseObject) BeforeCreate(_ *gorm.DB) error {
	if baseObject.ID == uuid.Nil {
		baseObject.ID = uuid.New()
	}
	return nil
}

func (BaseObject BaseObject) GetID() uuid.UUID {
	return BaseObject.ID
}
