package controllers

import (
	"context"

	"github.com/TheSandyDave/Media-Tags/conversion"
	"github.com/TheSandyDave/Media-Tags/domain"
	restgen "github.com/TheSandyDave/Media-Tags/generated/api"
	"github.com/TheSandyDave/Media-Tags/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TagController struct {
	TagService services.ITagService
}

func (controller *TagController) GetTags(c *gin.Context) {
	list(c, func(ctx context.Context, _ any) ([]*restgen.Tag, error) {

		tags, err := controller.TagService.Get(ctx)
		if err != nil {
			return nil, err
		}

		return conversion.EncodeSlice(tags, conversion.EncodeTag), nil
	})
}

func (controller *TagController) GetTagWithId(c *gin.Context) {
	getWithID(c, func(ctx context.Context, id uuid.UUID) (*restgen.Tag, error) {
		tag, err := controller.TagService.GetWithID(ctx, id)
		if err != nil {
			return nil, err
		}

		return conversion.EncodeTag(tag), nil
	})
}

func (controller *TagController) CreateTag(c *gin.Context) {
	create(c, func(ctx context.Context, input restgen.CreateTag) (*restgen.Tag, error) {
		tag := &domain.Tag{
			Name: input.Name,
		}

		if err := controller.TagService.Create(ctx, tag); err != nil {
			return nil, err
		}

		return conversion.EncodeTag(tag), nil
	})
}
