package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/TheSandyDave/Media-Tags/domain"
	restgen "github.com/TheSandyDave/Media-Tags/generated/api"
	mock_services "github.com/TheSandyDave/Media-Tags/generated/mock/services"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_TagController_Get_WritesCorrectOutput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		models []*domain.Tag
		output []restgen.Tag
	}{
		"single output": {
			models: []*domain.Tag{
				{
					Name: "testTag1",
				},
			},
			output: []restgen.Tag{
				{
					Name: "testTag1",
				},
			},
		},
		"multiple outputs": {
			models: []*domain.Tag{
				{
					Name: "testTag1",
				},
				{
					Name: "testTag2",
				},
			},
			output: []restgen.Tag{
				{
					Name: "testTag1",
				},
				{
					Name: "testTag2",
				},
			},
		},
	}

	for name, testData := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tagService := mock_services.NewMockITagService(ctrl)

			tagService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(testData.models, nil)

			TagController := TagController{
				TagService: tagService,
			}

			writer := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(writer)
			var err error
			context.Request, err = http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Error(err)
			}

			// act
			TagController.GetTags(context)

			// Assert
			var result []restgen.Tag

			if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusOK, writer.Result())) {
				assert.Equal(t, len(testData.output), len(result))
				equal := slices.EqualFunc(testData.output, result, func(a, b restgen.Tag) bool {
					return a.Name == b.Name
				})
				assert.True(t, equal)
			}
		})
	}
}

func Test_TagController_GetWithID_WritesCorrectOutput(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedTag := domain.Tag{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "expectedTag",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)

	tagService.EXPECT().GetWithID(gomock.Any(), expectedTag.ID, gomock.Any()).Return(&expectedTag, nil)

	TagController := TagController{
		TagService: tagService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	var err error
	context.Request, err = http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Error(err)
	}

	context.Params = append(context.Params, gin.Param{Key: "id", Value: expectedTag.ID.String()})
	// act
	TagController.GetTagWithId(context)

	// Assert
	var result restgen.Tag
	if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusOK, writer.Result())) {
		assert.Equal(t, expectedTag.Name, result.Name)
	}

}

func Test_TagController_Create_WritesCorrectOutput(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedName := "expected"

	ExpectedCreate := restgen.CreateTag{
		Name: expectedName,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)

	tagService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	TagController := TagController{
		TagService: tagService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	var err error
	body, err := json.Marshal(&ExpectedCreate)
	if err != nil {
		t.Error(err)
	}

	context.Request, err = http.NewRequest(http.MethodPost, "https://example.com", io.NopCloser(bytes.NewBuffer(body)))
	if err != nil {
		t.Error(err)
	}

	context.Request.Header.Set("Content-Type", "application/json")

	// act
	TagController.CreateTag(context)

	// Assert
	var result restgen.Tag
	if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusCreated, writer.Result())) {
		assert.Equal(t, expectedName, result.Name)
	}

}
