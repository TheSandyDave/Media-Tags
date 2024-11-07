package conversion

import (
	"github.com/TheSandyDave/Media-Tags/domain"
	restgen "github.com/TheSandyDave/Media-Tags/generated/api"
)

func EncodeMedia(source *domain.Media) *restgen.Media {
	tags := make([]string, len(source.Tags))
	for i, tag := range source.Tags {
		tags[i] = tag.Name
	}
	return &restgen.Media{
		Id:      source.ID.String(),
		Name:    source.Name,
		Tags:    tags,
		FileUrl: source.FileUrl,
	}
}
