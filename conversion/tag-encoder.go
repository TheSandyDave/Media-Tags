package conversion

import (
	"github.com/TheSandyDave/Media-Tags/domain"
	restgen "github.com/TheSandyDave/Media-Tags/generated/api"
)

func EncodeTag(source *domain.Tag) *restgen.Tag {
	return &restgen.Tag{
		Id:   source.ID.String(),
		Name: source.Name,
	}
}
