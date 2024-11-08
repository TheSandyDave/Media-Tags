package controllers

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"os"
	"slices"
	"testing"

	apierrors "github.com/TheSandyDave/Media-Tags/api_errors"
	"github.com/TheSandyDave/Media-Tags/domain"
	restgen "github.com/TheSandyDave/Media-Tags/generated/api"
	mock_services "github.com/TheSandyDave/Media-Tags/generated/mock/services"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_MediaController_Get_WritesCorrectOutput(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		models []*domain.Media
		output []restgen.Media
	}{
		"single output": {
			models: []*domain.Media{
				{
					Name: "testMedia1",
				},
			},
			output: []restgen.Media{
				{
					Name: "testMedia1",
				},
			},
		},
		"multiple outputs": {
			models: []*domain.Media{
				{
					Name: "testMedia1",
				},
				{
					Name: "testMedia2",
				},
			},
			output: []restgen.Media{
				{
					Name: "testMedia1",
				},
				{
					Name: "testMedia2",
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
			mediaService := mock_services.NewMockIMediaService(ctrl)

			mediaService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(testData.models, nil)

			MediaController := MediaController{
				TagService:   tagService,
				MediaService: mediaService,
			}

			writer := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(writer)
			var err error
			context.Request, err = http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Error(err)
			}
			// act
			MediaController.GetMedia(context)

			// Assert
			var result []restgen.Media
			if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusOK, writer.Result())) {
				assert.Equal(t, len(testData.output), len(result))
				equal := slices.EqualFunc(testData.output, result, func(a, b restgen.Media) bool {
					return a.Name == b.Name
				})
				assert.True(t, equal)
			}
		})
	}
}
func Test_MediaController_Get_CreatesAFilterWhenQueryParamIsPassed(t *testing.T) {
	t.Parallel()

	// Arrange

	models := []*domain.Media{
		{
			Name: "testMedia1",
		},
	}
	output := []restgen.Media{
		{
			Name: "testMedia1",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)
	mediaService := mock_services.NewMockIMediaService(ctrl)
	expectedFilter := "test filter tag"

	mediaService.EXPECT().FilterByTagOption(expectedFilter).Return(nil)
	mediaService.EXPECT().Get(gomock.Any(), gomock.Any()).Return(models, nil)

	MediaController := MediaController{
		TagService:   tagService,
		MediaService: mediaService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	var err error
	context.Request, err = http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Error(err)
	}
	query := url.Values{}
	query.Add("tag", expectedFilter)
	context.Request.URL.RawQuery = query.Encode()

	// act
	MediaController.GetMedia(context)

	// Assert
	var result []restgen.Media
	if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusOK, writer.Result())) {
		assert.Equal(t, len(output), len(result))
		equal := slices.EqualFunc(output, result, func(a, b restgen.Media) bool {
			return a.Name == b.Name
		})
		assert.True(t, equal)
	}
}

func Test_MediaController_GetWithID_WritesCorrectOutput(t *testing.T) {
	t.Parallel()

	// Arrange
	expectedMedia := domain.Media{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "expectedTag",
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)
	mediaService := mock_services.NewMockIMediaService(ctrl)

	mediaService.EXPECT().GetWithID(gomock.Any(), expectedMedia.ID, gomock.Any()).Return(&expectedMedia, nil)

	MediaController := MediaController{
		TagService:   tagService,
		MediaService: mediaService,
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	var err error
	context.Request, err = http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Error(err)
	}
	context.Params = append(context.Params, gin.Param{Key: "id", Value: expectedMedia.ID.String()})

	// act
	MediaController.GetMediaWithId(context)

	// Assert
	var result restgen.Media
	if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusOK, writer.Result())) {
		assert.Equal(t, expectedMedia.Name, result.Name)
	}

}

func Test_MediaController_Create_FailsIfNoFileIsProvided(t *testing.T) {
	t.Parallel()

	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)
	mediaService := mock_services.NewMockIMediaService(ctrl)

	MediaController := MediaController{
		TagService:   tagService,
		MediaService: mediaService,
	}

	body := new(bytes.Buffer)
	fileHeader := make(textproto.MIMEHeader)
	multipartWriter := multipart.NewWriter(body)
	defer multipartWriter.Close()
	_, err := multipartWriter.CreatePart(fileHeader)
	if err != nil {
		t.Error(err)
	}

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, err = http.NewRequest(http.MethodPost, "https://example.com", body)
	if err != nil {
		t.Error(err)
	}

	// act
	MediaController.CreateMedia(context)

	// Assert

	// the error handler middleware is not initialized here so the appropriate error response is not initialized
	// checking that the error is in the stack instead
	assert.Contains(t, context.Errors.Last().Err.Error(), "required Value file missing")
}

