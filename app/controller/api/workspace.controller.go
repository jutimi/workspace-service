package api

import (
	"context"
	"net/http"
	"time"
	"workspace-server/app/model"
	"workspace-server/app/service"
	_errors "workspace-server/package/errors"
	"workspace-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type workspaceHandler struct {
	services service.ServiceCollections
}

func NewApiWorkspaceController(router *gin.Engine, services service.ServiceCollections) {
	handler := workspaceHandler{services}

	group := router.Group("api/v1/workspaces")
	{
		group.POST("/create", handler.create)
	}
}

func (h *workspaceHandler) create(c *gin.Context) {
	var data model.CreateWorkspaceRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.WorkspaceSvc.CreateWorkspace(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
