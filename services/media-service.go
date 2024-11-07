package services

import (
	"github.com/TheSandyDave/Media-Tags/domain"
	"gorm.io/gorm"
)

// compile time check for the struct implementing the interface
var _ IMediaService = (*mediaService)(nil)

type IMediaService interface {
	IBaseService[domain.Media]
	FilterByTagOption(tag string) Option[domain.Media]
}

type mediaService struct {
	baseService[domain.Media]
}

func NewMediaService(db *gorm.DB) IMediaService {
	return &mediaService{
		baseService: baseService[domain.Media]{
			Database: db,
		},
	}
}

func (service *mediaService) FilterByTagOption(tag string) Option[domain.Media] {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins("INNER JOIN media_tags  as mt on media.id = mt.media_id").Joins("INNER JOIN tags on mt.tag_id = tags.id").Where("tags.name = ?", tag).Group("media.id")
	}
}