func Test_MediaController_Create_FailsIfUploadedFileIsNotAnImage(t *testing.T) {
	t.Parallel()

	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)
	mediaService := mock_services.NewMockIMediaService(ctrl)

	MediaController := MediaController{
		TagService:   tagService,
		MediaService: mediaService,
	}

	body := new(bytes.Buffer)
	fileHeader := make(textproto.MIMEHeader)
	multipartWriter := multipart.NewWriter(body)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "test.txt"))
	fileHeader.Set("Content-Type", "text/plain")
	fileWriter, err := multipartWriter.CreatePart(fileHeader)
	if err != nil {
		t.Error(err)
	}
	file, err := os.Open("test-resources/test.txt")
	if err != nil {
		t.Error(err)
	}

	io.Copy(fileWriter, file)
	multipartWriter.Close()

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, err = http.NewRequest(http.MethodPost, "https://example.com", body)
	if err != nil {
		t.Error(err)
	}
	context.Request.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	// act
	MediaController.CreateMedia(context)

	// Assert

	// the error handler middleware is not initialized here so the appropriate error response is not initialized
	// checking that the error is in the stack instead
	assert.Contains(t, context.Errors.Last().Err.Error(), "invalid file type for uploaded file,expected image")
}

func Test_MediaController_Create_FailsIfTagDoesNotExist(t *testing.T) {
	t.Parallel()

	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)
	mediaService := mock_services.NewMockIMediaService(ctrl)

	MediaController := MediaController{
		TagService:   tagService,
		MediaService: mediaService,
	}

	body := new(bytes.Buffer)
	fileHeader := make(textproto.MIMEHeader)
	multipartWriter := multipart.NewWriter(body)

	// we are relying on the contentType header to determine if it's an image so we can just pass any file for testing
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "test.png"))
	fileHeader.Set("Content-Type", "image/png")
	fileWriter, err := multipartWriter.CreatePart(fileHeader)
	if err != nil {
		t.Error(err)
	}
	file, err := os.Open("test-resources/test.txt")
	if err != nil {
		t.Error(err)
	}

	io.Copy(fileWriter, file)
	expectedTagKey := uuid.NewString()
	multipartWriter.WriteField("tags", expectedTagKey)
	multipartWriter.Close()

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, err = http.NewRequest(http.MethodPost, "https://example.com", body)
	if err != nil {
		t.Error(err)
	}
	context.Request.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	uuids := []uuid.UUID{uuid.MustParse(expectedTagKey)}
	tagService.EXPECT().GetWithIDs(gomock.Any(), uuids, gomock.Any()).
		Return(nil, apierrors.NewRecordsNotFoundWithIDs(uuids))

	// act
	MediaController.CreateMedia(context)

	// Assert

	// the error handler middleware is not initialized here so the appropriate error response is not initialized
	// checking that the error is in the stack instead
	assert.Contains(t, context.Errors.Last().Err.Error(), "Tags with following IDs could not be found:")
}

func Test_MediaController_Create_WritesCorrectOutput(t *testing.T) {
	t.Parallel()

	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagService := mock_services.NewMockITagService(ctrl)
	mediaService := mock_services.NewMockIMediaService(ctrl)

	MediaController := MediaController{
		TagService:   tagService,
		MediaService: mediaService,
		IsTest:       false,
	}

	body := new(bytes.Buffer)
	fileHeader := make(textproto.MIMEHeader)
	multipartWriter := multipart.NewWriter(body)

	// we are relying on the contentType header to determine if it's an image so we can just pass any file for testing
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "test.png"))
	fileHeader.Set("Content-Type", "image/png")
	fileWriter, err := multipartWriter.CreatePart(fileHeader)
	if err != nil {
		t.Error(err)
	}
	file, err := os.Open("test-resources/test.txt")
	if err != nil {
		t.Error(err)
	}

	io.Copy(fileWriter, file)
	expectedTag := domain.Tag{
		BaseObject: domain.BaseObject{
			ID: uuid.New(),
		},
		Name: "expectedTag",
	}

	multipartWriter.WriteField("tags", expectedTag.ID.String())

	expectedName := "expectedMedia"
	multipartWriter.WriteField("name", expectedName)
	multipartWriter.Close()

	writer := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(writer)
	context.Request, err = http.NewRequest(http.MethodPost, "https://example.com", body)
	if err != nil {
		t.Error(err)
	}
	context.Request.Header.Add("Content-Type", multipartWriter.FormDataContentType())
	tagService.EXPECT().GetWithIDs(gomock.Any(), []uuid.UUID{expectedTag.ID}, gomock.Any()).
		Return([]*domain.Tag{&expectedTag}, nil)

	mediaService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	// act
	MediaController.CreateMedia(context)

	// Assert
	var result restgen.Media
	if assert.True(t, utils.RetrieveResponse(t, &result, http.StatusCreated, writer.Result())) {
		assert.Equal(t, expectedName, result.Name)
		assert.Equal(t, expectedTag.Name, result.Tags[0])
	}

}
