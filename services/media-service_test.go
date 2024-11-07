package services

import (
	"context"
	"slices"
	"testing"

	"github.com/TheSandyDave/Media-Tags/domain"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FilterByTagOption_FilterTests(t *testing.T) {
	// Arrange
	tag1 := domain.Tag{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "TestTag1",
	}
	tag2 := domain.Tag{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "TestTag2",
	}
	media1 := domain.Media{
		Name: "testMedia1",
		Tags: []*domain.Tag{
			&tag1,
		},
	}
	media2 := domain.Media{
		Name: "testMedia2",
		Tags: []*domain.Tag{
			&tag1,
			&tag2,
		},
	}

	database := utils.NewInMemoryDatabase(t)
	require.NoError(t, database.Create(&[]*domain.Media{&media1, &media2}).Error)
	service := NewMediaService(database)
	testCases := map[string]struct {
		FilterOption   string
		expectedOutput []*domain.Media
	}{
		"multiple media match tag filter": {
			FilterOption: tag1.Name,
			expectedOutput: []*domain.Media{
				&media1,
				&media2,
			},
		},
		"single media matches tag filter": {
			FilterOption: tag2.Name,
			expectedOutput: []*domain.Media{
				&media2,
			},
		},
		"nothing matches tag filter": {
			FilterOption:   "nonExistantTag",
			expectedOutput: []*domain.Media{},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// Act
			res, err := service.Get(context.Background(), service.FilterByTagOption(testCase.FilterOption))

			// Assert
			assert.NoError(t, err)
			if assert.NotNil(t, res) {
				equal := slices.EqualFunc(testCase.expectedOutput, res, func(a, b *domain.Media) bool {
					return a.Name == b.Name
				})
				assert.True(t, equal)
			}
		})
	}

}
