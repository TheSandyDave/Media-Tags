package controllers

import (
	"context"
	"net/http"

	apierrors "github.com/TheSandyDave/Media-Tags/api_errors"
	"github.com/TheSandyDave/Media-Tags/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type listFunction[Input any, Output any] interface {
	func(ctx context.Context, input Input) ([]Output, error)
}

type function[Input any, Output any] interface {
	func(ctx context.Context, input Input) (*Output, error)
}

func list[Input, Output any, F listFunction[Input, Output]](c *gin.Context, callback F) {
	logger := utils.NewLogger(c.Request.Context())

	var input Input
	if err := c.BindQuery(&input); err != nil {
		logger.WithError(c.Error(err)).Error("failed Binding list input")
		return
	}
	output, err := callback(c.Request.Context(), input)
	if err != nil {
		logger.WithError(c.Error(err)).Error("list operation failed")
		return
	}

	c.JSON(http.StatusOK, output)
}

func create[Input, Output any, F function[Input, Output]](c *gin.Context, callback F) {
	logger := utils.NewLogger(c.Request.Context())

	var input Input

	if err := c.Bind(&input); err != nil {
		logger.WithError(c.Error(err)).Error("failed Binding create input")
		return
	}
	output, err := callback(c.Request.Context(), input)
	if err != nil {
		logger.WithError(c.Error(err)).Error("create operation failed")
		return
	}

	c.JSON(http.StatusCreated, output)
}

func getWithID[Output any, F function[uuid.UUID, Output]](c *gin.Context, callback F) {
	logger := utils.NewLogger(c.Request.Context())

	var input struct {
		ID string `uri:"id" binding:"required"`
	}

	if err := c.BindUri(&input); err != nil {
		logger.WithError(c.Error(err)).Error("failed Binding URI input")
		return
	}

	id, err := uuid.Parse(input.ID)
	if err != nil {
		logger.WithField("ID", input.ID).WithError(err).Error("failed parsing id")
		c.Error(apierrors.NewInvalidUUIDError(input.ID))
		return
	}
	output, err := callback(c.Request.Context(), id)
	if err != nil {
		logger.WithError(c.Error(err)).Error("getWithID operation failed")
		return
	}

	c.JSON(http.StatusOK, output)
}
