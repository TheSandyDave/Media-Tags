package services

import (
	"github.com/TheSandyDave/Media-Tags/domain"
	"gorm.io/gorm"
)

// compile time check for the struct implementing the interface
var _ ITagService = (*tagService)(nil)

type ITagService interface {
	IBaseService[domain.Tag]
}

type tagService struct {
	baseService[domain.Tag]
}

func NewTagService(db *gorm.DB) ITagService {
	return &tagService{
		baseService: baseService[domain.Tag]{
			Database: db,
		},
	}
}
