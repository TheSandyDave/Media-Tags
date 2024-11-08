/*
 * Tag and Media API
 *
 * API for managing tags and media items
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package restgen

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// list of generated routes
type Routes []Route

type Handlers struct {
	CreateMedia func(c *gin.Context)

	GetMedia func(c *gin.Context)

	GetMediaById func(c *gin.Context)

	CreateTag func(c *gin.Context)

	GetTagById func(c *gin.Context)

	GetTags func(c *gin.Context)
}

func GetRoutes(handlers Handlers) Routes {
	return Routes{

		{
			"CreateMedia",
			http.MethodPost,
			"/media",
			handlers.CreateMedia,
		},

		{
			"GetMedia",
			http.MethodGet,
			"/media",
			handlers.GetMedia,
		},

		{
			"GetMediaById",
			http.MethodGet,
			"/media/:id",
			handlers.GetMediaById,
		},

		{
			"CreateTag",
			http.MethodPost,
			"/tags",
			handlers.CreateTag,
		},

		{
			"GetTagById",
			http.MethodGet,
			"/tags/:id",
			handlers.GetTagById,
		},

		{
			"GetTags",
			http.MethodGet,
			"/tags",
			handlers.GetTags,
		},
	}
}

// decorate a gin Engine with routes
func Decorate(router *gin.Engine, routes Routes) {
	for _, route := range routes {
		router.Handle(route.Method, route.Pattern, route.HandlerFunc)
	}
}
