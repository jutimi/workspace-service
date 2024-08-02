package cms

import (
	"context"
	"net/http"
	"time"

	"workspace-server/app/middleware"
	"workspace-server/app/model"
	"workspace-server/app/service"
	"workspace-server/utils"

	_errors "workspace-server/package/errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type workspaceHandler struct {
	services   service.ServiceCollections
	middleware middleware.MiddlewareCollections
}

func NewApiWorkspaceController(
	router *gin.Engine,
	services service.ServiceCollections,
	middleware middleware.MiddlewareCollections,
) {
	handler := workspaceHandler{services, middleware}

	group := router.Group("api/v1/workspaces", middleware.WorkspaceMW.Handler())
	{
		group.PUT("/update", handler.update)
	}
}

func (h *workspaceHandler) update(c *gin.Context) {
	var data model.UpdateWorkspaceRequest
	if err := c.ShouldBindBodyWith(&data, binding.JSON); err != nil {
		resErr := _errors.NewValidatorError(err)
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(resErr))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, utils.GIN_CONTEXT_KEY, c)

	res, err := h.services.WorkspaceSvc.UpdateWorkspace(ctx, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.FormatErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.FormatSuccessResponse(res))
}
