package router

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/TheSandyDave/Media-Tags/api_errors"
	"github.com/TheSandyDave/Media-Tags/controllers"
	"github.com/TheSandyDave/Media-Tags/domain"
	restgen "github.com/TheSandyDave/Media-Tags/generated/api"
	"github.com/TheSandyDave/Media-Tags/services"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/flowchartsman/swaggerui"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/ing-bank/ginerr/v2"
	"gorm.io/gorm"
)

// format keeps logs consistent between logrus and gin logs
var loggerFormatString = "time=\"%s\" level=\"info\" statusCode=\"%v\" path=\"%s\" latency=\"%s\" method=\"%s\" origin=\"%s\" bodySize=\"%d\"\n"

type TaggedMediaAPI struct {
	Spec     []byte
	router   *gin.Engine
	database *gorm.DB

	// Controllers
	tagController   controllers.TagController
	mediaController controllers.MediaController
}

func (api *TaggedMediaAPI) Configure(ctx context.Context) *gin.Engine {
	logger := utils.NewLogger(ctx)
	logger.Info("configuring the API")

	api.router = gin.New()

	api.router.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: func(params gin.LogFormatterParams) string {
				timestamp := params.TimeStamp.Format("2006-01-02T15:04:05Z")
				return fmt.Sprintf(loggerFormatString, timestamp, params.StatusCode, params.Path, params.Latency, params.Method, params.ClientIP, params.BodySize)
			},

			// Output to stdout
			Output: gin.DefaultWriter,

			// Skip swagger path
			SkipPaths: []string{"/swagger"},
		}),
		gin.Recovery(),
	)

	api.configureDatabase(ctx)
	api.configureErrorHandlers(ctx)
	api.configureControllers()
	api.configureRoutes()

	if err := api.database.AutoMigrate(domain.Models...); err != nil {
		logger.Fatal(err)
	}
	return api.router
}

func (api *TaggedMediaAPI) configureDatabase(ctx context.Context) {
	logger := utils.NewLogger(ctx)

	if api.database == nil {
		database, err := gorm.Open(sqlite.Open("db"))
		if err != nil {
			logger.WithError(err).Fatal("failed to configure database")
		}

		api.database = database
	}
}

func (api *TaggedMediaAPI) configureErrorHandlers(ctx context.Context) {
	errorRegistry := ginerr.NewErrorRegistry()

	ginerr.RegisterErrorHandlerOn(errorRegistry, apierrors.HandleInvalidTagsError)
	ginerr.RegisterErrorHandlerOn(errorRegistry, apierrors.HandleInvalidUUIDError)
	ginerr.RegisterErrorHandlerOn(errorRegistry, apierrors.HandleRecordNotFoundError)
	ginerr.RegisterErrorHandlerOn(errorRegistry, apierrors.HandleInvalidFileTypeError)
	ginerr.RegisterErrorHandlerOn(errorRegistry, apierrors.HandleRequiredValueMissingError)

	errorRegistry.RegisterDefaultHandler(apierrors.DefaultErrorHandler)

	api.router.Use(func(c *gin.Context) {
		ctx = c.Request.Context()

		c.Next()

		ginErr := c.Errors.Last()
		if ginErr != nil {
			c.JSON(ginerr.NewErrorResponseFrom(errorRegistry, ctx, ginErr.Err))
		}
	})
}

func (api *TaggedMediaAPI) configureControllers() {
	var (
		tagService   = services.NewTagService(api.database)
		mediaService = services.NewMediaService(api.database)
	)

	api.tagController = controllers.TagController{
		TagService: tagService,
	}

	api.mediaController = controllers.MediaController{
		MediaService: mediaService,
		TagService:   tagService,
		IsTest:       false,
	}

}

func (api *TaggedMediaAPI) configureRoutes() {
	// Setup swagger UI
	api.router.GET("/swagger/*any", gin.WrapH(http.StripPrefix("/swagger", swaggerui.Handler(api.Spec))))

	//redirect default path to swagger
	api.router.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/swagger"
		api.router.HandleContext(c)
	})

	// Media file storage
	api.router.StaticFS("/files", gin.Dir("static", false))

	handlers := restgen.Handlers{
		// Tags
		CreateTag:  api.tagController.CreateTag,
		GetTags:    api.tagController.GetTags,
		GetTagById: api.tagController.GetTagWithId,

		// Media

		CreateMedia:  api.mediaController.CreateMedia,
		GetMedia:     api.mediaController.GetMedia,
		GetMediaById: api.mediaController.GetMediaWithId,
	}
	routes := restgen.GetRoutes(handlers)
	restgen.Decorate(api.router, routes)
}
